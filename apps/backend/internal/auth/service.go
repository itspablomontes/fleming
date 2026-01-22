package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/itspablomontes/fleming/pkg/protocol/crypto"
)

var (
	ErrInvalidSignature = errors.New("invalid signature")
	ErrChallengeExpired = errors.New("challenge expired or not found")
)

type ChallengeRequest struct {
	Address string
	Domain  string
	URI     string
	ChainID int
}

type Service struct {
	repo      Repository
	jwtSecret []byte
}

func NewService(repo Repository, jwtSecret string) *Service {
	return &Service{
		repo:      repo,
		jwtSecret: []byte(jwtSecret),
	}
}

func buildSIWEMessage(address, nonce, domain, uri string, chainID int) string {
	issuedAt := time.Now().UTC().Format(time.RFC3339)
	return fmt.Sprintf(`%s wants you to sign in with your Ethereum account:
%s

Sign in to Fleming for secure access to your medical data.

URI: %s
Version: 1
Chain ID: %d
Nonce: %s
Issued At: %s`, domain, address, uri, chainID, nonce, issuedAt)
}

func (s *Service) GenerateChallenge(ctx context.Context, req ChallengeRequest) (string, error) {
	nonceBytes := make([]byte, 32)
	if _, err := rand.Read(nonceBytes); err != nil {
		return "", fmt.Errorf("failed to generate random nonce: %w", err)
	}
	nonce := hex.EncodeToString(nonceBytes)

	message := buildSIWEMessage(req.Address, nonce, req.Domain, req.URI, req.ChainID)
	expiresAt := time.Now().Add(5 * time.Minute)

	challenge := &Challenge{
		Address:   req.Address,
		Message:   message,
		ExpiresAt: expiresAt,
	}

	if err := s.repo.SaveChallenge(ctx, challenge); err != nil {
		slog.Error("failed to store challenge", "address", req.Address, "error", err)
		return "", fmt.Errorf("failed to store challenge: %w", err)
	}

	slog.Debug("challenge generated", "address", req.Address, "expiresAt", expiresAt)
	return message, nil
}

func (s *Service) ValidateChallenge(ctx context.Context, address, signature string) (string, error) {
	challenge, err := s.repo.FindChallenge(ctx, address)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			slog.Warn("challenge not found", "address", address)
			return "", ErrChallengeExpired
		}
		slog.Error("failed to retrieve challenge", "address", address, "error", err)
		return "", fmt.Errorf("failed to retrieve challenge: %w", err)
	}

	if time.Now().After(challenge.ExpiresAt) {
		slog.Warn("challenge expired", "address", address, "expiresAt", challenge.ExpiresAt)
		s.deleteChallenge(ctx, address)
		return "", ErrChallengeExpired
	}

	if !crypto.VerifySignature(challenge.Message, signature, address) {
		slog.Warn("invalid signature", "address", address)
		return "", ErrInvalidSignature
	}
	s.deleteChallenge(ctx, address)

	slog.Info("auth successful", "address", address)
	return s.issueJWT(address)
}

func (s *Service) deleteChallenge(ctx context.Context, address string) {
	if err := s.repo.DeleteChallenge(ctx, address); err != nil {
		slog.Warn("failed to delete challenge", "address", address, "error", err)
	}
}

func (s *Service) StartCleanup(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Minute)
	go func() {
		for {
			select {
			case <-ticker.C:
				count, err := s.repo.DeleteExpiredChallenges(ctx)
				if err != nil {
					slog.Warn("challenge cleanup failed", "error", err)
					continue
				}
				if count > 0 {
					slog.Debug("cleaned up expired challenges", "count", count)
				}
			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()
}

func (s *Service) issueJWT(address string) (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": address,
		"exp": now.Add(24 * time.Hour).Unix(),
		"iat": now.Unix(),
	})

	return token.SignedString(s.jwtSecret)
}

func (s *Service) ValidateJWT(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		slog.Debug("JWT validation failed", "error", err)
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		sub, _ := claims["sub"].(string)
		slog.Debug("JWT validated", "address", sub)
		return sub, nil
	}

	slog.Debug("JWT invalid: claims mismatch or token invalid")
	return "", errors.New("invalid token")
}
