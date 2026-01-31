package audit

import (
	"context"
	"testing"
	"time"

	protocol "github.com/itspablomontes/fleming/pkg/protocol/audit"
	"github.com/itspablomontes/fleming/pkg/protocol/types"
)

type mockRepo struct {
	entries []AuditEntry
	batches []AuditBatch
}

func (m *mockRepo) Create(ctx context.Context, entry *AuditEntry) error {
	m.entries = append(m.entries, *entry)
	return nil
}

func (m *mockRepo) GetLatest(ctx context.Context) (*AuditEntry, error) {
	if len(m.entries) == 0 {
		return nil, nil
	}
	return &m.entries[len(m.entries)-1], nil
}

func (m *mockRepo) List(ctx context.Context, actor string, limit int) ([]AuditEntry, error) {
	return m.entries, nil
}

func (m *mockRepo) GetByResource(ctx context.Context, resourceID types.ID) ([]AuditEntry, error) {
	var result []AuditEntry
	for _, entry := range m.entries {
		if entry.ResourceID == resourceID.String() {
			result = append(result, entry)
		}
	}
	return result, nil
}

func (m *mockRepo) GetByActor(ctx context.Context, actor types.WalletAddress) ([]AuditEntry, error) {
	var result []AuditEntry
	for _, entry := range m.entries {
		if entry.Actor == actor.String() {
			result = append(result, entry)
		}
	}
	return result, nil
}

func (m *mockRepo) GetByID(ctx context.Context, id types.ID) (*AuditEntry, error) {
	for _, entry := range m.entries {
		if entry.ID == id.String() {
			found := entry
			return &found, nil
		}
	}
	return nil, nil
}

func (m *mockRepo) Query(ctx context.Context, filter protocol.QueryFilter) ([]AuditEntry, error) {
	var result []AuditEntry
	for _, entry := range m.entries {
		if !filter.Actor.IsEmpty() && entry.Actor != filter.Actor.String() {
			continue
		}
		if filter.StartTime != nil && entry.Timestamp.Before(filter.StartTime.Time) {
			continue
		}
		if filter.EndTime != nil && entry.Timestamp.After(filter.EndTime.Time) {
			continue
		}
		result = append(result, entry)
	}
	return result, nil
}

func (m *mockRepo) CreateBatch(ctx context.Context, batch *AuditBatch) error {
	if batch.ID == "" {
		batch.ID = "batch-1"
	}
	m.batches = append(m.batches, *batch)
	return nil
}

func (m *mockRepo) UpdateBatch(ctx context.Context, batch *AuditBatch) error {
	if batch == nil {
		return nil
	}
	for i := range m.batches {
		if m.batches[i].ID == batch.ID {
			m.batches[i] = *batch
			return nil
		}
	}
	m.batches = append(m.batches, *batch)
	return nil
}

func (m *mockRepo) GetBatchByIDForActor(ctx context.Context, actor string, id string) (*AuditBatch, error) {
	for _, batch := range m.batches {
		if batch.ID == id && batch.Actor == actor {
			found := batch
			return &found, nil
		}
	}
	return nil, nil
}

func (m *mockRepo) GetBatchByActorAndRoot(ctx context.Context, actor string, rootHash string) (*AuditBatch, error) {
	for _, batch := range m.batches {
		if batch.Actor == actor && batch.RootHash == rootHash {
			found := batch
			return &found, nil
		}
	}
	return nil, nil
}

func (m *mockRepo) ListBatchesByActor(ctx context.Context, actor string, limit int, offset int) ([]AuditBatch, error) {
	var out []AuditBatch
	for _, b := range m.batches {
		if b.Actor == actor {
			out = append(out, b)
		}
	}
	return out, nil
}

func (m *mockRepo) GetDistinctActorsWithEntries(ctx context.Context, startTime time.Time, endTime time.Time, limit int) ([]string, error) {
	seen := map[string]bool{}
	var actors []string
	for _, e := range m.entries {
		if !startTime.IsZero() && e.Timestamp.Before(startTime) {
			continue
		}
		if !endTime.IsZero() && e.Timestamp.After(endTime) {
			continue
		}
		if seen[e.Actor] {
			continue
		}
		seen[e.Actor] = true
		actors = append(actors, e.Actor)
		if limit > 0 && len(actors) >= limit {
			break
		}
	}
	return actors, nil
}

func TestService_BuildMerkleTreeAndVerifyProof(t *testing.T) {
	actor := "0x1234567890abcdef1234567890abcdef12345678"
	repo := &mockRepo{
		entries: []AuditEntry{
			{
				ID:        "entry-1",
				Actor:     actor,
				Hash:      "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
				Timestamp: time.Date(2026, 1, 25, 10, 0, 0, 0, time.UTC),
			},
			{
				ID:        "entry-2",
				Actor:     actor,
				Hash:      "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
				Timestamp: time.Date(2026, 1, 25, 11, 0, 0, 0, time.UTC),
			},
		},
	}
	service := NewService(repo)

	batch, tree, err := service.BuildMerkleTree(context.Background(), actor, time.Time{}, time.Time{})
	if err != nil {
		t.Fatalf("BuildMerkleTree() error = %v", err)
	}
	if batch == nil || tree == nil {
		t.Fatal("expected batch and tree to be returned")
	}
	if batch.EntryCount != 2 {
		t.Fatalf("expected entry count 2, got %d", batch.EntryCount)
	}
	if batch.RootHash != tree.Root {
		t.Fatalf("batch root mismatch: got %s want %s", batch.RootHash, tree.Root)
	}

	proof, err := protocol.GenerateProof(tree, repo.entries[0].Hash)
	if err != nil {
		t.Fatalf("GenerateProof() error = %v", err)
	}
	if !service.VerifyMerkleProof(tree.Root, repo.entries[0].Hash, proof) {
		t.Fatal("VerifyMerkleProof() expected true")
	}

	fetched, err := service.GetBatch(context.Background(), actor, batch.ID)
	if err != nil {
		t.Fatalf("GetBatch() error = %v", err)
	}
	if fetched == nil {
		t.Fatal("expected batch to be returned")
	}
	if fetched.RootHash != tree.Root {
		t.Fatalf("GetBatch() root mismatch: got %s want %s", fetched.RootHash, tree.Root)
	}

	byRoot, err := service.GetBatchByRoot(context.Background(), actor, tree.Root)
	if err != nil {
		t.Fatalf("GetBatchByRoot() error = %v", err)
	}
	if byRoot == nil {
		t.Fatal("expected GetBatchByRoot() to return a batch")
	}
	if byRoot.ID != batch.ID {
		t.Fatalf("expected GetBatchByRoot() id %q, got %q", batch.ID, byRoot.ID)
	}
}
