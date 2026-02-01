package audit

import (
	"context"
	"fmt"
	"log/slog"
	"maps"
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
	VerifyIntegrity(ctx context.Context, actor string) (bool, error)
	BuildMerkleTree(ctx context.Context, actor string, startTime time.Time, endTime time.Time) (*AuditBatch, *audit.MerkleTree, error)
	GetBatch(ctx context.Context, actor string, batchID string) (*AuditBatch, error)
	GetBatchByRoot(ctx context.Context, actor string, rootHash string) (*AuditBatch, error)
	ListBatches(ctx context.Context, actor string, limit int, offset int) ([]AuditBatch, error)
	AnchorBatch(ctx context.Context, actor string, batchID string, chainClient ChainAnchorer) (*AuditBatch, error)
	VerifyMerkleProof(root string, entryHash string, proof *audit.Proof) bool
	GetEntriesForMerkle(ctx context.Context, actor string, startTime time.Time, endTime time.Time) ([]AuditEntry, error)
	GetEntryByID(ctx context.Context, id string) (*AuditEntry, error)
	GetEntriesByResource(ctx context.Context, resourceID string) ([]AuditEntry, error)

	QueryEntries(ctx context.Context, filter audit.QueryFilter) (*audit.QueryResult, error)
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

	// Normalize to canonical lowercased address for stable chaining and lookup.
	addr, err := types.NewWalletAddress(actor)
	if err != nil {
		return fmt.Errorf("audit: invalid actor address: %w", err)
	}
	actor = addr.String()

	latest, err := s.repo.GetLatest(ctx, actor)
	if err != nil {
		return fmt.Errorf("audit: %w", err)
	}

	previousHash := "GENESIS"
	if latest != nil {
		previousHash = latest.Hash
	}

	// Truncate to Microsecond to match Postgres timestamptz precision
	ts := time.Now().UTC()
	ts = ts.Truncate(time.Microsecond)
	// Ensure monotonic timestamps per actor to avoid ambiguous ordering when multiple
	// entries share the same microsecond (UUID PK ordering is not time-correlated).
	if latest != nil && !ts.After(latest.Timestamp) {
		ts = latest.Timestamp.Add(time.Microsecond)
	}

	entry := audit.NewEntry(
		types.WalletAddress(actor),
		action,
		resourceType,
		types.ID(resourceID),
		previousHash,
	)
	entry.Timestamp = ts

	if metadata != nil {
		maps.Copy(entry.Metadata, metadata)
	}
	entry.SetHash()

	dbEntry := &AuditEntry{
		Actor:         actor,
		Action:        action,
		ResourceType:  resourceType,
		ResourceID:    resourceID,
		Timestamp:     entry.Timestamp,
		Metadata:      metadata,
		Hash:          entry.Hash,
		PreviousHash:  previousHash,
		SchemaVersion: entry.SchemaVersion,
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

	// Normalize address for consistent DB lookup
	addr, err := types.NewWalletAddress(actor)
	if err != nil {
		return nil, fmt.Errorf("invalid actor address: %w", err)
	}

	return s.repo.List(ctx, addr.String(), limit)
}

// VerifyIntegrity checks the per-actor hash chain for tampering.
func (s *service) VerifyIntegrity(ctx context.Context, actor string) (bool, error) {
	addr, err := types.NewWalletAddress(actor)
	if err != nil {
		return false, fmt.Errorf("invalid actor address: %w", err)
	}
	actor = addr.String()

	entries, err := s.repo.List(ctx, actor, 0)
	if err != nil {
		return false, err
	}
	if len(entries) == 0 {
		return true, nil
	}

	// Repo returns newest-first (timestamp DESC, id DESC).
	for i := 0; i < len(entries); i++ {
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

		// Oldest entry should link to genesis.
		if i == len(entries)-1 {
			if e.PreviousHash != "GENESIS" {
				slog.ErrorContext(ctx, "audit integrity failure: invalid genesis link", "id", e.ID, "previous_hash", e.PreviousHash)
				return false, nil
			}
			continue
		}

		nextOlder := entries[i+1]
		if e.PreviousHash != nextOlder.Hash {
			slog.ErrorContext(ctx, "audit integrity failure: chain broken", "id", e.ID, "previous_hash", e.PreviousHash, "expected_previous_hash", nextOlder.Hash)
			return false, nil
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

	entries, _, err := s.repo.Query(ctx, filter)
	return entries, err
}

func (s *service) GetEntryByID(ctx context.Context, id string) (*AuditEntry, error) {
	return s.repo.GetByID(ctx, types.ID(id))
}

func (s *service) GetEntriesByResource(ctx context.Context, resourceID string) ([]AuditEntry, error) {
	return s.repo.GetByResource(ctx, types.ID(resourceID))
}

func (s *service) QueryEntries(ctx context.Context, filter audit.QueryFilter) (*audit.QueryResult, error) {
	entries, total, err := s.repo.Query(ctx, filter)
	if err != nil {
		return nil, err
	}

	// Convert DB entries to Protocol entries
	protocolEntries := make([]audit.Entry, len(entries))
	for i, e := range entries {
		protocolEntries[i] = audit.Entry{
			Actor:        types.WalletAddress(e.Actor),
			Action:       e.Action,
			ResourceType: e.ResourceType,
			ResourceID:   types.ID(e.ResourceID),
			Timestamp:    e.Timestamp,
			Hash:         e.Hash,
			PreviousHash: e.PreviousHash,
			Metadata:     types.Metadata(e.Metadata),
		}
	}

	return &audit.QueryResult{
		Entries:    protocolEntries,
		TotalCount: total,
	}, nil
}

func (s *service) BuildMerkleTree(ctx context.Context, actor string, startTime time.Time, endTime time.Time) (*AuditBatch, *audit.MerkleTree, error) {
	// Normalize actor for consistent DB lookup
	addr, err := types.NewWalletAddress(actor)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid actor address: %w", err)
	}
	actor = addr.String()

	entries, err := s.GetEntriesForMerkle(ctx, actor, startTime, endTime)
	if err != nil {
		return nil, nil, fmt.Errorf("build merkle tree: %w", err)
	}
	if len(entries) == 0 {
		return nil, nil, fmt.Errorf("build merkle tree: no entries in range")
	}

	// Sort oldest-first for deterministic Merkle leaves and to validate chain linkage within the range.
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Timestamp.Equal(entries[j].Timestamp) {
			return entries[i].ID < entries[j].ID
		}
		return entries[i].Timestamp.Before(entries[j].Timestamp)
	})

	// Harden: validate entry hashes (and chain linkage where possible) before persisting/anchoring.
	// NOTE: For time-windowed queries, we cannot validate the oldest entry's PreviousHash unless
	// startTime is unset (meaning we include the full chain from genesis).
	for i := range entries {
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
			return nil, nil, fmt.Errorf("build merkle tree: integrity failure (hash mismatch) entry=%s", e.ID)
		}

		if i == 0 {
			if startTime.IsZero() && e.PreviousHash != "GENESIS" {
				return nil, nil, fmt.Errorf("build merkle tree: integrity failure (invalid genesis link) entry=%s", e.ID)
			}
			continue
		}

		prev := entries[i-1]
		if e.PreviousHash != prev.Hash {
			return nil, nil, fmt.Errorf("build merkle tree: integrity failure (chain broken) entry=%s", e.ID)
		}
	}

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
	// Normalize actor for consistent DB lookup
	addr, err := types.NewWalletAddress(actor)
	if err != nil {
		return nil, fmt.Errorf("invalid actor address: %w", err)
	}
	actor = addr.String()

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
	// Normalize actor for consistent DB lookup
	addr, err := types.NewWalletAddress(actor)
	if err != nil {
		return nil, fmt.Errorf("invalid actor address: %w", err)
	}
	actor = addr.String()

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
	// Normalize actor for consistent DB lookup
	addr, err := types.NewWalletAddress(actor)
	if err != nil {
		return nil, fmt.Errorf("invalid actor address: %w", err)
	}
	actor = addr.String()

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
