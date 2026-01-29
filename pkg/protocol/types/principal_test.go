package types

import (
	"testing"
)

func TestPrincipalType_IsValid(t *testing.T) {
	tests := []struct {
		pt   PrincipalType
		want bool
	}{
		{PrincipalPatient, true},
		{PrincipalProvider, true},
		{PrincipalResearcher, true},
		{PrincipalSystem, true},
		{"unknown", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(string(tt.pt), func(t *testing.T) {
			if got := tt.pt.IsValid(); got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewPrincipal(t *testing.T) {
	validAddr, _ := NewWalletAddress("0x1111111111111111111111111111111111111111")

	tests := []struct {
		name    string
		address WalletAddress
		roles   []PrincipalType
		wantErr bool
	}{
		{
			name:    "valid patient",
			address: validAddr,
			roles:   []PrincipalType{PrincipalPatient},
			wantErr: false,
		},
		{
			name:    "valid provider",
			address: validAddr,
			roles:   []PrincipalType{PrincipalProvider},
			wantErr: false,
		},
		{
			name:    "multiple roles",
			address: validAddr,
			roles:   []PrincipalType{PrincipalProvider, PrincipalResearcher},
			wantErr: false,
		},
		{
			name:    "empty address",
			address: WalletAddress(""),
			roles:   []PrincipalType{PrincipalPatient},
			wantErr: true,
		},
		{
			name:    "no roles",
			address: validAddr,
			roles:   nil,
			wantErr: true,
		},
		{
			name:    "invalid role",
			address: validAddr,
			roles:   []PrincipalType{"invalid"},
			wantErr: true,
		},
		{
			name:    "mixed valid and invalid roles",
			address: validAddr,
			roles:   []PrincipalType{PrincipalPatient, "invalid"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := NewPrincipal(tt.address, tt.roles...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPrincipal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if p.Address != tt.address {
					t.Errorf("NewPrincipal() Address = %v, want %v", p.Address, tt.address)
				}
				if len(p.Roles) != len(tt.roles) {
					t.Errorf("NewPrincipal() Roles length = %d, want %d", len(p.Roles), len(tt.roles))
				}
			}
		})
	}
}

func TestPrincipal_HasRole(t *testing.T) {
	validAddr, _ := NewWalletAddress("0x1111111111111111111111111111111111111111")

	p, _ := NewPrincipal(validAddr, PrincipalPatient, PrincipalProvider)

	if !p.HasRole(PrincipalPatient) {
		t.Error("HasRole() should return true for assigned role")
	}
	if !p.HasRole(PrincipalProvider) {
		t.Error("HasRole() should return true for assigned role")
	}
	if p.HasRole(PrincipalResearcher) {
		t.Error("HasRole() should return false for unassigned role")
	}
}

func TestPrincipal_IsPatient(t *testing.T) {
	validAddr, _ := NewWalletAddress("0x1111111111111111111111111111111111111111")

	patient, _ := NewPrincipal(validAddr, PrincipalPatient)
	if !patient.IsPatient() {
		t.Error("IsPatient() should return true for patient")
	}

	provider, _ := NewPrincipal(validAddr, PrincipalProvider)
	if provider.IsPatient() {
		t.Error("IsPatient() should return false for non-patient")
	}
}

func TestPrincipal_IsProvider(t *testing.T) {
	validAddr, _ := NewWalletAddress("0x1111111111111111111111111111111111111111")

	provider, _ := NewPrincipal(validAddr, PrincipalProvider)
	if !provider.IsProvider() {
		t.Error("IsProvider() should return true for provider")
	}

	patient, _ := NewPrincipal(validAddr, PrincipalPatient)
	if patient.IsProvider() {
		t.Error("IsProvider() should return false for non-provider")
	}
}

func TestPrincipal_CanOwn(t *testing.T) {
	validAddr, _ := NewWalletAddress("0x1111111111111111111111111111111111111111")

	patient, _ := NewPrincipal(validAddr, PrincipalPatient)
	if !patient.CanOwn() {
		t.Error("CanOwn() should return true for patient")
	}

	provider, _ := NewPrincipal(validAddr, PrincipalProvider)
	if provider.CanOwn() {
		t.Error("CanOwn() should return false for non-patient")
	}

	researcher, _ := NewPrincipal(validAddr, PrincipalResearcher)
	if researcher.CanOwn() {
		t.Error("CanOwn() should return false for researcher")
	}
}

func TestPrincipal_CanGenerate(t *testing.T) {
	validAddr, _ := NewWalletAddress("0x1111111111111111111111111111111111111111")

	provider, _ := NewPrincipal(validAddr, PrincipalProvider)
	if !provider.CanGenerate() {
		t.Error("CanGenerate() should return true for provider")
	}

	researcher, _ := NewPrincipal(validAddr, PrincipalResearcher)
	if !researcher.CanGenerate() {
		t.Error("CanGenerate() should return true for researcher")
	}

	patient, _ := NewPrincipal(validAddr, PrincipalPatient)
	if patient.CanGenerate() {
		t.Error("CanGenerate() should return false for patient")
	}

	system, _ := NewPrincipal(validAddr, PrincipalSystem)
	if system.CanGenerate() {
		t.Error("CanGenerate() should return false for system")
	}
}
