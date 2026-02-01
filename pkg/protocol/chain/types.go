// Package chain provides pure types for on-chain anchoring results.
//
// Implementation lives in apps/backend/internal/chain/ to keep the protocol
// layer free of infrastructure dependencies (go-ethereum, RPC clients, etc.).
package chain

// AnchorResult represents the result of anchoring a Merkle root.
type AnchorResult struct {
	TxHash      string
	BlockNumber uint64
	GasUsed     uint64
}

// RootAnchoredEvent is a minimal view of a RootAnchored on-chain event.
type RootAnchoredEvent struct {
	TxHash      string
	BlockNumber uint64
	Timestamp   uint64
	Anchorer    string // hex address
}
