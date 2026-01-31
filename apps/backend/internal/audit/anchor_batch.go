package audit

import (
	"context"
	"fmt"
	"strings"
	"time"
)

const (
	anchorStatusPending  = "pending"
	anchorStatusAnchored = "anchored"
	anchorStatusFailed   = "failed"
)

func sanitizeAnchorError(err error) string {
	if err == nil {
		return ""
	}
	msg := strings.TrimSpace(err.Error())
	if msg == "" {
		return "unknown error"
	}

	const maxLen = 500
	if len(msg) > maxLen {
		return msg[:maxLen] + "…"
	}
	return msg
}

func (s *service) AnchorBatch(ctx context.Context, actor string, batchID string, chainClient ChainAnchorer) (*AuditBatch, error) {
	if chainClient == nil {
		return nil, fmt.Errorf("anchor batch: chain client is nil")
	}

	batch, err := s.GetBatch(ctx, actor, batchID)
	if err != nil {
		return nil, fmt.Errorf("anchor batch: load batch: %w", err)
	}
	if batch == nil {
		return nil, nil
	}

	if batch.AnchorStatus == anchorStatusAnchored && batch.AnchoredAt != nil && batch.AnchorTxHash != nil && batch.AnchorBlockNumber != nil {
		return batch, nil
	}

	batch.AnchorStatus = anchorStatusPending
	batch.AnchorError = nil
	if err := s.repo.UpdateBatch(ctx, batch); err != nil {
		return nil, fmt.Errorf("anchor batch: persist pending: %w", err)
	}

	res, err := chainClient.AnchorRoot(ctx, batch.RootHash)
	if err != nil {
		msg := sanitizeAnchorError(err)
		batch.AnchorStatus = anchorStatusFailed
		batch.AnchorError = &msg
		_ = s.repo.UpdateBatch(ctx, batch)
		return batch, fmt.Errorf("anchor batch: anchor root: %w", err)
	}

	batch.AnchorTxHash = &res.TxHash
	batch.AnchorBlockNumber = &res.BlockNumber

	anchoredAtUnix, err := chainClient.VerifyRoot(ctx, batch.RootHash)
	if err != nil {
		msg := sanitizeAnchorError(err)
		batch.AnchorStatus = anchorStatusFailed
		batch.AnchorError = &msg
		_ = s.repo.UpdateBatch(ctx, batch)
		return batch, fmt.Errorf("anchor batch: verify root: %w", err)
	}
	if anchoredAtUnix == 0 {
		msg := "verify returned 0 after successful anchor"
		batch.AnchorStatus = anchorStatusFailed
		batch.AnchorError = &msg
		_ = s.repo.UpdateBatch(ctx, batch)
		return batch, fmt.Errorf("anchor batch: %s", msg)
	}

	anchoredAt := time.Unix(int64(anchoredAtUnix), 0).UTC()
	batch.AnchoredAt = &anchoredAt
	batch.AnchorStatus = anchorStatusAnchored
	batch.AnchorError = nil

	if err := s.repo.UpdateBatch(ctx, batch); err != nil {
		return nil, fmt.Errorf("anchor batch: persist anchored: %w", err)
	}

	return batch, nil
}
