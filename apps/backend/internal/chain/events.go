package chain

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

func (c *Client) FindRootAnchoredEvent(ctx context.Context, hexRoot string) (*RootAnchoredEvent, bool, error) {
	root, err := hexToBytes32(hexRoot)
	if err != nil {
		return nil, false, fmt.Errorf("chain: invalid root: %w", err)
	}

	filterer, err := NewFlemingAnchorFilterer(c.address, c.ethClient)
	if err != nil {
		return nil, false, fmt.Errorf("chain: bind contract filterer: %w", err)
	}

	it, err := filterer.FilterRootAnchored(
		&bind.FilterOpts{Start: 0, Context: ctx},
		[][32]byte{root},
		nil,
	)
	if err != nil {
		return nil, false, fmt.Errorf("chain: filter RootAnchored logs: %w", err)
	}
	defer it.Close()

	// Roots are anchored at most once; still, scan defensively.
	var last *FlemingAnchorRootAnchored
	for it.Next() {
		last = it.Event
	}
	if err := it.Error(); err != nil {
		return nil, false, fmt.Errorf("chain: iterate RootAnchored logs: %w", err)
	}

	if last == nil {
		return nil, false, nil
	}

	return &RootAnchoredEvent{
		TxHash:      last.Raw.TxHash.Hex(),
		BlockNumber: last.BlockNumber.Uint64(),
		Timestamp:   last.Timestamp.Uint64(),
		Anchorer:    last.Anchorer,
	}, true, nil
}
