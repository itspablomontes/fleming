package common

import (
	"encoding/hex"
	"strings"
)

// HexToBytes converts a hex string (potentially with 0x prefix) to a byte slice.
func HexToBytes(s string) ([]byte, error) {
	s = strings.TrimPrefix(s, "0x")
	return hex.DecodeString(s)
}

// BytesToHex converts a byte slice to a hex string with 0x prefix.
func BytesToHex(b []byte) string {
	return "0x" + hex.EncodeToString(b)
}
