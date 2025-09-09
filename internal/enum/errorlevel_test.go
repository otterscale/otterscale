package enum

import (
	"encoding"
	"reflect"
	"testing"
)

func TestErrorLevelStringAndIsValid(t *testing.T) {
	tests := []struct {
		level    ErrorLevel
		expected string
		valid    bool
	}{
		{ErrorLevelInfo, "Info", true},
		{ErrorLevelLow, "Low", true},
		{ErrorLevelMedium, "Medium", true},
		{ErrorLevelHigh, "High", true},
		{ErrorLevelCritical, "Critical", true},
		{ErrorLevel(99), "ErrorLevel(99)", false},
	}

	for _, tt := range tests {
		if got := tt.level.String(); got != tt.expected {
			t.Errorf("ErrorLevel(%d).String() = %q, want %q", tt.level, got, tt.expected)
		}
		if ok := tt.level.IsValid(); ok != tt.valid {
			t.Errorf("ErrorLevel(%d).IsValid() = %v, want %v", tt.level, ok, tt.valid)
		}
	}
}

func TestParseErrorLevel(t *testing.T) {
	for _, name := range ErrorLevelNames() {
		level, err := ParseErrorLevel(name)
		if err != nil {
			t.Fatalf("ParseErrorLevel(%q) unexpected error: %v", name, err)
		}
		if level.String() != name {
			t.Errorf("Parsed level .String() = %q, want %q", level.String(), name)
		}
	}

	if _, err := ParseErrorLevel("Unknown"); err == nil {
		t.Error("ParseErrorLevel(\"Unknown\") expected error, got nil")
	}
}

func TestErrorLevelMarshalUnmarshal(t *testing.T) {
	var _ encoding.TextMarshaler = ErrorLevelInfo
	var _ encoding.TextUnmarshaler = (*ErrorLevel)(nil)

	for _, lvl := range []ErrorLevel{
		ErrorLevelInfo,
		ErrorLevelLow,
		ErrorLevelMedium,
		ErrorLevelHigh,
		ErrorLevelCritical,
	} {
		// Marshal
		b, err := lvl.MarshalText()
		if err != nil {
			t.Fatalf("MarshalText(%v) error: %v", lvl, err)
		}
		if string(b) != lvl.String() {
			t.Errorf("MarshalText result = %q, want %q", string(b), lvl.String())
		}

		// Unmarshal
		var dst ErrorLevel
		if err := dst.UnmarshalText(b); err != nil {
			t.Fatalf("UnmarshalText(%q) error: %v", b, err)
		}
		if dst != lvl {
			t.Errorf("UnmarshalText result = %v, want %v", dst, lvl)
		}
	}
}

func TestErrorLevelUnmarshalInvalid(t *testing.T) {
	var lvl ErrorLevel
	if err := lvl.UnmarshalText([]byte("NotExist")); err == nil {
		t.Error("UnmarshalText with invalid name should return error")
	}
}

func TestErrorLevelRoundTripInterface(t *testing.T) {
	orig := ErrorLevelCritical
	data, _ := orig.MarshalText()

	var iface encoding.TextUnmarshaler = new(ErrorLevel)
	if err := iface.UnmarshalText(data); err != nil {
		t.Fatalf("interface UnmarshalText failed: %v", err)
	}
	if !reflect.DeepEqual(orig, *(iface.(*ErrorLevel))) {
		t.Errorf("round‑trip mismatch: got %v, want %v", *(iface.(*ErrorLevel)), orig)
	}
}
