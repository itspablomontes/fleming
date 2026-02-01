// Package chain provides blockchain integration for on-chain anchoring.
package chain

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Client wraps an Ethereum client for FlemingAnchor interactions.
type Client struct {
	ethClient *ethclient.Client
	contract  *FlemingAnchor
	address   common.Address
	chainID   *big.Int
	signer    *ecdsa.PrivateKey
}

// Config holds configuration for the chain client.
type Config struct {
	RPCURL          string
	ContractAddress string
	PrivateKey      string // hex-encoded, with or without 0x prefix
}

// NewClient creates a new chain client.
func NewClient(ctx context.Context, cfg Config) (*Client, error) {
	ethClient, err := ethclient.DialContext(ctx, cfg.RPCURL)
	if err != nil {
		return nil, fmt.Errorf("chain: dial rpc: %w", err)
	}

	chainID, err := ethClient.ChainID(ctx)
	if err != nil {
		return nil, fmt.Errorf("chain: get chain id: %w", err)
	}

	contractAddr := common.HexToAddress(cfg.ContractAddress)
	contract, err := NewFlemingAnchor(contractAddr, ethClient)
	if err != nil {
		return nil, fmt.Errorf("chain: bind contract: %w", err)
	}

	// Parse private key (strip 0x prefix if present)
	pkHex := cfg.PrivateKey
	if len(pkHex) >= 2 && pkHex[:2] == "0x" {
		pkHex = pkHex[2:]
	}
	privateKey, err := crypto.HexToECDSA(pkHex)
	if err != nil {
		return nil, fmt.Errorf("chain: parse private key: %w", err)
	}

	return &Client{
		ethClient: ethClient,
		contract:  contract,
		address:   contractAddr,
		chainID:   chainID,
		signer:    privateKey,
	}, nil
}

// Close closes the underlying Ethereum client.
func (c *Client) Close() {
	c.ethClient.Close()
}

// ContractAddress returns the FlemingAnchor contract address.
func (c *Client) ContractAddress() common.Address {
	return c.address
}
