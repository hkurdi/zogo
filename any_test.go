package zogo

import (
	"testing"
)

// Test Any accepts strings
func TestAnyString(t *testing.T) {
	schema := Any()

	result := schema.Parse("hello")
	if !result.Ok {
		t.Error("Expected Any to accept string")
	}
	if result.Value != "hello" {
		t.Errorf("Expected 'hello', got %v", result.Value)
	}
}

// Test Any accepts numbers
func TestAnyNumber(t *testing.T) {
	schema := Any()

	result := schema.Parse(42)
	if !result.Ok {
		t.Error("Expected Any to accept number")
	}
	if result.Value != 42 {
		t.Errorf("Expected 42, got %v", result.Value)
	}
}

// Test Any accepts booleans
func TestAnyBoolean(t *testing.T) {
	schema := Any()

	result := schema.Parse(true)
	if !result.Ok {
		t.Error("Expected Any to accept boolean")
	}
	if result.Value != true {
		t.Errorf("Expected true, got %v", result.Value)
	}
}

// Test Any accepts objects
func TestAnyObject(t *testing.T) {
	schema := Any()

	obj := map[string]interface{}{"key": "value"}
	result := schema.Parse(obj)
	if !result.Ok {
		t.Error("Expected Any to accept object")
	}
}

// Test Any accepts arrays
func TestAnyArray(t *testing.T) {
	schema := Any()

	arr := []interface{}{1, 2, 3}
	result := schema.Parse(arr)
	if !result.Ok {
		t.Error("Expected Any to accept array")
	}
}

// Test Any accepts nil by default
func TestAnyNilDefault(t *testing.T) {
	schema := Any()

	result := schema.Parse(nil)
	if !result.Ok {
		t.Error("Expected Any to accept nil by default")
	}
	if result.Value != nil {
		t.Errorf("Expected nil, got %v", result.Value)
	}
}

// Test Any with Optional
func TestAnyOptional(t *testing.T) {
	schema := Any().Optional()

	result := schema.Parse(nil)
	if !result.Ok {
		t.Error("Expected Any().Optional() to accept nil")
	}

	result = schema.Parse("value")
	if !result.Ok {
		t.Error("Expected Any().Optional() to accept value")
	}
}

// Test Any with Nullable
func TestAnyNullable(t *testing.T) {
	schema := Any().Nullable()

	result := schema.Parse(nil)
	if !result.Ok {
		t.Error("Expected Any().Nullable() to accept nil")
	}
}

// Test Any with Required
func TestAnyRequired(t *testing.T) {
	schema := Any().Required()

	// nil should fail
	result := schema.Parse(nil)
	if result.Ok {
		t.Error("Expected Any().Required() to reject nil")
	}

	// Any non-nil value should pass
	result = schema.Parse("value")
	if !result.Ok {
		t.Error("Expected Any().Required() to accept value")
	}

	result = schema.Parse(0)
	if !result.Ok {
		t.Error("Expected Any().Required() to accept 0")
	}

	result = schema.Parse(false)
	if !result.Ok {
		t.Error("Expected Any().Required() to accept false")
	}

	result = schema.Parse("")
	if !result.Ok {
		t.Error("Expected Any().Required() to accept empty string")
	}
}

// Test Any in object
func TestAnyInObject(t *testing.T) {
	schema := Object(Schema{
		"name":     String(),
		"metadata": Any(),
	})

	data := map[string]interface{}{
		"name": "John",
		"metadata": map[string]interface{}{
			"nested": "value",
			"number": 42,
		},
	}

	result := schema.Parse(data)
	if !result.Ok {
		t.Errorf("Expected object with Any field to pass. Errors: %v", result.Errors)
	}
}

// Test Any in array
func TestAnyInArray(t *testing.T) {
	schema := Array(Any())

	data := []interface{}{"string", 42, true, map[string]interface{}{"key": "value"}, nil}

	result := schema.Parse(data)
	if !result.Ok {
		t.Errorf("Expected array of Any to accept mixed types. Errors: %v", result.Errors)
	}
}

// Test Any accepts zero values
func TestAnyZeroValues(t *testing.T) {
	schema := Any()

	result := schema.Parse(0)
	if !result.Ok {
		t.Error("Expected Any to accept 0")
	}

	result = schema.Parse("")
	if !result.Ok {
		t.Error("Expected Any to accept empty string")
	}

	result = schema.Parse(false)
	if !result.Ok {
		t.Error("Expected Any to accept false")
	}
}
