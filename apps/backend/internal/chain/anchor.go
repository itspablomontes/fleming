package chain

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

// AnchorRoot submits a Merkle root to the FlemingAnchor contract.
// The root should be a 64-character hex string (SHA-256 hash from audit.MerkleTree).
func (c *Client) AnchorRoot(ctx context.Context, hexRoot string) (*AnchorResult, error) {
	// Convert hex string to bytes32
	root, err := hexToBytes32(hexRoot)
	if err != nil {
		return nil, fmt.Errorf("chain: invalid root: %w", err)
	}

	// Create transaction options
	auth, err := bind.NewKeyedTransactorWithChainID(c.signer, c.chainID)
	if err != nil {
		return nil, fmt.Errorf("chain: create transactor: %w", err)
	}
	auth.Context = ctx

	// Submit anchor transaction
	tx, err := c.contract.Anchor(auth, root)
	if err != nil {
		return nil, fmt.Errorf("chain: anchor tx: %w", err)
	}

	// Wait for transaction receipt
	receipt, err := bind.WaitMined(ctx, c.ethClient, tx)
	if err != nil {
		return nil, fmt.Errorf("chain: wait mined: %w", err)
	}

	if receipt.Status == 0 {
		return nil, fmt.Errorf("chain: transaction reverted (possibly duplicate root)")
	}

	return &AnchorResult{
		TxHash:      tx.Hash().Hex(),
		BlockNumber: receipt.BlockNumber.Uint64(),
		GasUsed:     receipt.GasUsed,
	}, nil
}

// VerifyRoot checks if a Merkle root has been anchored and returns the timestamp.
// Returns 0 if the root has not been anchored.
func (c *Client) VerifyRoot(ctx context.Context, hexRoot string) (uint64, error) {
	root, err := hexToBytes32(hexRoot)
	if err != nil {
		return 0, fmt.Errorf("chain: invalid root: %w", err)
	}

	opts := &bind.CallOpts{Context: ctx}
	timestamp, err := c.contract.GetAnchorTimestamp(opts, root)
	if err != nil {
		return 0, fmt.Errorf("chain: verify call: %w", err)
	}

	return timestamp.Uint64(), nil
}

// IsAnchored checks if a Merkle root has been anchored.
func (c *Client) IsAnchored(ctx context.Context, hexRoot string) (bool, error) {
	root, err := hexToBytes32(hexRoot)
	if err != nil {
		return false, fmt.Errorf("chain: invalid root: %w", err)
	}

	opts := &bind.CallOpts{Context: ctx}
	return c.contract.IsAnchored(opts, root)
}

// AnchorCount returns the total number of anchored roots.
func (c *Client) AnchorCount(ctx context.Context) (uint64, error) {
	opts := &bind.CallOpts{Context: ctx}
	count, err := c.contract.AnchorCount(opts)
	if err != nil {
		return 0, fmt.Errorf("chain: anchor count: %w", err)
	}
	return count.Uint64(), nil
}

// hexToBytes32 converts a 64-character hex string to [32]byte.
func hexToBytes32(hexStr string) ([32]byte, error) {
	var result [32]byte

	// Strip 0x prefix if present
	if len(hexStr) >= 2 && hexStr[:2] == "0x" {
		hexStr = hexStr[2:]
	}

	if len(hexStr) != 64 {
		return result, fmt.Errorf("expected 64 hex chars, got %d", len(hexStr))
	}

	bytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return result, fmt.Errorf("invalid hex: %w", err)
	}

	copy(result[:], bytes)
	return result, nil
}
