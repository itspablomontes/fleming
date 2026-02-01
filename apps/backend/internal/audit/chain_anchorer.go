package audit

import (
	"context"

	"github.com/itspablomontes/fleming/apps/backend/internal/chain"
)

// ChainAnchorer is the minimal interface the audit layer needs to anchor Merkle roots on-chain.
//
// This indirection keeps the handler + service testable without requiring a live chain.
type ChainAnchorer interface {
	AnchorRoot(ctx context.Context, hexRoot string) (*chain.AnchorResult, error)
	VerifyRoot(ctx context.Context, hexRoot string) (uint64, error)
	FindRootAnchoredEvent(ctx context.Context, hexRoot string) (*chain.RootAnchoredEvent, bool, error)
}
