package vc

import (
	"testing"
	"time"

	"github.com/itspablomontes/fleming/pkg/protocol/types"
)

func TestNewRevocationList(t *testing.T) {
	id, _ := types.NewID("list-1")
	issuer, _ := types.NewWalletAddress("0x1111111111111111111111111111111111111111")

	list := NewRevocationList(id, issuer)

	if list.ID != id {
		t.Errorf("NewRevocationList() ID = %v, want %v", list.ID, id)
	}
	if list.IssuerID != issuer {
		t.Errorf("NewRevocationList() IssuerID = %v, want %v", list.IssuerID, issuer)
	}
	if list.Size != DefaultRevocationListSize {
		t.Errorf("NewRevocationList() Size = %v, want %v", list.Size, DefaultRevocationListSize)
	}
	if len(list.Bitmap) != int(DefaultRevocationListSize/8) {
		t.Errorf("NewRevocationList() Bitmap length = %d, want %d", len(list.Bitmap), DefaultRevocationListSize/8)
	}
}

func TestNewRevocationListWithSize(t *testing.T) {
	id, _ := types.NewID("list-1")
	issuer, _ := types.NewWalletAddress("0x1111111111111111111111111111111111111111")
	customSize := uint64(1000)

	list := NewRevocationListWithSize(id, issuer, customSize)

	if list.Size != customSize {
		t.Errorf("NewRevocationListWithSize() Size = %v, want %v", list.Size, customSize)
	}
	// Should round up to nearest byte (1000 bits = 125 bytes)
	expectedBytes := (customSize + 7) / 8
	if len(list.Bitmap) != int(expectedBytes) {
		t.Errorf("NewRevocationListWithSize() Bitmap length = %d, want %d", len(list.Bitmap), expectedBytes)
	}
}

func TestRevocationList_IsRevoked(t *testing.T) {
	id, _ := types.NewID("list-1")
	issuer, _ := types.NewWalletAddress("0x1111111111111111111111111111111111111111")
	list := NewRevocationList(id, issuer)

	// Not revoked by default
	if list.IsRevoked(0) {
		t.Error("IsRevoked() should return false for unrevoked index")
	}

	// Revoke index 0
	list.Revoke(0)
	if !list.IsRevoked(0) {
		t.Error("IsRevoked() should return true for revoked index")
	}

	// Index out of bounds should return false
	if list.IsRevoked(list.Size+1) {
		t.Error("IsRevoked() should return false for out-of-bounds index")
	}
}

func TestRevocationList_Revoke(t *testing.T) {
	id, _ := types.NewID("list-1")
	issuer, _ := types.NewWalletAddress("0x1111111111111111111111111111111111111111")
	list := NewRevocationList(id, issuer)

	// Revoke valid index
	err := list.Revoke(5)
	if err != nil {
		t.Errorf("Revoke() error = %v", err)
	}
	if !list.IsRevoked(5) {
		t.Error("Revoke() should mark index as revoked")
	}

	// Revoke out of bounds
	err = list.Revoke(list.Size + 1)
	if err == nil {
		t.Error("Revoke() with out-of-bounds index should error")
	}

	// Should update LastUpdated
	before := list.LastUpdated
	time.Sleep(time.Millisecond)
	list.Revoke(10)
	if !list.LastUpdated.After(before) {
		t.Error("Revoke() should update LastUpdated")
	}
}

func TestRevocationList_Unrevoke(t *testing.T) {
	id, _ := types.NewID("list-1")
	issuer, _ := types.NewWalletAddress("0x1111111111111111111111111111111111111111")
	list := NewRevocationList(id, issuer)

	// Revoke then unrevoke
	list.Revoke(5)
	if !list.IsRevoked(5) {
		t.Error("Index should be revoked")
	}

	err := list.Unrevoke(5)
	if err != nil {
		t.Errorf("Unrevoke() error = %v", err)
	}
	if list.IsRevoked(5) {
		t.Error("Unrevoke() should clear revocation")
	}

	// Unrevoke out of bounds
	err = list.Unrevoke(list.Size + 1)
	if err == nil {
		t.Error("Unrevoke() with out-of-bounds index should error")
	}
}

func TestRevocationList_NextAvailableIndex(t *testing.T) {
	id, _ := types.NewID("list-1")
	issuer, _ := types.NewWalletAddress("0x1111111111111111111111111111111111111111")
	list := NewRevocationListWithSize(id, issuer, 100)

	// First available should be 0
	idx := list.NextAvailableIndex()
	if idx != 0 {
		t.Errorf("NextAvailableIndex() = %d, want 0", idx)
	}

	// Revoke first few indices
	list.Revoke(0)
	list.Revoke(1)
	list.Revoke(2)

	// Next available should be 3
	idx = list.NextAvailableIndex()
	if idx != 3 {
		t.Errorf("NextAvailableIndex() = %d, want 3", idx)
	}

	// Fill all indices
	for i := uint64(0); i < list.Size; i++ {
		list.Revoke(i)
	}

	// Should return -1 when full
	idx = list.NextAvailableIndex()
	if idx != -1 {
		t.Errorf("NextAvailableIndex() when full = %d, want -1", idx)
	}
}

func TestRevocationList_RevokedCount(t *testing.T) {
	id, _ := types.NewID("list-1")
	issuer, _ := types.NewWalletAddress("0x1111111111111111111111111111111111111111")
	list := NewRevocationListWithSize(id, issuer, 100)

	if list.RevokedCount() != 0 {
		t.Errorf("RevokedCount() = %d, want 0", list.RevokedCount())
	}

	list.Revoke(0)
	list.Revoke(5)
	list.Revoke(10)

	if list.RevokedCount() != 3 {
		t.Errorf("RevokedCount() = %d, want 3", list.RevokedCount())
	}
}

func TestRevocationList_EncodeBitmap(t *testing.T) {
	id, _ := types.NewID("list-1")
	issuer, _ := types.NewWalletAddress("0x1111111111111111111111111111111111111111")
	list := NewRevocationList(id, issuer)

	encoded := list.EncodeBitmap()
	if encoded == "" {
		t.Error("EncodeBitmap() returned empty string")
	}

	// Should be valid base64 - decode into a new list to verify
	list2 := NewRevocationList(id, issuer)
	err := list2.DecodeBitmap(encoded)
	if err != nil {
		t.Errorf("DecodeBitmap() error = %v", err)
		return
	}
	// Verify it decoded correctly by checking size matches
	if list2.Size != list.Size {
		t.Errorf("DecodeBitmap() size = %d, want %d", list2.Size, list.Size)
	}
}

func TestRevocationList_DecodeBitmap(t *testing.T) {
	id, _ := types.NewID("list-1")
	issuer, _ := types.NewWalletAddress("0x1111111111111111111111111111111111111111")
	list := NewRevocationList(id, issuer)

	// Revoke some indices
	list.Revoke(0)
	list.Revoke(5)
	originalCount := list.RevokedCount()

	// Encode and decode
	encoded := list.EncodeBitmap()
	list2 := NewRevocationList(id, issuer)
	err := list2.DecodeBitmap(encoded)
	if err != nil {
		t.Errorf("DecodeBitmap() error = %v", err)
		return
	}

	if list2.RevokedCount() != originalCount {
		t.Errorf("DecodeBitmap() revoked count = %d, want %d", list2.RevokedCount(), originalCount)
	}
	if !list2.IsRevoked(0) || !list2.IsRevoked(5) {
		t.Error("DecodeBitmap() should preserve revoked indices")
	}

	// Invalid base64
	list3 := NewRevocationList(id, issuer)
	err = list3.DecodeBitmap("invalid-base64!")
	if err == nil {
		t.Error("DecodeBitmap() with invalid base64 should error")
	}
}

func TestRevocationList_Validate(t *testing.T) {
	id, _ := types.NewID("list-1")
	issuer, _ := types.NewWalletAddress("0x1111111111111111111111111111111111111111")

	tests := []struct {
		name    string
		list    *RevocationList
		wantErr bool
	}{
		{
			name:    "valid list",
			list:    NewRevocationList(id, issuer),
			wantErr: false,
		},
		{
			name: "missing ID",
			list: &RevocationList{
				IssuerID: issuer,
				Bitmap:   make([]byte, 10),
				Size:     80,
			},
			wantErr: true,
		},
		{
			name: "missing issuer",
			list: &RevocationList{
				ID:     id,
				Bitmap: make([]byte, 10),
				Size:   80,
			},
			wantErr: true,
		},
		{
			name: "empty bitmap",
			list: &RevocationList{
				ID:       id,
				IssuerID: issuer,
				Bitmap:   nil,
				Size:    0,
			},
			wantErr: true,
		},
		{
			name: "zero size",
			list: &RevocationList{
				ID:       id,
				IssuerID: issuer,
				Bitmap:   make([]byte, 10),
				Size:     0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.list.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRevocationRegistry(t *testing.T) {
	registry := NewRevocationRegistry()

	id, _ := types.NewID("list-1")
	issuer, _ := types.NewWalletAddress("0x1111111111111111111111111111111111111111")
	list := NewRevocationList(id, issuer)

	registry.Register(list)

	retrieved, ok := registry.Get(id)
	if !ok {
		t.Error("Get() should find registered list")
	}
	if retrieved.ID != id {
		t.Errorf("Get() ID = %v, want %v", retrieved.ID, id)
	}

	// Check status
	status, err := registry.CheckStatus(id, 0)
	if err != nil {
		t.Errorf("CheckStatus() error = %v", err)
		return
	}
	if status.ListID != id {
		t.Errorf("CheckStatus() ListID = %v, want %v", status.ListID, id)
	}
	if status.Index != 0 {
		t.Errorf("CheckStatus() Index = %v, want 0", status.Index)
	}

	// Revoke and check again
	list.Revoke(0)
	status2, _ := registry.CheckStatus(id, 0)
	if !status2.IsRevoked {
		t.Error("CheckStatus() should reflect revocation")
	}

	// Nonexistent list
	_, err = registry.CheckStatus(types.ID("nonexistent"), 0)
	if err == nil {
		t.Error("CheckStatus() with nonexistent list should error")
	}
}
