package zogo

import (
	"testing"
)

// Test basic tuple with three elements
func TestTupleBasic(t *testing.T) {
	schema := Tuple(String(), Number(), Boolean())

	// Should pass - correct types in order
	result := schema.Parse([]interface{}{"hello", 42, true})
	if !result.Ok {
		t.Errorf("Expected valid tuple to pass. Errors: %v", result.Errors)
	}

	resultArr := result.Value.([]interface{})
	if len(resultArr) != 3 {
		t.Errorf("Expected 3 elements, got %d", len(resultArr))
	}
	if resultArr[0] != "hello" || resultArr[1] != float64(42) || resultArr[2] != true {
		t.Error("Expected tuple values to be preserved")
	}
}

// Test tuple with two elements
func TestTuplePair(t *testing.T) {
	schema := Tuple(String(), Number())

	// Should pass
	result := schema.Parse([]interface{}{"x", 10})
	if !result.Ok {
		t.Error("Expected valid pair to pass")
	}
}

// Test single element tuple
func TestTupleSingle(t *testing.T) {
	schema := Tuple(String())

	// Should pass
	result := schema.Parse([]interface{}{"hello"})
	if !result.Ok {
		t.Error("Expected single-element tuple to pass")
	}
}

// Test wrong type at position
func TestTupleWrongType(t *testing.T) {
	schema := Tuple(String(), Number(), Boolean())

	// Wrong type at position 0
	result := schema.Parse([]interface{}{123, 42, true})
	if result.Ok {
		t.Error("Expected wrong type at position 0 to fail")
	}

	// Wrong type at position 1
	result = schema.Parse([]interface{}{"hello", "world", true})
	if result.Ok {
		t.Error("Expected wrong type at position 1 to fail")
	}

	// Wrong type at position 2
	result = schema.Parse([]interface{}{"hello", 42, "not bool"})
	if result.Ok {
		t.Error("Expected wrong type at position 2 to fail")
	}
}

// Test wrong length - too short
func TestTupleTooShort(t *testing.T) {
	schema := Tuple(String(), Number(), Boolean())

	// Only 2 elements
	result := schema.Parse([]interface{}{"hello", 42})
	if result.Ok {
		t.Error("Expected too-short tuple to fail")
	}

	// Only 1 element
	result = schema.Parse([]interface{}{"hello"})
	if result.Ok {
		t.Error("Expected too-short tuple to fail")
	}

	// Empty
	result = schema.Parse([]interface{}{})
	if result.Ok {
		t.Error("Expected empty array to fail")
	}
}

// Test wrong length - too long
func TestTupleTooLong(t *testing.T) {
	schema := Tuple(String(), Number())

	// 3 elements when expecting 2
	result := schema.Parse([]interface{}{"hello", 42, true})
	if result.Ok {
		t.Error("Expected too-long tuple to fail")
	}
}

// Test with validators that have constraints
func TestTupleWithConstraints(t *testing.T) {
	schema := Tuple(
		String().Email(),
		Number().Min(18),
		Boolean(),
	)

	// Should pass
	result := schema.Parse([]interface{}{"user@example.com", 25, true})
	if !result.Ok {
		t.Error("Expected valid constrained tuple to pass")
	}

	// Invalid email
	result = schema.Parse([]interface{}{"notanemail", 25, true})
	if result.Ok {
		t.Error("Expected invalid email to fail")
	}

	// Number too small
	result = schema.Parse([]interface{}{"user@example.com", 15, true})
	if result.Ok {
		t.Error("Expected number < 18 to fail")
	}
}

// Test Rest validator
func TestTupleRest(t *testing.T) {
	// First two must be String and Number, rest can be any Number
	schema := Tuple(String(), Number()).Rest(Number())

	// Exact length (no rest elements)
	result := schema.Parse([]interface{}{"hello", 42})
	if !result.Ok {
		t.Error("Expected exact length to pass with Rest")
	}

	// With rest elements
	result = schema.Parse([]interface{}{"hello", 42, 100, 200, 300})
	if !result.Ok {
		t.Error("Expected tuple with valid rest elements to pass")
	}

	// Invalid rest element
	result = schema.Parse([]interface{}{"hello", 42, "not a number"})
	if result.Ok {
		t.Error("Expected invalid rest element to fail")
	}
}

// Test Rest with constraints
func TestTupleRestConstraints(t *testing.T) {
	schema := Tuple(String()).Rest(Number().Positive())

	// All rest elements positive
	result := schema.Parse([]interface{}{"name", 1, 2, 3})
	if !result.Ok {
		t.Error("Expected all positive rest elements to pass")
	}

	// Rest element not positive
	result = schema.Parse([]interface{}{"name", 1, -2, 3})
	if result.Ok {
		t.Error("Expected negative rest element to fail")
	}
}

// Test nil value
func TestTupleNil(t *testing.T) {
	schema := Tuple(String(), Number())

	result := schema.Parse(nil)
	if result.Ok {
		t.Error("Expected nil to fail")
	}
}

// Test Optional
func TestTupleOptional(t *testing.T) {
	schema := Tuple(String(), Number()).Optional()

	// nil should pass
	result := schema.Parse(nil)
	if !result.Ok {
		t.Error("Expected nil to pass with Optional()")
	}

	// Valid tuple should pass
	result = schema.Parse([]interface{}{"hello", 42})
	if !result.Ok {
		t.Error("Expected valid tuple to pass with Optional()")
	}
}

// Test Nullable
func TestTupleNullable(t *testing.T) {
	schema := Tuple(String(), Number()).Nullable()

	// nil should pass
	result := schema.Parse(nil)
	if !result.Ok {
		t.Error("Expected nil to pass with Nullable()")
	}
}

// Test Required
func TestTupleRequired(t *testing.T) {
	schema := Tuple(String(), Number()).Required()

	// nil should fail
	result := schema.Parse(nil)
	if result.Ok {
		t.Error("Expected nil to fail with Required()")
	}

	// Valid tuple should pass
	result = schema.Parse([]interface{}{"hello", 42})
	if !result.Ok {
		t.Error("Expected valid tuple to pass with Required()")
	}
}

// Test invalid type (not array)
func TestTupleInvalidType(t *testing.T) {
	schema := Tuple(String(), Number())

	result := schema.Parse("not an array")
	if result.Ok {
		t.Error("Expected string to fail tuple validation")
	}

	result = schema.Parse(42)
	if result.Ok {
		t.Error("Expected number to fail tuple validation")
	}

	result = schema.Parse(map[string]interface{}{"key": "value"})
	if result.Ok {
		t.Error("Expected object to fail tuple validation")
	}
}

// Test error paths
func TestTupleErrorPaths(t *testing.T) {
	schema := Tuple(String().Email(), Number().Min(18))

	// Invalid email at position 0
	result := schema.Parse([]interface{}{"notanemail", 25})
	if result.Ok {
		t.Error("Expected invalid email to fail")
	}
	if len(result.Errors) == 0 || result.Errors[0].Path != "[0]" {
		t.Errorf("Expected error path '[0]', got '%s'", result.Errors[0].Path)
	}

	// Invalid number at position 1
	result = schema.Parse([]interface{}{"user@example.com", 15})
	if result.Ok {
		t.Error("Expected number < 18 to fail")
	}
	if len(result.Errors) == 0 || result.Errors[0].Path != "[1]" {
		t.Errorf("Expected error path '[1]', got '%s'", result.Errors[0].Path)
	}
}

// Test tuple in object
func TestTupleInObject(t *testing.T) {
	schema := Object(Schema{
		"name":   String(),
		"coords": Tuple(Number(), Number()),
	})

	data := map[string]interface{}{
		"name":   "Point A",
		"coords": []interface{}{10.5, 20.3},
	}

	result := schema.Parse(data)
	if !result.Ok {
		t.Errorf("Expected object with tuple to pass. Errors: %v", result.Errors)
	}
}

// Test tuple in array
func TestTupleInArray(t *testing.T) {
	schema := Array(Tuple(String(), Number()))

	data := []interface{}{
		[]interface{}{"a", 1},
		[]interface{}{"b", 2},
		[]interface{}{"c", 3},
	}

	result := schema.Parse(data)
	if !result.Ok {
		t.Errorf("Expected array of tuples to pass. Errors: %v", result.Errors)
	}
}

// Test nested tuples
func TestTupleNested(t *testing.T) {
	schema := Tuple(
		String(),
		Tuple(Number(), Number()),
	)

	data := []interface{}{
		"name",
		[]interface{}{10, 20},
	}

	result := schema.Parse(data)
	if !result.Ok {
		t.Errorf("Expected nested tuple to pass. Errors: %v", result.Errors)
	}
}

// Test empty tuple
func TestTupleEmpty(t *testing.T) {
	schema := Tuple()

	// Should only accept empty array
	result := schema.Parse([]interface{}{})
	if !result.Ok {
		t.Error("Expected empty array to pass empty tuple")
	}

	// Should fail on non-empty
	result = schema.Parse([]interface{}{"anything"})
	if result.Ok {
		t.Error("Expected non-empty array to fail empty tuple")
	}
}

// Test multiple errors
func TestTupleMultipleErrors(t *testing.T) {
	schema := Tuple(String().Email(), Number().Min(18), Boolean())

	// All invalid
	result := schema.Parse([]interface{}{"notanemail", 15, "notbool"})
	if result.Ok {
		t.Error("Expected all invalid to fail")
	}

	// Should have 3 errors
	if len(result.Errors) != 3 {
		t.Errorf("Expected 3 errors, got %d", len(result.Errors))
	}
}

// Test coordinates pattern (common use case)
func TestTupleCoordinates(t *testing.T) {
	coordSchema := Tuple(Number(), Number())

	result := coordSchema.Parse([]interface{}{40.7128, -74.0060})
	if !result.Ok {
		t.Error("Expected coordinates to pass")
	}
}

// Test RGB color pattern (common use case)
func TestTupleRGB(t *testing.T) {
	rgbSchema := Tuple(
		Number().Min(0).Max(255),
		Number().Min(0).Max(255),
		Number().Min(0).Max(255),
	)

	// Valid RGB
	result := rgbSchema.Parse([]interface{}{255, 128, 0})
	if !result.Ok {
		t.Error("Expected valid RGB to pass")
	}

	// Invalid - out of range
	result = rgbSchema.Parse([]interface{}{256, 128, 0})
	if result.Ok {
		t.Error("Expected out-of-range RGB to fail")
	}
}

// Test key-value pair pattern (common use case)
func TestTupleKeyValue(t *testing.T) {
	kvSchema := Tuple(String(), Any())

	result := kvSchema.Parse([]interface{}{"name", "John"})
	if !result.Ok {
		t.Error("Expected key-value pair to pass")
	}

	result = kvSchema.Parse([]interface{}{"age", 30})
	if !result.Ok {
		t.Error("Expected key-value pair with number to pass")
	}
}
