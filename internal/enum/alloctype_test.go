package enum

import (
	"encoding"
	"reflect"
	"testing"
)

func TestAllocTypeStringAndIsValid(t *testing.T) {
	tests := []struct {
		typ      AllocType
		expected string
		valid    bool
	}{
		{AllocTypeAutomatic, "Automatic", true},
		{AllocTypeSticky, "Sticky", true},
		{AllocTypeUserReserved, "User Reserved", true},
		{AllocTypeDHCP, "DHCP", true},
		{AllocTypeDiscovered, "Discovered", true},
		{AllocType(99), "AllocType(99)", false},
	}

	for _, tt := range tests {
		if got := tt.typ.String(); got != tt.expected {
			t.Errorf("AllocType(%d).String() = %q, want %q", tt.typ, got, tt.expected)
		}
		if ok := tt.typ.IsValid(); ok != tt.valid {
			t.Errorf("AllocType(%d).IsValid() = %v, want %v", tt.typ, ok, tt.valid)
		}
	}
}

func TestParseAllocType(t *testing.T) {
	validNames := AllocTypeNames()
	for _, name := range validNames {
		typ, err := ParseAllocType(name)
		if err != nil {
			t.Fatalf("ParseAllocType(%q) unexpected error: %v", name, err)
		}
		if got := typ.String(); got != name {
			t.Errorf("Parsed type .String() = %q, want %q", got, name)
		}
	}

	if _, err := ParseAllocType("InvalidName"); err == nil {
		t.Error("ParseAllocType(\"InvalidName\") expected error, got nil")
	}
}

func TestAllocTypeMarshalUnmarshal(t *testing.T) {
	var _ encoding.TextMarshaler = AllocTypeAutomatic // compile‑time 介面驗證
	var _ encoding.TextUnmarshaler = (*AllocType)(nil)

	for _, typ := range []AllocType{
		AllocTypeAutomatic,
		AllocTypeSticky,
		AllocTypeUserReserved,
		AllocTypeDHCP,
		AllocTypeDiscovered,
	} {
		// Marshal
		b, err := typ.MarshalText()
		if err != nil {
			t.Fatalf("MarshalText(%v) error: %v", typ, err)
		}
		if string(b) != typ.String() {
			t.Errorf("MarshalText result = %q, want %q", string(b), typ.String())
		}

		// Unmarshal
		var dst AllocType
		if err := dst.UnmarshalText(b); err != nil {
			t.Fatalf("UnmarshalText(%q) error: %v", b, err)
		}
		if dst != typ {
			t.Errorf("UnmarshalText result = %v, want %v", dst, typ)
		}
	}
}

func TestAllocTypeUnmarshalInvalid(t *testing.T) {
	var a AllocType
	err := a.UnmarshalText([]byte("NotExist"))
	if err == nil {
		t.Error("UnmarshalText with invalid name should return error")
	}
}

func TestAllocTypeRoundTripInterface(t *testing.T) {
	original := AllocTypeUserReserved
	marshaled, _ := original.MarshalText()

	var iface encoding.TextUnmarshaler = new(AllocType)
	if err := iface.UnmarshalText(marshaled); err != nil {
		t.Fatalf("interface UnmarshalText failed: %v", err)
	}
	if !reflect.DeepEqual(original, *(iface.(*AllocType))) {
		t.Errorf("round‑trip mismatch: got %v, want %v", *(iface.(*AllocType)), original)
	}
}
