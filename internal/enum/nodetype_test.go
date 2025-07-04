package enum

import (
	"strings"
	"testing"
)

func TestNodeType_Constants(t *testing.T) {
	// Test that all constants are properly defined
	tests := []struct {
		name     string
		nodeType NodeType
		expected int64
	}{
		{"NodeTypeMachine", NodeTypeMachine, 0},
		{"NodeTypeDevice", NodeTypeDevice, 1},
		{"NodeTypeRackController", NodeTypeRackController, 2},
		{"NodeTypeRegionRackController", NodeTypeRegionRackController, 3},
		{"NodeTypeRegionController", NodeTypeRegionController, 4},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if int64(tt.nodeType) != tt.expected {
				t.Errorf("Expected %s to be %d, got %d", tt.name, tt.expected, int64(tt.nodeType))
			}
		})
	}
}

func TestNodeType_String(t *testing.T) {
	tests := []struct {
		name     string
		nodeType NodeType
		expected string
	}{
		{"Machine", NodeTypeMachine, "Machine"},
		{"Device", NodeTypeDevice, "Device"},
		{"RackController", NodeTypeRackController, "Rack Controller"},
		{"RegionRackController", NodeTypeRegionRackController, "Region & Rack Controller"},
		{"RegionController", NodeTypeRegionController, "Region Controller"},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.nodeType.String()
			if result != tt.expected {
				t.Errorf("Expected %s.String() to be %q, got %q", tt.name, tt.expected, result)
			}
		})
	}
}

func TestNodeType_String_Invalid(t *testing.T) {
	// Test invalid NodeType values
	tests := []struct {
		name     string
		nodeType NodeType
		contains string
	}{
		{"negative", NodeType(-1), "NodeType(-1)"},
		{"large", NodeType(999), "NodeType(999)"},
		{"out_of_range", NodeType(10), "NodeType(10)"},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.nodeType.String()
			if !strings.Contains(result, tt.contains) {
				t.Errorf("Expected %s.String() to contain %q, got %q", tt.name, tt.contains, result)
			}
		})
	}
}

func TestNodeType_IsValid(t *testing.T) {
	// Test valid NodeTypes
	validTypes := []NodeType{
		NodeTypeMachine,
		NodeTypeDevice,
		NodeTypeRackController,
		NodeTypeRegionRackController,
		NodeTypeRegionController,
	}
	
	for _, nodeType := range validTypes {
		t.Run("valid_"+nodeType.String(), func(t *testing.T) {
			if !nodeType.IsValid() {
				t.Errorf("Expected %s to be valid", nodeType.String())
			}
		})
	}
	
	// Test invalid NodeTypes
	invalidTypes := []NodeType{
		NodeType(-1),
		NodeType(999),
		NodeType(10),
		NodeType(100),
	}
	
	for _, nodeType := range invalidTypes {
		t.Run("invalid", func(t *testing.T) {
			if nodeType.IsValid() {
				t.Errorf("Expected NodeType(%d) to be invalid", int64(nodeType))
			}
		})
	}
}

func TestParseNodeType(t *testing.T) {
	// Test valid parsing
	tests := []struct {
		name     string
		input    string
		expected NodeType
		wantErr  bool
	}{
		{"Machine", "Machine", NodeTypeMachine, false},
		{"Device", "Device", NodeTypeDevice, false},
		{"RackController", "Rack Controller", NodeTypeRackController, false},
		{"RegionRackController", "Region & Rack Controller", NodeTypeRegionRackController, false},
		{"RegionController", "Region Controller", NodeTypeRegionController, false},
		{"Invalid", "Invalid", NodeType(0), true},
		{"Empty", "", NodeType(0), true},
		{"CaseSensitive", "machine", NodeType(0), true},
		{"ExtraSpaces", " Machine ", NodeType(0), true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseNodeType(tt.input)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error for input %q, got none", tt.input)
				}
				if err != nil && !strings.Contains(err.Error(), "not a valid NodeType") {
					t.Errorf("Expected error to contain 'not a valid NodeType', got: %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error for input %q, got: %v", tt.input, err)
				}
				if result != tt.expected {
					t.Errorf("Expected ParseNodeType(%q) to be %s, got %s", tt.input, tt.expected.String(), result.String())
				}
			}
		})
	}
}

func TestNodeType_MarshalText(t *testing.T) {
	tests := []struct {
		name     string
		nodeType NodeType
		expected string
	}{
		{"Machine", NodeTypeMachine, "Machine"},
		{"Device", NodeTypeDevice, "Device"},
		{"RackController", NodeTypeRackController, "Rack Controller"},
		{"RegionRackController", NodeTypeRegionRackController, "Region & Rack Controller"},
		{"RegionController", NodeTypeRegionController, "Region Controller"},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.nodeType.MarshalText()
			if err != nil {
				t.Errorf("Expected no error for MarshalText(), got: %v", err)
			}
			if string(result) != tt.expected {
				t.Errorf("Expected MarshalText() to be %q, got %q", tt.expected, string(result))
			}
		})
	}
}

func TestNodeType_UnmarshalText(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected NodeType
		wantErr  bool
	}{
		{"Machine", "Machine", NodeTypeMachine, false},
		{"Device", "Device", NodeTypeDevice, false},
		{"RackController", "Rack Controller", NodeTypeRackController, false},
		{"RegionRackController", "Region & Rack Controller", NodeTypeRegionRackController, false},
		{"RegionController", "Region Controller", NodeTypeRegionController, false},
		{"Invalid", "Invalid", NodeType(0), true},
		{"Empty", "", NodeType(0), true},
		{"CaseSensitive", "machine", NodeType(0), true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result NodeType
			err := result.UnmarshalText([]byte(tt.input))
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error for UnmarshalText(%q), got none", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error for UnmarshalText(%q), got: %v", tt.input, err)
				}
				if result != tt.expected {
					t.Errorf("Expected UnmarshalText(%q) to be %s, got %s", tt.input, tt.expected.String(), result.String())
				}
			}
		})
	}
}

func TestNodeTypeNames(t *testing.T) {
	names := NodeTypeNames()
	
	expectedNames := []string{
		"Machine",
		"Device",
		"Rack Controller",
		"Region & Rack Controller",
		"Region Controller",
	}
	
	if len(names) != len(expectedNames) {
		t.Errorf("Expected %d names, got %d", len(expectedNames), len(names))
	}
	
	for i, expected := range expectedNames {
		if i >= len(names) {
			t.Errorf("Missing name at index %d: %s", i, expected)
			continue
		}
		if names[i] != expected {
			t.Errorf("Expected name at index %d to be %q, got %q", i, expected, names[i])
		}
	}
	
	// Test that modifying the returned slice doesn't affect the original
	originalFirstName := names[0]
	names[0] = "Modified"
	newNames := NodeTypeNames()
	if newNames[0] != originalFirstName {
		t.Error("NodeTypeNames() should return a copy, not the original slice")
	}
}

func TestNodeType_Values(t *testing.T) {
	// Test Values method for each NodeType
	nodeTypes := []NodeType{
		NodeTypeMachine,
		NodeTypeDevice,
		NodeTypeRackController,
		NodeTypeRegionRackController,
		NodeTypeRegionController,
	}
	
	for _, nodeType := range nodeTypes {
		t.Run("values_"+nodeType.String(), func(t *testing.T) {
			values := nodeType.Values()
			
			expectedValues := []string{
				"Machine",
				"Device",
				"Rack Controller",
				"Region & Rack Controller",
				"Region Controller",
			}
			
			if len(values) != len(expectedValues) {
				t.Errorf("Expected %d values, got %d", len(expectedValues), len(values))
			}
			
			for i, expected := range expectedValues {
				if i >= len(values) {
					t.Errorf("Missing value at index %d: %s", i, expected)
					continue
				}
				if values[i] != expected {
					t.Errorf("Expected value at index %d to be %q, got %q", i, expected, values[i])
				}
			}
		})
	}
}

func TestErrInvalidNodeType(t *testing.T) {
	// Test that the error contains expected content
	errMsg := ErrInvalidNodeType.Error()
	
	if !strings.Contains(errMsg, "not a valid NodeType") {
		t.Errorf("Expected error message to contain 'not a valid NodeType', got: %s", errMsg)
	}
	
	// Should contain all valid names
	expectedNames := []string{
		"Machine",
		"Device", 
		"Rack Controller",
		"Region & Rack Controller",
		"Region Controller",
	}
	
	for _, name := range expectedNames {
		if !strings.Contains(errMsg, name) {
			t.Errorf("Expected error message to contain %q, got: %s", name, errMsg)
		}
	}
}

// Test round-trip: Marshal then Unmarshal
func TestNodeType_RoundTrip(t *testing.T) {
	nodeTypes := []NodeType{
		NodeTypeMachine,
		NodeTypeDevice,
		NodeTypeRackController,
		NodeTypeRegionRackController,
		NodeTypeRegionController,
	}
	
	for _, original := range nodeTypes {
		t.Run("roundtrip_"+original.String(), func(t *testing.T) {
			// Marshal
			marshaled, err := original.MarshalText()
			if err != nil {
				t.Fatalf("MarshalText failed: %v", err)
			}
			
			// Unmarshal
			var unmarshaled NodeType
			err = unmarshaled.UnmarshalText(marshaled)
			if err != nil {
				t.Fatalf("UnmarshalText failed: %v", err)
			}
			
			// Compare
			if unmarshaled != original {
				t.Errorf("Round-trip failed: original=%s, unmarshaled=%s", 
					original.String(), unmarshaled.String())
			}
		})
	}
}

// Test concurrent access safety
func TestNodeType_ConcurrentAccess(t *testing.T) {
	done := make(chan bool, 20)
	
	// Run multiple goroutines that use the NodeType methods
	for i := 0; i < 20; i++ {
		go func(id int) {
			defer func() { done <- true }()
			
			nodeType := NodeType(id % 5) // Will create valid and invalid types
			
			// Call various methods concurrently
			_ = nodeType.String()
			_ = nodeType.IsValid()
			_ = nodeType.Values()
			_, _ = nodeType.MarshalText()
			
			var nt NodeType
			_ = nt.UnmarshalText([]byte("Machine"))
			
			_, _ = ParseNodeType("Device")
			_ = NodeTypeNames()
		}(i)
	}
	
	// Wait for all goroutines to complete
	for i := 0; i < 20; i++ {
		<-done
	}
}

// Benchmark tests
func BenchmarkNodeType_String(b *testing.B) {
	nodeType := NodeTypeMachine
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = nodeType.String()
	}
}

func BenchmarkNodeType_IsValid(b *testing.B) {
	nodeType := NodeTypeMachine
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = nodeType.IsValid()
	}
}

func BenchmarkParseNodeType(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ParseNodeType("Machine")
	}
}

func BenchmarkNodeType_MarshalText(b *testing.B) {
	nodeType := NodeTypeMachine
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = nodeType.MarshalText()
	}
}

func BenchmarkNodeType_UnmarshalText(b *testing.B) {
	text := []byte("Machine")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var nt NodeType
		_ = nt.UnmarshalText(text)
	}
}

func BenchmarkNodeTypeNames(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NodeTypeNames()
	}
}
