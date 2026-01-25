package audit

import (
	"context"
	"fmt"
	"log/slog"
	"sort"
	"sync"
	"time"

	"github.com/itspablomontes/fleming/apps/backend/internal/common"
	"github.com/itspablomontes/fleming/pkg/protocol/audit"
	"github.com/itspablomontes/fleming/pkg/protocol/types"
)

// Service defines the business logic for the audit protocol.
type Service interface {
	Record(ctx context.Context, actor string, action audit.Action, resourceType audit.ResourceType, resourceID string, metadata common.JSONMap) error
	GetLatestEntries(ctx context.Context, actor string, limit int) ([]AuditEntry, error)
	VerifyIntegrity(ctx context.Context) (bool, error)
	BuildMerkleTree(ctx context.Context, startTime time.Time, endTime time.Time) (*AuditBatch, *audit.MerkleTree, error)
	GetMerkleRoot(ctx context.Context, batchID string) (string, error)
	VerifyMerkleProof(root string, entryHash string, proof *audit.Proof) bool
	GetEntriesForMerkle(ctx context.Context, startTime time.Time, endTime time.Time) ([]AuditEntry, error)
	GetEntryByID(ctx context.Context, id string) (*AuditEntry, error)
	GetEntriesByResource(ctx context.Context, resourceID string) ([]AuditEntry, error)
	QueryEntries(ctx context.Context, filter audit.QueryFilter) ([]AuditEntry, error)
}

type service struct {
	repo Repository
	mu   sync.Mutex // Ensure sequential hashing if multiple records happen at once
}

// NewService creates a new audit service.
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// Record generates a new cryptographically chained audit entry.
func (s *service) Record(ctx context.Context, actor string, action audit.Action, resourceType audit.ResourceType, resourceID string, metadata common.JSONMap) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	latest, err := s.repo.GetLatest(ctx)
	if err != nil {
		return fmt.Errorf("audit: %w", err)
	}

	previousHash := "GENESIS"
	if latest != nil {
		previousHash = latest.Hash
	}

	protocolEntry := audit.NewEntry(
		types.WalletAddress(actor),
		action,
		resourceType,
		types.ID(resourceID),
		previousHash,
	)

	if metadata != nil {
		for k, v := range metadata {
			protocolEntry.Metadata[k] = v
		}
		protocolEntry.SetHash()
	}

	dbEntry := &AuditEntry{
		Actor:        actor,
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		Timestamp:    protocolEntry.Timestamp,
		Metadata:     metadata,
		Hash:         protocolEntry.Hash,
		PreviousHash: protocolEntry.PreviousHash,
	}

	if err := s.repo.Create(ctx, dbEntry); err != nil {
		return fmt.Errorf("audit: %w", err)
	}

	slog.DebugContext(ctx, "audit entry recorded", "action", action, "hash", dbEntry.Hash)
	return nil
}

// GetLatestEntries returns the most recent audit logs.
func (s *service) GetLatestEntries(ctx context.Context, actor string, limit int) ([]AuditEntry, error) {
	if limit <= 0 {
		limit = 100
	}
	return s.repo.List(ctx, actor, limit)
}

// VerifyIntegrity checks the entire hash chain for tampering.
func (s *service) VerifyIntegrity(ctx context.Context) (bool, error) {
	entries, err := s.repo.List(ctx, "", 0)
	if err != nil {
		return false, err
	}

	for i := len(entries) - 1; i >= 0; i-- {
		e := entries[i]

		protocolEntry := audit.Entry{
			Actor:        types.WalletAddress(e.Actor),
			Action:       e.Action,
			ResourceType: e.ResourceType,
			ResourceID:   types.ID(e.ResourceID),
			Timestamp:    e.Timestamp,
			PreviousHash: e.PreviousHash,
		}

		computed := protocolEntry.ComputeHash()
		if computed != e.Hash {
			slog.ErrorContext(ctx, "audit integrity failure: hash mismatch", "id", e.ID, "expected", e.Hash, "computed", computed)
			return false, nil
		}

		if i < len(entries)-1 {
			prev := entries[i+1]
			if e.PreviousHash != prev.Hash {
				slog.ErrorContext(ctx, "audit integrity failure: chain broken", "id", e.ID, "previous_hash", e.PreviousHash, "prev_entry_hash", prev.Hash)
				return false, nil
			}
		}
	}

	return true, nil
}

func (s *service) GetEntriesForMerkle(ctx context.Context, startTime time.Time, endTime time.Time) ([]AuditEntry, error) {
	filter := audit.NewQueryFilter()
	if !startTime.IsZero() {
		ts := types.NewTimestamp(startTime)
		filter.StartTime = &ts
	}
	if !endTime.IsZero() {
		ts := types.NewTimestamp(endTime)
		filter.EndTime = &ts
	}
	filter.Limit = 0

	return s.repo.Query(ctx, filter)
}

func (s *service) GetEntryByID(ctx context.Context, id string) (*AuditEntry, error) {
	return s.repo.GetByID(ctx, types.ID(id))
}

func (s *service) GetEntriesByResource(ctx context.Context, resourceID string) ([]AuditEntry, error) {
	return s.repo.GetByResource(ctx, types.ID(resourceID))
}

func (s *service) QueryEntries(ctx context.Context, filter audit.QueryFilter) ([]AuditEntry, error) {
	return s.repo.Query(ctx, filter)
}

func (s *service) BuildMerkleTree(ctx context.Context, startTime time.Time, endTime time.Time) (*AuditBatch, *audit.MerkleTree, error) {
	entries, err := s.GetEntriesForMerkle(ctx, startTime, endTime)
	if err != nil {
		return nil, nil, fmt.Errorf("build merkle tree: %w", err)
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Timestamp.Equal(entries[j].Timestamp) {
			return entries[i].ID < entries[j].ID
		}
		return entries[i].Timestamp.Before(entries[j].Timestamp)
	})

	protocolEntries := make([]audit.Entry, 0, len(entries))
	for _, entry := range entries {
		protocolEntries = append(protocolEntries, audit.Entry{
			Hash: entry.Hash,
		})
	}

	tree, err := audit.BuildMerkleTree(protocolEntries)
	if err != nil {
		return nil, nil, fmt.Errorf("build merkle tree: %w", err)
	}

	batch := &AuditBatch{
		RootHash:   tree.Root,
		StartTime:  startTime.UTC(),
		EndTime:    endTime.UTC(),
		EntryCount: len(entries),
		CreatedAt:  time.Now().UTC(),
	}
	if err := s.repo.CreateBatch(ctx, batch); err != nil {
		return nil, nil, fmt.Errorf("create audit batch: %w", err)
	}

	return batch, tree, nil
}

func (s *service) GetMerkleRoot(ctx context.Context, batchID string) (string, error) {
	batch, err := s.repo.GetBatchByID(ctx, batchID)
	if err != nil {
		return "", err
	}
	if batch == nil {
		return "", nil
	}
	return batch.RootHash, nil
}

func (s *service) VerifyMerkleProof(root string, entryHash string, proof *audit.Proof) bool {
	return audit.VerifyProof(root, entryHash, proof)
}
