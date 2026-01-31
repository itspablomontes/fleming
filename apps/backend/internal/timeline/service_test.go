package timeline

import (
	"context"
	"fmt"
	"testing"
	"time"

	"io"

	"github.com/itspablomontes/fleming/apps/backend/internal/audit"
	"github.com/itspablomontes/fleming/apps/backend/internal/common"
	"github.com/itspablomontes/fleming/apps/backend/internal/storage"
	protocol "github.com/itspablomontes/fleming/pkg/protocol/audit"
	"github.com/itspablomontes/fleming/pkg/protocol/timeline"
	"github.com/itspablomontes/fleming/pkg/protocol/types"
)

type MockAuditService struct{}

func (m *MockAuditService) Record(ctx context.Context, actor string, action protocol.Action, resourceType protocol.ResourceType, resourceID string, metadata common.JSONMap) error {
	return nil
}
func (m *MockAuditService) GetLatestEntries(ctx context.Context, actor string, limit int) ([]audit.AuditEntry, error) {
	return nil, nil
}
func (m *MockAuditService) VerifyIntegrity(ctx context.Context) (bool, error) {
	return true, nil
}
func (m *MockAuditService) BuildMerkleTree(ctx context.Context, startTime time.Time, endTime time.Time) (*audit.AuditBatch, *protocol.MerkleTree, error) {
	return nil, nil, nil
}
func (m *MockAuditService) GetMerkleRoot(ctx context.Context, batchID string) (string, error) {
	return "", nil
}
func (m *MockAuditService) VerifyMerkleProof(root string, entryHash string, proof *protocol.Proof) bool {
	return true
}
func (m *MockAuditService) GetEntriesForMerkle(ctx context.Context, startTime time.Time, endTime time.Time) ([]audit.AuditEntry, error) {
	return nil, nil
}
func (m *MockAuditService) GetEntryByID(ctx context.Context, id string) (*audit.AuditEntry, error) {
	return nil, nil
}
func (m *MockAuditService) GetEntriesByResource(ctx context.Context, resourceID string) ([]audit.AuditEntry, error) {
	return nil, nil
}
func (m *MockAuditService) QueryEntries(ctx context.Context, filter protocol.QueryFilter) ([]audit.AuditEntry, error) {
	return nil, nil
}

type MockStorage struct{}

func (m *MockStorage) Put(ctx context.Context, bucketName, objectName string, reader io.Reader, objectSize int64, contentType string) (string, error) {
	return objectName, nil
}
func (m *MockStorage) Get(ctx context.Context, bucketName, objectName string) (io.ReadCloser, error) {
	return nil, nil
}
func (m *MockStorage) Delete(ctx context.Context, bucketName, objectName string) error {
	return nil
}
func (m *MockStorage) GetURL(ctx context.Context, bucketName, objectName string) (string, error) {
	return "http://localhost:9000/" + objectName, nil
}
func (m *MockStorage) CreateMultipartUpload(ctx context.Context, bucketName, objectName, contentType string) (string, error) {
	return "upload-id", nil
}
func (m *MockStorage) UploadPart(ctx context.Context, bucketName, objectName, uploadID string, partNumber int, reader io.Reader, objectSize int64) (string, error) {
	return "etag", nil
}
func (m *MockStorage) CompleteMultipartUpload(ctx context.Context, bucketName, objectName, uploadID string, parts []storage.Part) (string, error) {
	return objectName, nil
}
func (m *MockStorage) AbortMultipartUpload(ctx context.Context, bucketName, objectName, uploadID string) error {
	return nil
}

type MockRepo struct {
	nextID int
	events []timeline.Event
	edges  []timeline.Edge
}

func (m *MockRepo) GetEvent(ctx context.Context, id types.ID) (*timeline.Event, error) {
	for i := range m.events {
		if m.events[i].ID == id {
			evt := m.events[i]
			return &evt, nil
		}
	}
	return nil, nil
}

func (m *MockRepo) GetTimeline(ctx context.Context, patientID types.WalletAddress) ([]timeline.Event, error) {
	out := make([]timeline.Event, 0, len(m.events))
	for _, e := range m.events {
		if e.PatientID == patientID {
			out = append(out, e)
		}
	}
	return out, nil
}

func (m *MockRepo) GetRelated(ctx context.Context, eventID types.ID, depth int) ([]timeline.Event, []timeline.Edge, error) {
	// Tests don't require graph traversal; return no related edges so filtering doesn't remove events.
	return []timeline.Event{}, []timeline.Edge{}, nil
}

func (m *MockRepo) CreateEvent(ctx context.Context, event *timeline.Event) error {
	if event == nil {
		return fmt.Errorf("event is nil")
	}
	if event.ID.IsEmpty() {
		m.nextID++
		event.ID = types.ID(fmt.Sprintf("evt-%d", m.nextID))
	}
	m.events = append(m.events, *event)
	return nil
}

func (m *MockRepo) UpdateEvent(ctx context.Context, event *timeline.Event) error {
	if event == nil {
		return fmt.Errorf("event is nil")
	}
	for i := range m.events {
		if m.events[i].ID == event.ID {
			m.events[i] = *event
			return nil
		}
	}
	m.events = append(m.events, *event)
	return nil
}

func (m *MockRepo) DeleteEvent(ctx context.Context, id types.ID) error {
	for i := range m.events {
		if m.events[i].ID == id {
			m.events = append(m.events[:i], m.events[i+1:]...)
			return nil
		}
	}
	return nil
}

func (m *MockRepo) CreateEdge(ctx context.Context, edge *timeline.Edge) error {
	if edge == nil {
		return fmt.Errorf("edge is nil")
	}
	if edge.ID.IsEmpty() {
		m.nextID++
		edge.ID = types.ID(fmt.Sprintf("edge-%d", m.nextID))
	}
	m.edges = append(m.edges, *edge)
	return nil
}

func (m *MockRepo) DeleteEdge(ctx context.Context, id types.ID) error {
	for i := range m.edges {
		if m.edges[i].ID == id {
			m.edges = append(m.edges[:i], m.edges[i+1:]...)
			return nil
		}
	}
	return nil
}

func (m *MockRepo) CreateFile(ctx context.Context, file *EventFile) error { return nil }
func (m *MockRepo) GetFileByID(ctx context.Context, id string) (*EventFile, error) {
	return nil, nil
}
func (m *MockRepo) GetFilesByEventID(ctx context.Context, eventID string) ([]EventFile, error) {
	return nil, nil
}
func (m *MockRepo) UpsertFileAccess(ctx context.Context, confirmations *EventFileAccess) error { return nil }
func (m *MockRepo) GetFileAccess(ctx context.Context, fileID string, grantee string) (*EventFileAccess, error) {
	return nil, nil
}
func (m *MockRepo) GetGraphData(ctx context.Context, patientID string) ([]TimelineEvent, []EventEdge, error) {
	return []TimelineEvent{}, []EventEdge{}, nil
}
func (m *MockRepo) Transaction(ctx context.Context, fn func(repo Repository) error) error { return fn(m) }

func TestService_CreateEvent(t *testing.T) {
	repo := &MockRepo{}
	auditSvc := &MockAuditService{}
	storageSvc := &MockStorage{}
	svc := NewService(repo, auditSvc, storageSvc, "test-bucket")

	patientID, err := types.NewWalletAddress("0x0000000000000000000000000000000000000123")
	if err != nil {
		t.Fatalf("unexpected patient id error: %v", err)
	}

	event, err := timeline.NewEventBuilder().
		WithPatientID(patientID).
		WithType(timeline.EventLabResult).
		WithTitle("Blood Test").
		WithTimestamp(time.Now()).
		Build()
	if err != nil {
		t.Fatalf("unexpected event build error: %v", err)
	}

	if err := svc.CreateEvent(context.Background(), event); err != nil {
		t.Fatalf("CreateEvent() error = %v", err)
	}

	got, err := svc.GetTimelineForPatient(context.Background(), patientID)
	if err != nil {
		t.Fatalf("GetTimelineForPatient() error = %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("GetTimelineForPatient() count = %d, want %d", len(got), 1)
	}
}

func TestService_GetTimelineForPatient_FiltersByPatient(t *testing.T) {
	repo := &MockRepo{}
	auditSvc := &MockAuditService{}
	storageSvc := &MockStorage{}
	svc := NewService(repo, auditSvc, storageSvc, "test-bucket")

	p1, _ := types.NewWalletAddress("0x0000000000000000000000000000000000000123")
	p2, _ := types.NewWalletAddress("0x0000000000000000000000000000000000000456")

	e1, _ := timeline.NewEventBuilder().WithPatientID(p1).WithType(timeline.EventConsultation).WithTitle("Event 1").WithTimestamp(time.Now()).Build()
	e2, _ := timeline.NewEventBuilder().WithPatientID(p2).WithType(timeline.EventConsultation).WithTitle("Event 2").WithTimestamp(time.Now()).Build()
	_ = svc.CreateEvent(context.Background(), e1)
	_ = svc.CreateEvent(context.Background(), e2)

	got, err := svc.GetTimelineForPatient(context.Background(), p1)
	if err != nil {
		t.Fatalf("GetTimelineForPatient() error = %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("GetTimelineForPatient() count = %d, want %d", len(got), 1)
	}
}
