package audit

import (
	"strings"
	"testing"
	"time"

	"github.com/itspablomontes/fleming/pkg/protocol/types"
)

func TestBuildMerkleTree_ComputeRootMatches(t *testing.T) {
	entries := []Entry{
		{Hash: strings.Repeat("a", 64)},
		{Hash: strings.Repeat("b", 64)},
		{Hash: strings.Repeat("c", 64)},
		{Hash: strings.Repeat("d", 64)},
	}

	tree, err := BuildMerkleTree(entries)
	if err != nil {
		t.Fatalf("BuildMerkleTree() error = %v", err)
	}

	root, err := ComputeRoot(tree.Leaves)
	if err != nil {
		t.Fatalf("ComputeRoot() error = %v", err)
	}

	if tree.Root != root {
		t.Fatalf("tree root mismatch: got %s want %s", tree.Root, root)
	}

	if len(tree.Levels) < 2 {
		t.Fatalf("expected multiple levels, got %d", len(tree.Levels))
	}
}

func TestBuildMerkleTree_UsesComputedHashWhenMissing(t *testing.T) {
	entry := Entry{
		Actor:        types.WalletAddress("0x1234567890abcdef1234567890abcdef12345678"),
		Action:       ActionCreate,
		ResourceType: ResourceEvent,
		ResourceID:   "event-1",
		Timestamp:    time.Date(2026, 1, 25, 12, 0, 0, 0, time.UTC),
	}
	expectedHash := entry.ComputeHash()

	tree, err := BuildMerkleTree([]Entry{entry})
	if err != nil {
		t.Fatalf("BuildMerkleTree() error = %v", err)
	}

	if tree.Leaves[0] != expectedHash {
		t.Fatalf("expected computed hash leaf, got %s want %s", tree.Leaves[0], expectedHash)
	}
}

func TestGenerateProof_VerifyProof(t *testing.T) {
	entries := []Entry{
		{Hash: strings.Repeat("a", 64)},
		{Hash: strings.Repeat("b", 64)},
		{Hash: strings.Repeat("c", 64)},
	}

	tree, err := BuildMerkleTree(entries)
	if err != nil {
		t.Fatalf("BuildMerkleTree() error = %v", err)
	}

	target := entries[1].Hash
	proof, err := GenerateProof(tree, target)
	if err != nil {
		t.Fatalf("GenerateProof() error = %v", err)
	}

	if !VerifyProof(tree.Root, target, proof) {
		t.Fatal("VerifyProof() expected true")
	}
}

func TestGenerateProof_NotFound(t *testing.T) {
	tree, err := BuildMerkleTree([]Entry{{Hash: strings.Repeat("a", 64)}})
	if err != nil {
		t.Fatalf("BuildMerkleTree() error = %v", err)
	}

	_, err = GenerateProof(tree, strings.Repeat("f", 64))
	if err == nil {
		t.Fatal("expected error for missing leaf")
	}
}

func TestVerifyProof_TamperFails(t *testing.T) {
	entries := []Entry{
		{Hash: strings.Repeat("a", 64)},
		{Hash: strings.Repeat("b", 64)},
	}

	tree, err := BuildMerkleTree(entries)
	if err != nil {
		t.Fatalf("BuildMerkleTree() error = %v", err)
	}

	target := entries[0].Hash
	proof, err := GenerateProof(tree, target)
	if err != nil {
		t.Fatalf("GenerateProof() error = %v", err)
	}

	if VerifyProof(strings.Repeat("f", 64), target, proof) {
		t.Fatal("VerifyProof() should fail with wrong root")
	}
	if VerifyProof(tree.Root, strings.Repeat("f", 64), proof) {
		t.Fatal("VerifyProof() should fail with wrong entry hash")
	}
}

func TestComputeRoot_Empty(t *testing.T) {
	if _, err := ComputeRoot([]string{}); err == nil {
		t.Fatal("expected error for empty leaves")
	}
}
