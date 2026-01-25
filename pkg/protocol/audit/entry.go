package audit

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"slices"
	"time"

	"github.com/itspablomontes/fleming/pkg/protocol/types"
)

type Action string

const (
	ActionCreate Action = "create"
	ActionRead   Action = "read"
	ActionUpdate Action = "update"
	ActionDelete Action = "delete"

	ActionConsentRequest Action = "consent.request"
	ActionConsentApprove Action = "consent.approve"
	ActionConsentDeny    Action = "consent.deny"
	ActionConsentRevoke  Action = "consent.revoke"
	ActionConsentExpire  Action = "consent.expire"

	ActionLogin  Action = "auth.login"
	ActionLogout Action = "auth.logout"

	ActionUpload   Action = "file.upload"
	ActionDownload Action = "file.download"
	ActionShare    Action = "file.share"
)

func ValidActions() []Action {
	return []Action{
		ActionCreate, ActionRead, ActionUpdate, ActionDelete,
		ActionConsentRequest, ActionConsentApprove, ActionConsentDeny,
		ActionConsentRevoke, ActionConsentExpire,
		ActionLogin, ActionLogout,
		ActionUpload, ActionDownload, ActionShare,
	}
}

func (a Action) IsValid() bool {
	return slices.Contains(ValidActions(), a)
}

type ResourceType string

const (
	ResourceEvent   ResourceType = "event"
	ResourceFile    ResourceType = "file"
	ResourceConsent ResourceType = "consent"
	ResourceSession ResourceType = "session"
)

type Entry struct {
	ID types.ID `json:"id"`

	Actor types.WalletAddress `json:"actor"`

	Action Action `json:"action"`

	ResourceType ResourceType `json:"resourceType"`

	ResourceID types.ID `json:"resourceId"`

	Timestamp time.Time `json:"timestamp"`

	Metadata types.Metadata `json:"metadata,omitempty"`

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
		Actor:        actor,
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		Timestamp:    time.Now().UTC(),
		PreviousHash: previousHash,
		Metadata:     types.NewMetadata(),
	}
	entry.SetHash()
	return entry
}
