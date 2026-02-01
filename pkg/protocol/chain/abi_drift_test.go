package chain

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

func TestFlemingAnchor_ABIJSONMatchesGeneratedBinding(t *testing.T) {
	t.Helper()

	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}

	abiPath := filepath.Join(filepath.Dir(thisFile), "abi", "FlemingAnchor.abi.json")
	raw, err := os.ReadFile(abiPath)
	if err != nil {
		t.Fatalf("read ABI json %q: %v", abiPath, err)
	}

	fileABI, err := abi.JSON(bytes.NewReader(raw))
	if err != nil {
		t.Fatalf("parse ABI json %q: %v", abiPath, err)
	}

	genABI, err := abi.JSON(strings.NewReader(FlemingAnchorMetaData.ABI))
	if err != nil {
		t.Fatalf("parse generated ABI: %v", err)
	}

	assertSameMethods(t, "methods", fileABI.Methods, genABI.Methods)
	assertSameEvents(t, "events", fileABI.Events, genABI.Events)

	// go-ethereum's ABI parser supports custom errors; keep these in lockstep too.
	// (If this ever becomes unavailable in the ABI type, this assertion can be dropped.)
	assertSameErrors(t, "errors", fileABI.Errors, genABI.Errors)
}

func assertSameMethods(t *testing.T, label string, a, b map[string]abi.Method) {
	t.Helper()
	if len(a) != len(b) {
		t.Fatalf("ABI %s mismatch: file=%d generated=%d", label, len(a), len(b))
	}
	for name, am := range a {
		bm, ok := b[name]
		if !ok {
			t.Fatalf("ABI %s mismatch: %q missing from generated binding", label, name)
		}
		if am.Sig != bm.Sig {
			t.Fatalf("ABI %s mismatch: %q signature file=%q generated=%q", label, name, am.Sig, bm.Sig)
		}
		if !bytes.Equal(am.ID, bm.ID) {
			t.Fatalf("ABI %s mismatch: %q selector differs (file=%#x generated=%#x)", label, name, am.ID, bm.ID)
		}
	}
	for name := range b {
		if _, ok := a[name]; !ok {
			t.Fatalf("ABI %s mismatch: %q present in generated binding but missing from file", label, name)
		}
	}
}

func assertSameEvents(t *testing.T, label string, a, b map[string]abi.Event) {
	t.Helper()
	if len(a) != len(b) {
		t.Fatalf("ABI %s mismatch: file=%d generated=%d", label, len(a), len(b))
	}
	for name, ae := range a {
		be, ok := b[name]
		if !ok {
			t.Fatalf("ABI %s mismatch: %q missing from generated binding", label, name)
		}
		if ae.Sig != be.Sig {
			t.Fatalf("ABI %s mismatch: %q signature file=%q generated=%q", label, name, ae.Sig, be.Sig)
		}
		if ae.ID != be.ID {
			t.Fatalf("ABI %s mismatch: %q topic0 differs (file=%s generated=%s)", label, name, ae.ID, be.ID)
		}
	}
	for name := range b {
		if _, ok := a[name]; !ok {
			t.Fatalf("ABI %s mismatch: %q present in generated binding but missing from file", label, name)
		}
	}
}

func assertSameErrors(t *testing.T, label string, a, b map[string]abi.Error) {
	t.Helper()
	if len(a) != len(b) {
		t.Fatalf("ABI %s mismatch: file=%d generated=%d", label, len(a), len(b))
	}
	for name, ae := range a {
		be, ok := b[name]
		if !ok {
			t.Fatalf("ABI %s mismatch: %q missing from generated binding", label, name)
		}
		if ae.Sig != be.Sig {
			t.Fatalf("ABI %s mismatch: %q signature file=%q generated=%q", label, name, ae.Sig, be.Sig)
		}
		if ae.ID != be.ID {
			t.Fatalf("ABI %s mismatch: %q selector differs (file=%#x generated=%#x)", label, name, ae.ID, be.ID)
		}
	}
	for name := range b {
		if _, ok := a[name]; !ok {
			t.Fatalf("ABI %s mismatch: %q present in generated binding but missing from file", label, name)
		}
	}
}

