package types

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

type ID string

func NewID(s string) (ID, error) {
	if s == "" {
		return "", ErrInvalidID
	}
	return ID(s), nil
}

func (id ID) String() string {
	return string(id)
}

func (id ID) IsEmpty() bool {
	return id == ""
}

type WalletAddress string

var walletAddressRegex = regexp.MustCompile(`^0x[a-fA-F0-9]{40}$`)

func NewWalletAddress(s string) (WalletAddress, error) {
	if !walletAddressRegex.MatchString(s) {
		return "", ErrInvalidAddress
	}
	return WalletAddress(strings.ToLower(s)), nil
}

func (w WalletAddress) String() string {
	return string(w)
}

func (w WalletAddress) IsEmpty() bool {
	return w == ""
}

func (w WalletAddress) Equals(other WalletAddress) bool {
	return strings.EqualFold(string(w), string(other))
}

type Metadata map[string]any

func NewMetadata() Metadata {
	return make(Metadata)
}

func (m Metadata) Get(key string) (any, bool) {
	if m == nil {
		return nil, false
	}
	v, ok := m[key]
	return v, ok
}

func (m Metadata) GetString(key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func (m Metadata) GetInt(key string) int {
	if v, ok := m[key]; ok {
		switch n := v.(type) {
		case int:
			return n
		case int64:
			return int(n)
		case float64:
			return int(n)
		}
	}
	return 0
}

func (m Metadata) Set(key string, value any) Metadata {
	if m == nil {
		m = make(Metadata)
	}
	m[key] = value
	return m
}

type Timestamp struct {
	time.Time
}

func Now() Timestamp {
	return Timestamp{time.Now().UTC()}
}

func NewTimestamp(t time.Time) Timestamp {
	return Timestamp{t.UTC()}
}

func ParseTimestamp(s string) (Timestamp, error) {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return Timestamp{}, fmt.Errorf("invalid timestamp format: %w", err)
	}
	return NewTimestamp(t), nil
}

func (ts Timestamp) String() string {
	return ts.Format(time.RFC3339)
}

func (ts Timestamp) IsZero() bool {
	return ts.Time.IsZero()
}

func (ts Timestamp) Before(other Timestamp) bool {
	return ts.Time.Before(other.Time)
}

func (ts Timestamp) After(other Timestamp) bool {
	return ts.Time.After(other.Time)
}
