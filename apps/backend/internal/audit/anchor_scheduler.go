package audit

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"
)

type AnchorScheduler struct {
	repo        Repository
	service     Service
	chainClient ChainAnchorer

	interval  time.Duration
	window    time.Duration
	maxActors int
}

func NewAnchorScheduler(repo Repository, service Service, chainClient ChainAnchorer) (*AnchorScheduler, error) {
	if repo == nil {
		return nil, fmt.Errorf("audit: anchor scheduler: repo is nil")
	}
	if service == nil {
		return nil, fmt.Errorf("audit: anchor scheduler: service is nil")
	}
	if chainClient == nil {
		return nil, fmt.Errorf("audit: anchor scheduler: chain client is nil")
	}

	interval, err := parseDurationEnv("MERKLE_AUTO_ANCHOR_INTERVAL", 24*time.Hour)
	if err != nil {
		return nil, err
	}
	window, err := parseDurationEnv("MERKLE_AUTO_ANCHOR_WINDOW", 24*time.Hour)
	if err != nil {
		return nil, err
	}
	if interval <= 0 {
		return nil, fmt.Errorf("audit: anchor scheduler: interval must be > 0")
	}
	if window <= 0 {
		return nil, fmt.Errorf("audit: anchor scheduler: window must be > 0")
	}

	maxActors := 500
	if v := strings.TrimSpace(os.Getenv("MERKLE_AUTO_ANCHOR_MAX_ACTORS")); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n <= 0 {
			return nil, fmt.Errorf("audit: anchor scheduler: invalid MERKLE_AUTO_ANCHOR_MAX_ACTORS")
		}
		maxActors = n
	}

	return &AnchorScheduler{
		repo:        repo,
		service:     service,
		chainClient: chainClient,
		interval:    interval,
		window:      window,
		maxActors:   maxActors,
	}, nil
}

func parseDurationEnv(key string, defaultValue time.Duration) (time.Duration, error) {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return defaultValue, nil
	}
	d, err := time.ParseDuration(raw)
	if err != nil {
		return 0, fmt.Errorf("audit: anchor scheduler: invalid %s: %w", key, err)
	}
	return d, nil
}

func (s *AnchorScheduler) Start(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	go func() {
		defer ticker.Stop()

		// Run once at startup, then on interval.
		s.runOnce(ctx)

		for {
			select {
			case <-ticker.C:
				s.runOnce(ctx)
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (s *AnchorScheduler) runOnce(ctx context.Context) {
	end := alignToInterval(time.Now().UTC(), s.window)
	start := end.Add(-s.window)

	slog.Info(
		"audit: auto-anchor tick",
		"windowStart", start.Format(time.RFC3339),
		"windowEnd", end.Format(time.RFC3339),
	)

	actors, err := s.repo.GetDistinctActorsWithEntries(ctx, start, end, s.maxActors)
	if err != nil {
		slog.Error("audit: auto-anchor: list actors failed", "error", err)
		return
	}
	if len(actors) == 0 {
		slog.Debug("audit: auto-anchor: no actors with entries in window")
		return
	}

	for _, actor := range actors {
		// Best effort per actor; do not stop whole job.
		if err := s.anchorActorWindow(ctx, actor, start, end); err != nil {
			slog.Warn("audit: auto-anchor: actor window failed", "actor", actor, "error", err)
		}
	}
}

func alignToInterval(t time.Time, interval time.Duration) time.Time {
	// Align by truncating relative to Unix epoch.
	sec := t.Unix()
	step := int64(interval.Seconds())
	if step <= 0 {
		return t
	}
	aligned := (sec / step) * step
	return time.Unix(aligned, 0).UTC()
}

func (s *AnchorScheduler) anchorActorWindow(ctx context.Context, actor string, start time.Time, end time.Time) error {
	// Build batch for this actor & window (idempotent via actor+root uniqueness).
	batch, _, err := s.service.BuildMerkleTree(ctx, actor, start, end)
	if err != nil {
		// Common case: no entries in this range (e.g. actor list might be stale).
		return err
	}
	if batch == nil {
		return nil
	}

	// Anchor (idempotent).
	updated, err := s.service.AnchorBatch(ctx, actor, batch.ID, s.chainClient)
	if err != nil {
		return err
	}
	if updated != nil && updated.AnchorStatus == anchorStatusAnchored {
		slog.Info(
			"audit: auto-anchor: anchored batch",
			"actor", actor,
			"batchId", updated.ID,
			"root", updated.RootHash,
			"blockNumber", updated.AnchorBlockNumber,
			"txHash", updated.AnchorTxHash,
		)
	}

	return nil
}
