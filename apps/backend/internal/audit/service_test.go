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

func (m *mockRepo) GetBatchByID(ctx context.Context, id string) (*AuditBatch, error) {
	for _, batch := range m.batches {
		if batch.ID == id {
			found := batch
			return &found, nil
		}
	}
	return nil, nil
}

func (m *mockRepo) GetBatchByRoot(ctx context.Context, rootHash string) (*AuditBatch, error) {
	for _, batch := range m.batches {
		if batch.RootHash == rootHash {
			found := batch
			return &found, nil
		}
	}
	return nil, nil
}

func TestService_BuildMerkleTreeAndVerifyProof(t *testing.T) {
	repo := &mockRepo{
		entries: []AuditEntry{
			{
				ID:        "entry-1",
				Hash:      "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
				Timestamp: time.Date(2026, 1, 25, 10, 0, 0, 0, time.UTC),
			},
			{
				ID:        "entry-2",
				Hash:      "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
				Timestamp: time.Date(2026, 1, 25, 11, 0, 0, 0, time.UTC),
			},
		},
	}
	service := NewService(repo)

	batch, tree, err := service.BuildMerkleTree(context.Background(), time.Time{}, time.Time{})
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

	root, err := service.GetMerkleRoot(context.Background(), batch.ID)
	if err != nil {
		t.Fatalf("GetMerkleRoot() error = %v", err)
	}
	if root != tree.Root {
		t.Fatalf("GetMerkleRoot() mismatch: got %s want %s", root, tree.Root)
	}
}
