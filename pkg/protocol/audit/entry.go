package audit

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/itspablomontes/fleming/pkg/protocol"
	"github.com/itspablomontes/fleming/pkg/protocol/types"
)

const SchemaVersionAudit = protocol.SchemaVersionAudit

type Action string

const (
	// CRUD operations
	ActionCreate Action = "create"
	ActionRead   Action = "read"
	ActionUpdate Action = "update"
	ActionDelete Action = "delete"

	// Consent operations
	ActionConsentRequest Action = "consent.request"
	ActionConsentApprove Action = "consent.approve"
	ActionConsentDeny    Action = "consent.deny"
	ActionConsentRevoke  Action = "consent.revoke"
	ActionConsentExpire  Action = "consent.expire"
	ActionConsentSuspend Action = "consent.suspend"
	ActionConsentResume  Action = "consent.resume"

	// Authentication
	ActionLogin  Action = "auth.login"
	ActionLogout Action = "auth.logout"

	// File operations
	ActionUpload   Action = "file.upload"
	ActionDownload Action = "file.download"
	ActionShare    Action = "file.share"

	// Verifiable Credentials
	ActionVCIssue   Action = "vc.issue"
	ActionVCRevoke  Action = "vc.revoke"
	ActionVCVerify  Action = "vc.verify"
	ActionVCPresent Action = "vc.present"

	// Zero-Knowledge Proofs
	ActionZKGenerate Action = "zk.generate"
	ActionZKVerify   Action = "zk.verify"

	// Attestation (Post-MVP)
	ActionCosign Action = "attestation.cosign"
	ActionAttest Action = "attestation.attest"
)

func (a Action) IsValid() bool {
	return GetActionRegistry().IsValid(a)
}

type ResourceType string

const (
	// Core resources
	ResourceEvent   ResourceType = "event"   // Timeline event
	ResourceFile    ResourceType = "file"    // File attachment
	ResourceConsent ResourceType = "consent" // Consent grant
	ResourceSession ResourceType = "session" // User session

	// Verifiable Credentials
	ResourceVC ResourceType = "vc" // Verifiable credential

	// Zero-Knowledge Proofs
	ResourceZKProof ResourceType = "zk_proof" // Zero-knowledge proof

	// Attestation
	ResourceAttestation ResourceType = "attestation" // Provider attestation
)

func (rt ResourceType) IsValid() bool {
	return GetResourceTypeRegistry().IsValid(rt)
}

type Entry struct {
	ID types.ID `json:"id"`

	Actor types.WalletAddress `json:"actor"`

	Action Action `json:"action"`

	ResourceType ResourceType `json:"resourceType"`

	ResourceID types.ID `json:"resourceId"`

	Timestamp time.Time `json:"timestamp"`

	Metadata types.Metadata `json:"metadata,omitempty"`

	SchemaVersion string `json:"schemaVersion,omitempty"`

	Hash string `json:"hash,omitempty"`

	PreviousHash string `json:"previousHash,omitempty"`
}

func (e *Entry) Validate() error {
	var errs types.ValidationErrors

	if e.Actor.IsEmpty() {
		errs.Add("actor", "actor is required")
	}

	if !e.Action.IsValid() {
		errs.Add("action", "invalid action")
	}

	if !e.ResourceType.IsValid() {
		errs.Add("resourceType", "invalid resource type")
	}

	if e.ResourceID.IsEmpty() {
		errs.Add("resourceId", "resource ID is required")
	}

	if e.Timestamp.IsZero() {
		errs.Add("timestamp", "timestamp is required")
	}

	if errs.HasErrors() {
		return errs
	}
	return nil
}

func (e *Entry) ComputeHash() string {
	data := struct {
		Actor        string       `json:"actor"`
		Action       Action       `json:"action"`
		ResourceType ResourceType `json:"resourceType"`
		ResourceID   string       `json:"resourceId"`
		Timestamp    string       `json:"timestamp"`
		PreviousHash string       `json:"previousHash"`
	}{
		Actor:        e.Actor.String(),
		Action:       e.Action,
		ResourceType: e.ResourceType,
		ResourceID:   e.ResourceID.String(),
		Timestamp:    e.Timestamp.UTC().Format(time.RFC3339Nano),
		PreviousHash: e.PreviousHash,
	}

	bytes, _ := json.Marshal(data)
	hash := sha256.Sum256(bytes)
	return hex.EncodeToString(hash[:])
}

func (e *Entry) SetHash() {
	e.Hash = e.ComputeHash()
}

func (e *Entry) VerifyHash() bool {
	if e.Hash == "" {
		return false
	}
	return e.Hash == e.ComputeHash()
}

func NewEntry(
	actor types.WalletAddress,
	action Action,
	resourceType ResourceType,
	resourceID types.ID,
	previousHash string,
) *Entry {
	entry := &Entry{
		Actor:         actor,
		Action:        action,
		ResourceType:  resourceType,
		ResourceID:    resourceID,
		Timestamp:     time.Now().UTC(),
		PreviousHash:  previousHash,
		Metadata:      types.NewMetadata(),
		SchemaVersion: SchemaVersionAudit,
	}
	entry.SetHash()
	return entry
}
