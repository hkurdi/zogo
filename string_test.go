package zogo

import "testing"

// Test basic string validation
func TestStringBasic(t *testing.T) {
	schema := String()
	result := schema.Parse("hello")

	if !result.Ok {
		t.Error("Expected valid string to pass")
	}

	if result.Value != "hello" {
		t.Errorf("Expected 'hello', got %v", result.Value)
	}
}

// Test non-string input
func TestStringInvalidType(t *testing.T) {
	schema := String()
	result := schema.Parse(123)

	if result.Ok {
		t.Error("Expected number to fail string validation")
	}

	if len(result.Errors) == 0 {
		t.Error("Expected error for invalid type")
	}
}

// Test nil value
func TestStringNil(t *testing.T) {
	schema := String()
	result := schema.Parse(nil)

	if result.Ok {
		t.Error("Expected nil to fail")
	}
}

// Test Min length
func TestStringMin(t *testing.T) {
	schema := String().Min(5)

	// Should pass
	result := schema.Parse("hello")
	if !result.Ok {
		t.Error("Expected 5 char string to pass Min(5)")
	}

	// Should fail
	result = schema.Parse("hi")
	if result.Ok {
		t.Error("Expected 2 char string to fail Min(5)")
	}
}

// Test Max length
func TestStringMax(t *testing.T) {
	schema := String().Max(5)

	// Should pass
	result := schema.Parse("hello")
	if !result.Ok {
		t.Error("Expected 5 char string to pass Max(5)")
	}

	// Should fail
	result = schema.Parse("hello world")
	if result.Ok {
		t.Error("Expected 11 char string to fail Max(5)")
	}
}

// Test exact Length
func TestStringLength(t *testing.T) {
	schema := String().Length(5)

	// Should pass
	result := schema.Parse("hello")
	if !result.Ok {
		t.Error("Expected 5 char string to pass Length(5)")
	}

	// Should fail - too short
	result = schema.Parse("hi")
	if result.Ok {
		t.Error("Expected 2 char string to fail Length(5)")
	}

	// Should fail - too long
	result = schema.Parse("hello world")
	if result.Ok {
		t.Error("Expected 11 char string to fail Length(5)")
	}
}

// Test chaining Min and Max
func TestStringMinMax(t *testing.T) {
	schema := String().Min(3).Max(10)

	// Should pass
	result := schema.Parse("hello")
	if !result.Ok {
		t.Error("Expected 5 char string to pass Min(3).Max(10)")
	}

	// Should fail - too short
	result = schema.Parse("hi")
	if result.Ok {
		t.Error("Expected 2 char string to fail Min(3)")
	}

	// Should fail - too long
	result = schema.Parse("hello world!")
	if result.Ok {
		t.Error("Expected 12 char string to fail Max(10)")
	}
}
