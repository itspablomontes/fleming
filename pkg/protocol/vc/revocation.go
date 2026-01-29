package vc

import (
	"encoding/base64"
	"fmt"
	"sync"
	"time"

	"github.com/itspablomontes/fleming/pkg/protocol/types"
)

// RevocationList implements a bitmap-based revocation status list
// following the W3C VC Status List 2021 pattern.
// Each credential is assigned an index in the bitmap, and the bit at that index
// indicates revocation status (1 = revoked, 0 = not revoked).
type RevocationList struct {
	mu sync.RWMutex

	// ID is the unique identifier for this revocation list
	ID types.ID `json:"id"`

	// IssuerID is the wallet address of the issuer who controls this list
	IssuerID types.WalletAddress `json:"issuerId"`

	// Purpose describes what this list is for (e.g., "revocation", "suspension")
	Purpose string `json:"purpose"`

	// Bitmap is the bit array where each bit represents a credential's status
	// Encoded as base64 for storage/transmission
	Bitmap []byte `json:"bitmap"`

	// Size is the number of credentials this list can track
	Size uint64 `json:"size"`

	// LastUpdated is when the list was last modified
	LastUpdated time.Time `json:"lastUpdated"`

	// SchemaVersion is the protocol schema version
	SchemaVersion string `json:"schemaVersion"`
}

// DefaultRevocationListSize is the default size (16KB = 131,072 bits = credentials)
const DefaultRevocationListSize = 16 * 1024 * 8

// NewRevocationList creates a new revocation list.
func NewRevocationList(id types.ID, issuerID types.WalletAddress) *RevocationList {
	return NewRevocationListWithSize(id, issuerID, DefaultRevocationListSize)
}

// NewRevocationListWithSize creates a new revocation list with custom size.
func NewRevocationListWithSize(id types.ID, issuerID types.WalletAddress, size uint64) *RevocationList {
	// Round up to nearest byte
	byteSize := (size + 7) / 8
	return &RevocationList{
		ID:            id,
		IssuerID:      issuerID,
		Purpose:       "revocation",
		Bitmap:        make([]byte, byteSize),
		Size:          size,
		LastUpdated:   time.Now().UTC(),
		SchemaVersion: SchemaVersionVC,
	}
}

// IsRevoked checks if a credential at the given index is revoked.
func (r *RevocationList) IsRevoked(index uint64) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if index >= r.Size {
		return false
	}

	byteIndex := index / 8
	bitIndex := index % 8

	return (r.Bitmap[byteIndex] & (1 << bitIndex)) != 0
}

// Revoke marks a credential at the given index as revoked.
func (r *RevocationList) Revoke(index uint64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if index >= r.Size {
		return fmt.Errorf("index %d exceeds list size %d", index, r.Size)
	}

	byteIndex := index / 8
	bitIndex := index % 8

	r.Bitmap[byteIndex] |= (1 << bitIndex)
	r.LastUpdated = time.Now().UTC()

	return nil
}

// Unrevoke clears the revocation status for a credential.
// Use with caution - typically credentials should not be unrevoked.
func (r *RevocationList) Unrevoke(index uint64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if index >= r.Size {
		return fmt.Errorf("index %d exceeds list size %d", index, r.Size)
	}

	byteIndex := index / 8
	bitIndex := index % 8

	r.Bitmap[byteIndex] &= ^(1 << bitIndex)
	r.LastUpdated = time.Now().UTC()

	return nil
}

// NextAvailableIndex finds the next available (not revoked) index.
// This can be used to allocate indices to new credentials.
// Returns -1 if no available index exists.
func (r *RevocationList) NextAvailableIndex() int64 {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for i := uint64(0); i < r.Size; i++ {
		byteIndex := i / 8
		bitIndex := i % 8
		if (r.Bitmap[byteIndex] & (1 << bitIndex)) == 0 {
			return int64(i)
		}
	}

	return -1 // No available index
}

// RevokedCount returns the number of revoked credentials.
func (r *RevocationList) RevokedCount() uint64 {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var count uint64
	for _, b := range r.Bitmap {
		// Count set bits (Brian Kernighan's algorithm)
		for b != 0 {
			count++
			b &= b - 1
		}
	}
	return count
}

// EncodeBitmap returns the base64-encoded bitmap.
func (r *RevocationList) EncodeBitmap() string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return base64.StdEncoding.EncodeToString(r.Bitmap)
}

// DecodeBitmap decodes a base64-encoded bitmap.
func (r *RevocationList) DecodeBitmap(encoded string) error {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return fmt.Errorf("failed to decode bitmap: %w", err)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.Bitmap = decoded
	r.Size = uint64(len(decoded) * 8)
	r.LastUpdated = time.Now().UTC()

	return nil
}

// Validate validates the revocation list structure.
func (r *RevocationList) Validate() error {
	var errs types.ValidationErrors

	if r.ID.IsEmpty() {
		errs.Add("id", "ID is required")
	}

	if r.IssuerID.IsEmpty() {
		errs.Add("issuerId", "issuer ID is required")
	}

	if len(r.Bitmap) == 0 {
		errs.Add("bitmap", "bitmap cannot be empty")
	}

	if r.Size == 0 {
		errs.Add("size", "size must be positive")
	}

	if errs.HasErrors() {
		return errs
	}
	return nil
}

// RevocationStatus represents the status of a credential in a revocation list.
type RevocationStatus struct {
	// ListID is the ID of the revocation list
	ListID types.ID `json:"listId"`

	// Index is the credential's position in the list
	Index uint64 `json:"index"`

	// IsRevoked indicates whether the credential is revoked
	IsRevoked bool `json:"isRevoked"`

	// CheckedAt is when the status was checked
	CheckedAt time.Time `json:"checkedAt"`
}

// CheckRevocationStatus checks a credential's revocation status.
func CheckRevocationStatus(list *RevocationList, index uint64) *RevocationStatus {
	return &RevocationStatus{
		ListID:    list.ID,
		Index:     index,
		IsRevoked: list.IsRevoked(index),
		CheckedAt: time.Now().UTC(),
	}
}

// RevocationRegistry manages multiple revocation lists.
type RevocationRegistry struct {
	mu    sync.RWMutex
	lists map[types.ID]*RevocationList
}

// NewRevocationRegistry creates a new revocation registry.
func NewRevocationRegistry() *RevocationRegistry {
	return &RevocationRegistry{
		lists: make(map[types.ID]*RevocationList),
	}
}

// Register adds a revocation list to the registry.
func (r *RevocationRegistry) Register(list *RevocationList) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.lists[list.ID] = list
}

// Get retrieves a revocation list by ID.
func (r *RevocationRegistry) Get(id types.ID) (*RevocationList, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list, ok := r.lists[id]
	return list, ok
}

// CheckStatus checks a credential's revocation status by list ID and index.
func (r *RevocationRegistry) CheckStatus(listID types.ID, index uint64) (*RevocationStatus, error) {
	list, ok := r.Get(listID)
	if !ok {
		return nil, fmt.Errorf("revocation list not found: %s", listID)
	}
	return CheckRevocationStatus(list, index), nil
}
