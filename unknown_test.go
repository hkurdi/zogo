package zogo

import (
	"testing"
)

// Test Unknown accepts strings
func TestUnknownString(t *testing.T) {
	schema := Unknown()

	result := schema.Parse("hello")
	if !result.Ok {
		t.Error("Expected Unknown to accept string")
	}
	if result.Value != "hello" {
		t.Errorf("Expected 'hello', got %v", result.Value)
	}
}

// Test Unknown accepts numbers
func TestUnknownNumber(t *testing.T) {
	schema := Unknown()

	result := schema.Parse(42)
	if !result.Ok {
		t.Error("Expected Unknown to accept number")
	}
	if result.Value != 42 {
		t.Errorf("Expected 42, got %v", result.Value)
	}
}

// Test Unknown accepts booleans
func TestUnknownBoolean(t *testing.T) {
	schema := Unknown()

	result := schema.Parse(true)
	if !result.Ok {
		t.Error("Expected Unknown to accept boolean")
	}
	if result.Value != true {
		t.Errorf("Expected true, got %v", result.Value)
	}
}

// Test Unknown accepts objects
func TestUnknownObject(t *testing.T) {
	schema := Unknown()

	obj := map[string]interface{}{"key": "value"}
	result := schema.Parse(obj)
	if !result.Ok {
		t.Error("Expected Unknown to accept object")
	}
}

// Test Unknown accepts arrays
func TestUnknownArray(t *testing.T) {
	schema := Unknown()

	arr := []interface{}{1, 2, 3}
	result := schema.Parse(arr)
	if !result.Ok {
		t.Error("Expected Unknown to accept array")
	}
}

// Test Unknown accepts nil by default
func TestUnknownNilDefault(t *testing.T) {
	schema := Unknown()

	result := schema.Parse(nil)
	if !result.Ok {
		t.Error("Expected Unknown to accept nil by default")
	}
	if result.Value != nil {
		t.Errorf("Expected nil, got %v", result.Value)
	}
}

// Test Unknown with Optional
func TestUnknownOptional(t *testing.T) {
	schema := Unknown().Optional()

	result := schema.Parse(nil)
	if !result.Ok {
		t.Error("Expected Unknown().Optional() to accept nil")
	}

	result = schema.Parse("value")
	if !result.Ok {
		t.Error("Expected Unknown().Optional() to accept value")
	}
}

// Test Unknown with Nullable
func TestUnknownNullable(t *testing.T) {
	schema := Unknown().Nullable()

	result := schema.Parse(nil)
	if !result.Ok {
		t.Error("Expected Unknown().Nullable() to accept nil")
	}
}

// Test Unknown with Required
func TestUnknownRequired(t *testing.T) {
	schema := Unknown().Required()

	// nil should fail
	result := schema.Parse(nil)
	if result.Ok {
		t.Error("Expected Unknown().Required() to reject nil")
	}

	// Any non-nil value should pass
	result = schema.Parse("value")
	if !result.Ok {
		t.Error("Expected Unknown().Required() to accept value")
	}

	result = schema.Parse(0)
	if !result.Ok {
		t.Error("Expected Unknown().Required() to accept 0")
	}

	result = schema.Parse(false)
	if !result.Ok {
		t.Error("Expected Unknown().Required() to accept false")
	}

	result = schema.Parse("")
	if !result.Ok {
		t.Error("Expected Unknown().Required() to accept empty string")
	}
}

// Test Unknown in object
func TestUnknownInObject(t *testing.T) {
	schema := Object(Schema{
		"name": String(),
		"data": Unknown(),
	})

	data := map[string]interface{}{
		"name": "John",
		"data": map[string]interface{}{
			"anything": "goes",
			"here":     123,
		},
	}

	result := schema.Parse(data)
	if !result.Ok {
		t.Errorf("Expected object with Unknown field to pass. Errors: %v", result.Errors)
	}
}

// Test Unknown in array
func TestUnknownInArray(t *testing.T) {
	schema := Array(Unknown())

	data := []interface{}{"string", 42, true, map[string]interface{}{"key": "value"}, nil}

	result := schema.Parse(data)
	if !result.Ok {
		t.Errorf("Expected array of Unknown to accept mixed types. Errors: %v", result.Errors)
	}
}

// Test Unknown accepts zero values
func TestUnknownZeroValues(t *testing.T) {
	schema := Unknown()

	result := schema.Parse(0)
	if !result.Ok {
		t.Error("Expected Unknown to accept 0")
	}

	result = schema.Parse("")
	if !result.Ok {
		t.Error("Expected Unknown to accept empty string")
	}

	result = schema.Parse(false)
	if !result.Ok {
		t.Error("Expected Unknown to accept false")
	}
}

// Test Unknown as passthrough validator
func TestUnknownPassthrough(t *testing.T) {
	schema := Object(Schema{
		"validated":   String().Email(),
		"passthrough": Unknown(),
	})

	data := map[string]interface{}{
		"validated":   "user@example.com",
		"passthrough": "literally anything",
	}

	result := schema.Parse(data)
	if !result.Ok {
		t.Error("Expected object with Unknown passthrough to pass")
	}

	// Invalid validated field should fail
	data["validated"] = "not-an-email"
	result = schema.Parse(data)
	if result.Ok {
		t.Error("Expected object with invalid validated field to fail")
	}
}
