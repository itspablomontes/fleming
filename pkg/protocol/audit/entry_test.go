package audit

import (
	"testing"
	"time"

	"github.com/itspablomontes/fleming/pkg/protocol/types"
)

func TestAction_IsValid(t *testing.T) {
	tests := []struct {
		action Action
		want   bool
	}{
		// CRUD
		{ActionCreate, true},
		{ActionRead, true},
		{ActionUpdate, true},
		{ActionDelete, true},
		// Consent
		{ActionConsentRequest, true},
		{ActionConsentApprove, true},
		{ActionConsentDeny, true},
		{ActionConsentRevoke, true},
		{ActionConsentExpire, true},
		{ActionConsentSuspend, true},
		{ActionConsentResume, true},
		// Auth
		{ActionLogin, true},
		{ActionLogout, true},
		// Files
		{ActionUpload, true},
		{ActionDownload, true},
		{ActionShare, true},
		// VC
		{ActionVCIssue, true},
		{ActionVCRevoke, true},
		{ActionVCVerify, true},
		{ActionVCPresent, true},
		// ZK
		{ActionZKGenerate, true},
		{ActionZKVerify, true},
		// Attestation
		{ActionCosign, true},
		{ActionAttest, true},
		// Invalid
		{"unknown", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(string(tt.action), func(t *testing.T) {
			if got := tt.action.IsValid(); got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResourceType_IsValid(t *testing.T) {
	tests := []struct {
		rt   ResourceType
		want bool
	}{
		{ResourceEvent, true},
		{ResourceFile, true},
		{ResourceConsent, true},
		{ResourceSession, true},
		{ResourceVC, true},
		{ResourceZKProof, true},
		{ResourceAttestation, true},
		{"unknown", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(string(tt.rt), func(t *testing.T) {
			if got := tt.rt.IsValid(); got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEntry_Validate(t *testing.T) {
	validActor, _ := types.NewWalletAddress("0x1234567890abcdef1234567890abcdef12345678")

	tests := []struct {
		name    string
		entry   Entry
		wantErr bool
	}{
		{
			name: "valid entry",
			entry: Entry{
				Actor:        validActor,
				Action:       ActionCreate,
				ResourceType: ResourceEvent,
				ResourceID:   "event-1",
				Timestamp:    time.Now(),
			},
			wantErr: false,
		},
		{
			name: "missing actor",
			entry: Entry{
				Action:       ActionCreate,
				ResourceType: ResourceEvent,
				ResourceID:   "event-1",
				Timestamp:    time.Now(),
			},
			wantErr: true,
		},
		{
			name: "invalid action",
			entry: Entry{
				Actor:        validActor,
				Action:       "invalid",
				ResourceType: ResourceEvent,
				ResourceID:   "event-1",
				Timestamp:    time.Now(),
			},
			wantErr: true,
		},
		{
			name: "missing resource ID",
			entry: Entry{
				Actor:        validActor,
				Action:       ActionCreate,
				ResourceType: ResourceEvent,
				Timestamp:    time.Now(),
			},
			wantErr: true,
		},
		{
			name: "missing timestamp",
			entry: Entry{
				Actor:        validActor,
				Action:       ActionCreate,
				ResourceType: ResourceEvent,
				ResourceID:   "event-1",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.entry.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEntry_Hash(t *testing.T) {
	actor, _ := types.NewWalletAddress("0x1234567890abcdef1234567890abcdef12345678")
	timestamp := time.Date(2026, 1, 23, 12, 0, 0, 0, time.UTC)

	entry := &Entry{
		Actor:        actor,
		Action:       ActionCreate,
		ResourceType: ResourceEvent,
		ResourceID:   "event-1",
		Timestamp:    timestamp,
	}

	entry.SetHash()

	if entry.Hash == "" {
		t.Error("Hash should not be empty after SetHash()")
	}

	if !entry.VerifyHash() {
		t.Error("Hash verification should pass for unmodified entry")
	}

	originalHash := entry.Hash
	entry.ResourceID = "tampered-event"

	if entry.VerifyHash() {
		t.Error("Hash verification should fail for tampered entry")
	}

	entry.ResourceID = "event-1"
	if entry.ComputeHash() != originalHash {
		t.Error("Same inputs should produce same hash")
	}
}

func TestNewEntry(t *testing.T) {
	actor, _ := types.NewWalletAddress("0x1234567890abcdef1234567890abcdef12345678")

	entry := NewEntry(actor, ActionCreate, ResourceEvent, "event-1", "")

	if entry.Actor.IsEmpty() {
		t.Error("Actor should be set")
	}

	if entry.Hash == "" {
		t.Error("Hash should be auto-computed")
	}

	if entry.Timestamp.IsZero() {
		t.Error("Timestamp should be auto-set")
	}
}

func TestEntry_ChainIntegrity(t *testing.T) {
	actor, _ := types.NewWalletAddress("0x1234567890abcdef1234567890abcdef12345678")

	entry1 := NewEntry(actor, ActionCreate, ResourceEvent, "event-1", "")
	entry2 := NewEntry(actor, ActionUpdate, ResourceEvent, "event-1", entry1.Hash)
	entry3 := NewEntry(actor, ActionRead, ResourceEvent, "event-1", entry2.Hash)

	if entry2.PreviousHash != entry1.Hash {
		t.Error("Entry 2 should link to Entry 1")
	}

	if entry3.PreviousHash != entry2.Hash {
		t.Error("Entry 3 should link to Entry 2")
	}

	if !entry1.VerifyHash() || !entry2.VerifyHash() || !entry3.VerifyHash() {
		t.Error("All entries in chain should verify")
	}
}

func TestQueryFilter(t *testing.T) {
	actor, _ := types.NewWalletAddress("0x1234567890abcdef1234567890abcdef12345678")

	filter := NewQueryFilter().
		WithActor(actor).
		WithResource("event-1").
		WithAction(ActionRead).
		WithLimit(50)

	if filter.Actor != actor {
		t.Error("Actor filter not set")
	}

	if filter.ResourceID != "event-1" {
		t.Error("Resource filter not set")
	}

	if filter.Action != ActionRead {
		t.Error("Action filter not set")
	}

	if filter.Limit != 50 {
		t.Error("Limit not set")
	}
}
