package audit

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	protocolaudit "github.com/itspablomontes/fleming/pkg/protocol/audit"
	protocolchain "github.com/itspablomontes/fleming/pkg/protocol/chain"
	"github.com/itspablomontes/fleming/pkg/protocol/types"
)

type memRepo struct {
	batches map[string]*AuditBatch
}

func (m *memRepo) Create(ctx context.Context, entry *AuditEntry) error { return nil }
func (m *memRepo) GetLatest(ctx context.Context) (*AuditEntry, error)  { return nil, nil }
func (m *memRepo) List(ctx context.Context, actor string, limit int) ([]AuditEntry, error) {
	return nil, nil
}
func (m *memRepo) GetByResource(ctx context.Context, resourceID types.ID) ([]AuditEntry, error) {
	return nil, nil
}
func (m *memRepo) GetByActor(ctx context.Context, actor types.WalletAddress) ([]AuditEntry, error) {
	return nil, nil
}
func (m *memRepo) GetByID(ctx context.Context, id types.ID) (*AuditEntry, error) { return nil, nil }
func (m *memRepo) Query(ctx context.Context, filter protocolaudit.QueryFilter) ([]AuditEntry, error) {
	return nil, nil
}

func (m *memRepo) CreateBatch(ctx context.Context, batch *AuditBatch) error {
	if m.batches == nil {
		m.batches = make(map[string]*AuditBatch)
	}
	m.batches[batch.ID] = batch
	return nil
}

func (m *memRepo) UpdateBatch(ctx context.Context, batch *AuditBatch) error {
	if m.batches == nil {
		m.batches = make(map[string]*AuditBatch)
	}
	cpy := *batch
	m.batches[batch.ID] = &cpy
	return nil
}

func (m *memRepo) GetBatchByIDForActor(ctx context.Context, actor string, id string) (*AuditBatch, error) {
	if m.batches == nil {
		return nil, nil
	}
	b := m.batches[id]
	if b == nil || b.Actor != actor {
		return nil, nil
	}
	cpy := *b
	return &cpy, nil
}

func (m *memRepo) GetBatchByActorAndRoot(ctx context.Context, actor string, rootHash string) (*AuditBatch, error) {
	for _, b := range m.batches {
		if b != nil && b.Actor == actor && b.RootHash == rootHash {
			cpy := *b
			return &cpy, nil
		}
	}
	return nil, nil
}

func (m *memRepo) ListBatchesByActor(ctx context.Context, actor string, limit int, offset int) ([]AuditBatch, error) {
	return nil, nil
}

func (m *memRepo) GetDistinctActorsWithEntries(ctx context.Context, startTime time.Time, endTime time.Time, limit int) ([]string, error) {
	return nil, nil
}

type mockChainClient struct {
	anchorCalls int
	verifyCalls int
	anchorRes   *protocolchain.AnchorResult
	anchorErr   error
	verifyTs    uint64
	verifyErr   error
}

func (m *mockChainClient) AnchorRoot(ctx context.Context, hexRoot string) (*protocolchain.AnchorResult, error) {
	m.anchorCalls++
	return m.anchorRes, m.anchorErr
}
func (m *mockChainClient) VerifyRoot(ctx context.Context, hexRoot string) (uint64, error) {
	m.verifyCalls++
	return m.verifyTs, m.verifyErr
}
func (m *mockChainClient) FindRootAnchoredEvent(ctx context.Context, hexRoot string) (*protocolchain.RootAnchoredEvent, bool, error) {
	return nil, false, nil
}

func TestHandleAnchorMerkleBatch_NotConfigured_ReturnsNotImplemented(t *testing.T) {
	t.Setenv("ENV", "dev")
	gin.SetMode(gin.TestMode)

	h := NewHandler(nil, nil)

	r := gin.New()
	h.RegisterRoutes(r.Group(""))

	req := httptest.NewRequest(http.MethodPost, "/audit/merkle/batch-1/anchor", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotImplemented {
		t.Fatalf("expected %d, got %d (%s)", http.StatusNotImplemented, rec.Code, rec.Body.String())
	}
}

func TestHandleAnchorMerkleBatch_Success(t *testing.T) {
	t.Setenv("ENV", "dev")
	gin.SetMode(gin.TestMode)

	actor := "0x1234567890abcdef1234567890abcdef12345678"
	root := "0000000000000000000000000000000000000000000000000000000000000001"
	repo := &memRepo{
		batches: map[string]*AuditBatch{
			"batch-1": {
				ID:           "batch-1",
				Actor:        actor,
				RootHash:     root,
				StartTime:    time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
				EndTime:      time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC),
				EntryCount:   1,
				CreatedAt:    time.Date(2026, 1, 2, 0, 0, 1, 0, time.UTC),
				AnchorStatus: "pending",
			},
		},
	}
	svc := NewService(repo)
	chain := &mockChainClient{
		anchorRes: &protocolchain.AnchorResult{TxHash: "0xabc", BlockNumber: 123, GasUsed: 456},
		verifyTs:  1700000000,
	}

	h := NewHandler(svc, chain)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_address", actor)
		c.Next()
	})
	h.RegisterRoutes(r.Group(""))

	req := httptest.NewRequest(http.MethodPost, "/audit/merkle/batch-1/anchor", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d (%s)", http.StatusOK, rec.Code, rec.Body.String())
	}

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid json response: %v (%s)", err, rec.Body.String())
	}

	batch, ok := body["batch"].(map[string]any)
	if !ok {
		t.Fatalf("expected batch object, got %T (%v)", body["batch"], body["batch"])
	}

	if batch["rootHash"] != root {
		t.Fatalf("expected rootHash %q, got %v", root, batch["rootHash"])
	}
	if batch["anchorTxHash"] != "0xabc" {
		t.Fatalf("expected anchorTxHash %q, got %v", "0xabc", batch["anchorTxHash"])
	}
	if batch["anchorBlockNumber"] != float64(123) {
		t.Fatalf("expected anchorBlockNumber %d, got %v", 123, batch["anchorBlockNumber"])
	}
	if batch["anchorStatus"] != "anchored" {
		t.Fatalf("expected anchorStatus %q, got %v", "anchored", batch["anchorStatus"])
	}

	// Idempotency: second call should not send a second tx.
	req2 := httptest.NewRequest(http.MethodPost, "/audit/merkle/batch-1/anchor", nil)
	rec2 := httptest.NewRecorder()
	r.ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d (%s)", http.StatusOK, rec2.Code, rec2.Body.String())
	}
	if chain.anchorCalls != 1 {
		t.Fatalf("expected 1 AnchorRoot call, got %d", chain.anchorCalls)
	}
}
