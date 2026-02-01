package audit

import (
	"context"
	"sort"
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

func (m *mockRepo) GetLatest(ctx context.Context, actor string) (*AuditEntry, error) {
	if len(m.entries) == 0 {
		return nil, nil
	}
	// Entries are appended in creation order in tests; scan from the end.
	for i := len(m.entries) - 1; i >= 0; i-- {
		if actor == "" || m.entries[i].Actor == actor {
			e := m.entries[i]
			return &e, nil
		}
	}
	return nil, nil
}

func (m *mockRepo) List(ctx context.Context, actor string, limit int) ([]AuditEntry, error) {
	var out []AuditEntry
	for _, e := range m.entries {
		if actor == "" || e.Actor == actor {
			out = append(out, e)
		}
	}
	// Match production ordering: newest-first.
	sort.Slice(out, func(i, j int) bool {
		if out[i].Timestamp.Equal(out[j].Timestamp) {
			return out[i].ID > out[j].ID
		}
		return out[i].Timestamp.After(out[j].Timestamp)
	})
	if limit > 0 && len(out) > limit {
		out = out[:limit]
	}
	return out, nil
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

func (m *mockRepo) Query(ctx context.Context, filter protocol.QueryFilter) ([]AuditEntry, int64, error) {
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
	return result, int64(len(result)), nil
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
	ts1 := time.Date(2026, 1, 25, 10, 0, 0, 0, time.UTC)
	ts2 := time.Date(2026, 1, 25, 11, 0, 0, 0, time.UTC)

	e1 := protocol.Entry{
		Actor:        types.WalletAddress(actor),
		Action:       protocol.ActionCreate,
		ResourceType: protocol.ResourceEvent,
		ResourceID:   types.ID("res-1"),
		Timestamp:    ts1,
		PreviousHash: "GENESIS",
	}
	h1 := e1.ComputeHash()

	e2 := protocol.Entry{
		Actor:        types.WalletAddress(actor),
		Action:       protocol.ActionUpdate,
		ResourceType: protocol.ResourceEvent,
		ResourceID:   types.ID("res-1"),
		Timestamp:    ts2,
		PreviousHash: h1,
	}
	h2 := e2.ComputeHash()

	repo := &mockRepo{
		entries: []AuditEntry{
			{
				ID:           "entry-1",
				Actor:        actor,
				Action:       e1.Action,
				ResourceType: e1.ResourceType,
				ResourceID:   e1.ResourceID.String(),
				Timestamp:    ts1,
				Hash:         h1,
				PreviousHash: e1.PreviousHash,
			},
			{
				ID:           "entry-2",
				Actor:        actor,
				Action:       e2.Action,
				ResourceType: e2.ResourceType,
				ResourceID:   e2.ResourceID.String(),
				Timestamp:    ts2,
				Hash:         h2,
				PreviousHash: e2.PreviousHash,
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

	proof, err := protocol.GenerateProof(tree, h1)
	if err != nil {
		t.Fatalf("GenerateProof() error = %v", err)
	}
	if !service.VerifyMerkleProof(tree.Root, h1, proof) {
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

func TestService_Record_ChainsPerActor(t *testing.T) {
	repo := &mockRepo{}
	service := NewService(repo)

	a := "0x1234567890abcdef1234567890abcdef12345678"
	b := "0x0000000000000000000000000000000000000abc"

	if err := service.Record(context.Background(), a, protocol.ActionCreate, protocol.ResourceEvent, "res-1", nil); err != nil {
		t.Fatalf("Record(a) error = %v", err)
	}
	if err := service.Record(context.Background(), b, protocol.ActionCreate, protocol.ResourceEvent, "res-2", nil); err != nil {
		t.Fatalf("Record(b) error = %v", err)
	}

	var lastA, lastB *AuditEntry
	for i := range repo.entries {
		e := repo.entries[i]
		if e.Actor == a {
			cpy := e
			lastA = &cpy
		}
		if e.Actor == b {
			cpy := e
			lastB = &cpy
		}
	}
	if lastA == nil || lastB == nil {
		t.Fatalf("expected entries for both actors")
	}
	if lastA.PreviousHash != "GENESIS" {
		t.Fatalf("expected actor A genesis link, got %q", lastA.PreviousHash)
	}
	if lastB.PreviousHash != "GENESIS" {
		t.Fatalf("expected actor B genesis link, got %q", lastB.PreviousHash)
	}
}

func TestService_VerifyIntegrity_PerActor(t *testing.T) {
	repo := &mockRepo{}
	service := NewService(repo)

	actor := "0x1234567890abcdef1234567890abcdef12345678"
	if err := service.Record(context.Background(), actor, protocol.ActionCreate, protocol.ResourceEvent, "res-1", nil); err != nil {
		t.Fatalf("Record() error = %v", err)
	}
	if err := service.Record(context.Background(), actor, protocol.ActionUpdate, protocol.ResourceEvent, "res-1", nil); err != nil {
		t.Fatalf("Record() error = %v", err)
	}

	ok, err := service.VerifyIntegrity(context.Background(), actor)
	if err != nil {
		t.Fatalf("VerifyIntegrity() error = %v", err)
	}
	if !ok {
		t.Fatalf("expected VerifyIntegrity() true")
	}
}

func TestService_Record_MonotonicTimestampPerActor(t *testing.T) {
	actor := "0x1234567890abcdef1234567890abcdef12345678"
	now := time.Now().UTC().Truncate(time.Microsecond)
	repo := &mockRepo{
		entries: []AuditEntry{
			{
				ID:           "entry-0",
				Actor:        actor,
				Action:       protocol.ActionCreate,
				ResourceType: protocol.ResourceEvent,
				ResourceID:   "res-0",
				Timestamp:    now.Add(1 * time.Second), // force bump path
				Hash:         "deadbeef",               // unused by Record()
				PreviousHash: "GENESIS",
			},
		},
	}
	service := NewService(repo)

	if err := service.Record(context.Background(), actor, protocol.ActionCreate, protocol.ResourceEvent, "res-1", nil); err != nil {
		t.Fatalf("Record() error = %v", err)
	}

	created := repo.entries[len(repo.entries)-1]
	if !created.Timestamp.After(repo.entries[0].Timestamp) {
		t.Fatalf("expected monotonic timestamp: got %v, want > %v", created.Timestamp, repo.entries[0].Timestamp)
	}
}
