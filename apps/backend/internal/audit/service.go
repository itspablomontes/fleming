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
	BuildMerkleTree(ctx context.Context, actor string, startTime time.Time, endTime time.Time) (*AuditBatch, *audit.MerkleTree, error)
	GetBatch(ctx context.Context, actor string, batchID string) (*AuditBatch, error)
	GetBatchByRoot(ctx context.Context, actor string, rootHash string) (*AuditBatch, error)
	ListBatches(ctx context.Context, actor string, limit int, offset int) ([]AuditBatch, error)
	AnchorBatch(ctx context.Context, actor string, batchID string, chainClient ChainAnchorer) (*AuditBatch, error)
	VerifyMerkleProof(root string, entryHash string, proof *audit.Proof) bool
	GetEntriesForMerkle(ctx context.Context, actor string, startTime time.Time, endTime time.Time) ([]AuditEntry, error)
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
		Actor:         actor,
		Action:        action,
		ResourceType:  resourceType,
		ResourceID:    resourceID,
		Timestamp:     protocolEntry.Timestamp,
		Metadata:      metadata,
		Hash:          protocolEntry.Hash,
		PreviousHash:  protocolEntry.PreviousHash,
		SchemaVersion: protocolEntry.SchemaVersion,
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

func (s *service) GetEntriesForMerkle(ctx context.Context, actor string, startTime time.Time, endTime time.Time) ([]AuditEntry, error) {
	address, err := types.NewWalletAddress(actor)
	if err != nil {
		return nil, fmt.Errorf("audit: invalid actor address: %w", err)
	}

	filter := audit.NewQueryFilter()
	filter.Actor = address
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

func (s *service) BuildMerkleTree(ctx context.Context, actor string, startTime time.Time, endTime time.Time) (*AuditBatch, *audit.MerkleTree, error) {
	if actor == "" {
		return nil, nil, fmt.Errorf("build merkle tree: actor is required")
	}

	entries, err := s.GetEntriesForMerkle(ctx, actor, startTime, endTime)
	if err != nil {
		return nil, nil, fmt.Errorf("build merkle tree: %w", err)
	}
	if len(entries) == 0 {
		return nil, nil, fmt.Errorf("build merkle tree: no entries in range")
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

	existing, err := s.repo.GetBatchByActorAndRoot(ctx, actor, tree.Root)
	if err != nil {
		return nil, nil, fmt.Errorf("get audit batch by root: %w", err)
	}
	if existing != nil {
		return existing, tree, nil
	}

	batch := &AuditBatch{
		Actor:        actor,
		RootHash:     tree.Root,
		StartTime:    startTime.UTC(),
		EndTime:      endTime.UTC(),
		EntryCount:   len(entries),
		CreatedAt:    time.Now().UTC(),
		AnchorStatus: "pending",
	}
	if err := s.repo.CreateBatch(ctx, batch); err != nil {
		return nil, nil, fmt.Errorf("create audit batch: %w", err)
	}

	return batch, tree, nil
}

func (s *service) GetBatch(ctx context.Context, actor string, batchID string) (*AuditBatch, error) {
	if actor == "" {
		return nil, fmt.Errorf("get audit batch: actor is required")
	}
	if batchID == "" {
		return nil, fmt.Errorf("get audit batch: batch id is required")
	}

	batch, err := s.repo.GetBatchByIDForActor(ctx, actor, batchID)
	if err != nil {
		return nil, err
	}
	return batch, nil
}

func (s *service) GetBatchByRoot(ctx context.Context, actor string, rootHash string) (*AuditBatch, error) {
	if actor == "" {
		return nil, fmt.Errorf("get audit batch by root: actor is required")
	}
	if rootHash == "" {
		return nil, fmt.Errorf("get audit batch by root: root hash is required")
	}
	batch, err := s.repo.GetBatchByActorAndRoot(ctx, actor, rootHash)
	if err != nil {
		return nil, err
	}
	return batch, nil
}

func (s *service) ListBatches(ctx context.Context, actor string, limit int, offset int) ([]AuditBatch, error) {
	if actor == "" {
		return nil, fmt.Errorf("list audit batches: actor is required")
	}
	if limit <= 0 {
		limit = 25
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.ListBatchesByActor(ctx, actor, limit, offset)
}

func (s *service) VerifyMerkleProof(root string, entryHash string, proof *audit.Proof) bool {
	return audit.VerifyProof(root, entryHash, proof)
}
