package vc

import (
	"testing"
	"time"

	"github.com/itspablomontes/fleming/pkg/protocol/types"
)

func TestCredentialBuilder_WithID(t *testing.T) {
	validID, _ := types.NewID("cred-1")
	builder := NewCredentialBuilder()

	builder.WithID(validID)
	if builder.cred.ID != validID {
		t.Error("WithID() did not set ID")
	}

	// Empty ID should add error
	emptyID := types.ID("")
	builder2 := NewCredentialBuilder()
	builder2.WithID(emptyID)
	if !builder2.errs.HasErrors() {
		t.Error("WithID() with empty ID should add error")
	}
}

func TestCredentialBuilder_WithIssuer(t *testing.T) {
	validAddr, _ := types.NewWalletAddress("0x1111111111111111111111111111111111111111")
	builder := NewCredentialBuilder()

	builder.WithIssuer(validAddr)
	if builder.cred.Issuer != validAddr {
		t.Error("WithIssuer() did not set issuer")
	}

	// Empty address should add error
	emptyAddr := types.WalletAddress("")
	builder2 := NewCredentialBuilder()
	builder2.WithIssuer(emptyAddr)
	if !builder2.errs.HasErrors() {
		t.Error("WithIssuer() with empty address should add error")
	}
}

func TestCredentialBuilder_WithSubject(t *testing.T) {
	validAddr, _ := types.NewWalletAddress("0x2222222222222222222222222222222222222222")
	builder := NewCredentialBuilder()

	builder.WithSubject(validAddr)
	if builder.cred.Subject != validAddr {
		t.Error("WithSubject() did not set subject")
	}

	// Empty address should add error
	emptyAddr := types.WalletAddress("")
	builder2 := NewCredentialBuilder()
	builder2.WithSubject(emptyAddr)
	if !builder2.errs.HasErrors() {
		t.Error("WithSubject() with empty address should add error")
	}
}

func TestCredentialBuilder_WithClaimType(t *testing.T) {
	builder := NewCredentialBuilder()

	builder.WithClaimType(ClaimBloodworkRange)
	if builder.cred.ClaimType != ClaimBloodworkRange {
		t.Error("WithClaimType() did not set claim type")
	}

	// Invalid claim type should add error
	builder2 := NewCredentialBuilder()
	builder2.WithClaimType("invalid")
	if !builder2.errs.HasErrors() {
		t.Error("WithClaimType() with invalid type should add error")
	}
}

func TestCredentialBuilder_AddClaim(t *testing.T) {
	builder := NewCredentialBuilder()

	builder.AddClaim("marker", "718-7", false)
	if builder.cred.Claims["marker"] != "718-7" {
		t.Error("AddClaim() did not add claim")
	}
	if len(builder.cred.Disclosures) != 0 {
		t.Error("AddClaim() with disclosed=false should not add disclosure")
	}

	builder.AddClaim("value", 15.0, true)
	if len(builder.cred.Disclosures) != 1 {
		t.Error("AddClaim() with disclosed=true should add disclosure")
	}

	// Empty key should add error
	builder2 := NewCredentialBuilder()
	builder2.AddClaim("", "value", false)
	if !builder2.errs.HasErrors() {
		t.Error("AddClaim() with empty key should add error")
	}
}

func TestCredentialBuilder_AddBloodworkClaim(t *testing.T) {
	validClaim := &BloodworkRangeClaim{
		Marker:       "718-7",
		RangeMin:     13.5,
		RangeMax:     17.5,
		WindowMonths: 6,
		AllInRange:   true,
		SampleCount:  5,
	}

	builder := NewCredentialBuilder()
	builder.AddBloodworkClaim(validClaim)

	if builder.cred.ClaimType != ClaimBloodworkRange {
		t.Error("AddBloodworkClaim() should set claim type")
	}
	if builder.cred.Claims["marker"] != "718-7" {
		t.Error("AddBloodworkClaim() should add claim fields")
	}

	// Invalid claim should add error
	invalidClaim := &BloodworkRangeClaim{
		Marker: "", // Missing required field
	}
	builder2 := NewCredentialBuilder()
	builder2.AddBloodworkClaim(invalidClaim)
	if !builder2.errs.HasErrors() {
		t.Error("AddBloodworkClaim() with invalid claim should add error")
	}
}

func TestCredentialBuilder_WithSourceEvents(t *testing.T) {
	eventID1, _ := types.NewID("event-1")
	eventID2, _ := types.NewID("event-2")

	builder := NewCredentialBuilder()
	builder.WithSourceEvents(eventID1, eventID2)

	if len(builder.cred.SourceEventIDs) != 2 {
		t.Errorf("WithSourceEvents() expected 2 events, got %d", len(builder.cred.SourceEventIDs))
	}

	// Empty events should add error
	builder2 := NewCredentialBuilder()
	builder2.WithSourceEvents()
	if !builder2.errs.HasErrors() {
		t.Error("WithSourceEvents() with no events should add error")
	}
}

func TestCredentialBuilder_WithExpiresAt(t *testing.T) {
	future := time.Now().Add(time.Hour)
	builder := NewCredentialBuilder()

	builder.WithExpiresAt(future)
	if builder.cred.ExpiresAt == nil || *builder.cred.ExpiresAt != future {
		t.Error("WithExpiresAt() did not set expiration")
	}
}

func TestCredentialBuilder_WithTTL(t *testing.T) {
	builder := NewCredentialBuilder()
	ttl := 24 * time.Hour

	builder.WithTTL(ttl)
	if builder.cred.ExpiresAt == nil {
		t.Error("WithTTL() should set ExpiresAt")
	}
	if builder.cred.ExpiresAt.Before(time.Now()) {
		t.Error("WithTTL() should set future expiration")
	}
}

func TestCredentialBuilder_Build(t *testing.T) {
	validIssuer, _ := types.NewWalletAddress("0x1111111111111111111111111111111111111111")
	validSubject, _ := types.NewWalletAddress("0x2222222222222222222222222222222222222222")
	eventID, _ := types.NewID("event-1")

	tests := []struct {
		name    string
		builder func() *CredentialBuilder
		wantErr bool
	}{
		{
			name: "valid credential",
			builder: func() *CredentialBuilder {
				return NewCredentialBuilder().
					WithIssuer(validIssuer).
					WithSubject(validSubject).
					WithClaimType(ClaimBloodworkRange).
					AddClaim("marker", "718-7", false).
					WithSourceEvents(eventID)
			},
			wantErr: false,
		},
		{
			name: "missing issuer",
			builder: func() *CredentialBuilder {
				return NewCredentialBuilder().
					WithSubject(validSubject).
					WithClaimType(ClaimBloodworkRange).
					AddClaim("marker", "718-7", false)
			},
			wantErr: true,
		},
		{
			name: "missing subject",
			builder: func() *CredentialBuilder {
				return NewCredentialBuilder().
					WithIssuer(validIssuer).
					WithClaimType(ClaimBloodworkRange).
					AddClaim("marker", "718-7", false)
			},
			wantErr: true,
		},
		{
			name: "invalid claim type",
			builder: func() *CredentialBuilder {
				return NewCredentialBuilder().
					WithIssuer(validIssuer).
					WithSubject(validSubject).
					WithClaimType("invalid").
					AddClaim("marker", "718-7", false)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := tt.builder()
			cred, err := builder.Build()
			if (err != nil) != tt.wantErr {
				t.Errorf("Build() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if cred == nil {
					t.Error("Build() returned nil for valid credential")
				}
				if cred.Status != StatusActive {
					t.Errorf("Build() status = %v, want %v", cred.Status, StatusActive)
				}
			}
		})
	}
}

func TestPresentationBuilder_DiscloseKey(t *testing.T) {
	validIssuer, _ := types.NewWalletAddress("0x1111111111111111111111111111111111111111")
	validSubject, _ := types.NewWalletAddress("0x2222222222222222222222222222222222222222")
	eventID, _ := types.NewID("event-1")

	cred, _ := NewCredentialBuilder().
		WithIssuer(validIssuer).
		WithSubject(validSubject).
		WithClaimType(ClaimBloodworkRange).
		AddClaim("marker", "718-7", false).
		AddClaim("value", 15.0, false).
		WithSourceEvents(eventID).
		Build()

	pb := NewPresentationBuilder(cred)
	pb.DiscloseKey("marker")

	if !pb.disclosedKeys["marker"] {
		t.Error("DiscloseKey() should mark key as disclosed")
	}
	if pb.disclosedKeys["value"] {
		t.Error("DiscloseKey() should not mark other keys")
	}

	// Invalid key should add error
	pb2 := NewPresentationBuilder(cred)
	pb2.DiscloseKey("nonexistent")
	if !pb2.errs.HasErrors() {
		t.Error("DiscloseKey() with nonexistent key should add error")
	}
}

func TestPresentationBuilder_DiscloseAll(t *testing.T) {
	validIssuer, _ := types.NewWalletAddress("0x1111111111111111111111111111111111111111")
	validSubject, _ := types.NewWalletAddress("0x2222222222222222222222222222222222222222")
	eventID, _ := types.NewID("event-1")

	cred, _ := NewCredentialBuilder().
		WithIssuer(validIssuer).
		WithSubject(validSubject).
		WithClaimType(ClaimBloodworkRange).
		AddClaim("marker", "718-7", false).
		AddClaim("value", 15.0, false).
		WithSourceEvents(eventID).
		Build()

	pb := NewPresentationBuilder(cred)
	pb.DiscloseAll()

	if len(pb.disclosedKeys) != 2 {
		t.Errorf("DiscloseAll() expected 2 keys, got %d", len(pb.disclosedKeys))
	}
	if !pb.disclosedKeys["marker"] || !pb.disclosedKeys["value"] {
		t.Error("DiscloseAll() should mark all keys as disclosed")
	}
}

func TestPresentationBuilder_Build(t *testing.T) {
	validIssuer, _ := types.NewWalletAddress("0x1111111111111111111111111111111111111111")
	validSubject, _ := types.NewWalletAddress("0x2222222222222222222222222222222222222222")
	eventID, _ := types.NewID("event-1")

	cred, _ := NewCredentialBuilder().
		WithIssuer(validIssuer).
		WithSubject(validSubject).
		WithClaimType(ClaimBloodworkRange).
		AddClaim("marker", "718-7", false).
		AddClaim("value", 15.0, false).
		WithSourceEvents(eventID).
		Build()

	pb := NewPresentationBuilder(cred)
	pb.DiscloseKey("marker")

	presentation, err := pb.Build()
	if err != nil {
		t.Errorf("Build() error = %v", err)
		return
	}

	if len(presentation.Claims) != 1 {
		t.Errorf("Build() expected 1 disclosed claim, got %d", len(presentation.Claims))
	}
	if presentation.Claims["marker"] != "718-7" {
		t.Error("Build() should include disclosed claims")
	}
	if presentation.Claims["value"] != nil {
		t.Error("Build() should not include undisclosed claims")
	}

	// Unusable credential should error
	cred.Status = StatusRevoked
	pb2 := NewPresentationBuilder(cred)
	_, err = pb2.Build()
	if err == nil {
		t.Error("Build() with unusable credential should error")
	}
}
