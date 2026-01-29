package types

import (
	"sync"
	"testing"
)

func TestNewTypeRegistry(t *testing.T) {
	reg := NewTypeRegistry[string]()

	if reg == nil {
		t.Error("NewTypeRegistry() returned nil")
	}

	// Should start empty
	if len(reg.ValidTypes()) != 0 {
		t.Error("NewTypeRegistry() should start with no types")
	}
}

func TestRegistry_Register(t *testing.T) {
	reg := NewTypeRegistry[string]()

	metadata := TypeMetadata{
		Name:        "TestType",
		Description: "A test type",
		Since:       "1.0.0",
	}

	err := reg.Register("test", metadata)
	if err != nil {
		t.Errorf("Register() error = %v", err)
	}

	if !reg.IsValid("test") {
		t.Error("Register() should make type valid")
	}

	// Duplicate registration should error
	err = reg.Register("test", metadata)
	if err == nil {
		t.Error("Register() with duplicate type should error")
	}

	// Auto-generate name if empty
	reg2 := NewTypeRegistry[string]()
	err = reg2.Register("auto", TypeMetadata{Description: "Auto-named"})
	if err != nil {
		t.Errorf("Register() with empty name error = %v", err)
	}
	meta, _ := reg2.GetMetadata("auto")
	if meta.Name == "" {
		t.Error("Register() should auto-generate name if empty")
	}
}

func TestRegistry_IsValid(t *testing.T) {
	reg := NewTypeRegistry[string]()

	reg.Register("valid", TypeMetadata{Name: "Valid"})
	reg.Register("deprecated", TypeMetadata{Name: "Deprecated", Deprecated: true})

	if !reg.IsValid("valid") {
		t.Error("IsValid() should return true for registered type")
	}

	if reg.IsValid("deprecated") {
		t.Error("IsValid() should return false for deprecated type")
	}

	if reg.IsValid("unregistered") {
		t.Error("IsValid() should return false for unregistered type")
	}
}

func TestRegistry_ValidTypes(t *testing.T) {
	reg := NewTypeRegistry[string]()

	reg.Register("first", TypeMetadata{Name: "First"})
	reg.Register("second", TypeMetadata{Name: "Second"})
	reg.Register("third", TypeMetadata{Name: "Third"})

	valid := reg.ValidTypes()
	if len(valid) != 3 {
		t.Errorf("ValidTypes() returned %d types, want 3", len(valid))
	}

	// Should preserve registration order
	if valid[0] != "first" || valid[1] != "second" || valid[2] != "third" {
		t.Error("ValidTypes() should preserve registration order")
	}
}

func TestRegistry_GetMetadata(t *testing.T) {
	reg := NewTypeRegistry[string]()

	metadata := TypeMetadata{
		Name:        "TestType",
		Description: "A test type",
		Since:       "1.0.0",
		Deprecated:  false,
	}

	reg.Register("test", metadata)

	meta, ok := reg.GetMetadata("test")
	if !ok {
		t.Error("GetMetadata() should find registered type")
	}
	if meta.Name != "TestType" {
		t.Errorf("GetMetadata() Name = %v, want TestType", meta.Name)
	}
	if meta.Description != "A test type" {
		t.Errorf("GetMetadata() Description = %v, want A test type", meta.Description)
	}

	_, ok = reg.GetMetadata("nonexistent")
	if ok {
		t.Error("GetMetadata() should not find unregistered type")
	}
}

func TestRegisterBatch(t *testing.T) {
	reg := NewTypeRegistry[string]()

	types := map[string]TypeMetadata{
		"type1": {Name: "Type 1"},
		"type2": {Name: "Type 2"},
		"type3": {Name: "Type 3"},
	}

	err := RegisterBatch(reg, types)
	if err != nil {
		t.Errorf("RegisterBatch() error = %v", err)
	}

	if len(reg.ValidTypes()) != 3 {
		t.Errorf("RegisterBatch() registered %d types, want 3", len(reg.ValidTypes()))
	}

	// Duplicate in batch should error
	types["type1"] = TypeMetadata{Name: "Duplicate"}
	err = RegisterBatch(reg, types)
	if err == nil {
		t.Error("RegisterBatch() with duplicate should error")
	}
}

func TestRegistry_ThreadSafety(t *testing.T) {
	reg := NewTypeRegistry[string]()

	var wg sync.WaitGroup
	numGoroutines := 10
	typesPerGoroutine := 10

	// Concurrent registration
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(prefix int) {
			defer wg.Done()
			for j := 0; j < typesPerGoroutine; j++ {
				typeName := string(rune('a' + prefix)) + string(rune('0'+j))
				reg.Register(typeName, TypeMetadata{Name: typeName})
			}
		}(i)
	}

	// Concurrent reads
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				reg.IsValid("test")
				reg.ValidTypes()
			}
		}()
	}

	wg.Wait()

	// Should have registered all types
	valid := reg.ValidTypes()
	expectedCount := numGoroutines * typesPerGoroutine
	if len(valid) != expectedCount {
		t.Errorf("ThreadSafety test: registered %d types, want %d", len(valid), expectedCount)
	}
}
