package chain

import "github.com/ethereum/go-ethereum/common"

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
	Anchorer    common.Address
}
