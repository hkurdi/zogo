package zogo

import (
	"testing"
)

// Test basic union - string or number
func TestUnionStringOrNumber(t *testing.T) {
	schema := Union(String(), Number())

	// Should pass - string
	result := schema.Parse("hello")
	if !result.Ok {
		t.Errorf("Expected string to pass Union(String, Number). Errors: %v", result.Errors)
	}
	if result.Value != "hello" {
		t.Errorf("Expected 'hello', got %v", result.Value)
	}

	// Should pass - number
	result = schema.Parse(42)
	if !result.Ok {
		t.Errorf("Expected number to pass Union(String, Number). Errors: %v", result.Errors)
	}
	if result.Value != float64(42) {
		t.Errorf("Expected 42, got %v", result.Value)
	}

	// Should fail - boolean
	result = schema.Parse(true)
	if result.Ok {
		t.Error("Expected boolean to fail Union(String, Number)")
	}
}

// Test union with three types
func TestUnionThreeTypes(t *testing.T) {
	schema := Union(String(), Number(), Boolean())

	// All three should pass
	result := schema.Parse("hello")
	if !result.Ok {
		t.Error("Expected string to pass")
	}

	result = schema.Parse(42)
	if !result.Ok {
		t.Error("Expected number to pass")
	}

	result = schema.Parse(true)
	if !result.Ok {
		t.Error("Expected boolean to pass")
	}

	// Object should fail
	result = schema.Parse(map[string]interface{}{"key": "value"})
	if result.Ok {
		t.Error("Expected object to fail")
	}
}

// Test union with validation constraints
func TestUnionWithConstraints(t *testing.T) {
	schema := Union(
		String().Email(),
		Number().Min(100),
	)

	// Valid email should pass
	result := schema.Parse("user@example.com")
	if !result.Ok {
		t.Error("Expected valid email to pass")
	}

	// Large number should pass
	result = schema.Parse(150)
	if !result.Ok {
		t.Error("Expected number >= 100 to pass")
	}

	// Invalid email should fail
	result = schema.Parse("notanemail")
	if result.Ok {
		t.Error("Expected invalid email to fail")
	}

	// Small number should fail
	result = schema.Parse(50)
	if result.Ok {
		t.Error("Expected number < 100 to fail")
	}
}

// Test union with enum
func TestUnionWithEnum(t *testing.T) {
	schema := Union(
		Enum([]interface{}{"active", "inactive"}),
		Number(),
	)

	// Enum value should pass
	result := schema.Parse("active")
	if !result.Ok {
		t.Error("Expected enum value to pass")
	}

	// Number should pass
	result = schema.Parse(42)
	if !result.Ok {
		t.Error("Expected number to pass")
	}

	// Invalid enum and not a number should fail
	result = schema.Parse("deleted")
	if result.Ok {
		t.Error("Expected invalid value to fail")
	}
}

// Test union with literal
func TestUnionWithLiteral(t *testing.T) {
	schema := Union(
		Literal("success"),
		Literal("error"),
		Literal("pending"),
	)

	// All literals should pass
	result := schema.Parse("success")
	if !result.Ok {
		t.Error("Expected 'success' to pass")
	}

	result = schema.Parse("error")
	if !result.Ok {
		t.Error("Expected 'error' to pass")
	}

	result = schema.Parse("pending")
	if !result.Ok {
		t.Error("Expected 'pending' to pass")
	}

	// Other values should fail
	result = schema.Parse("failed")
	if result.Ok {
		t.Error("Expected 'failed' to fail")
	}
}

// Test union with objects
func TestUnionObjects(t *testing.T) {
	userSchema := Object(Schema{
		"type": Literal("user"),
		"name": String(),
	})

	adminSchema := Object(Schema{
		"type":        Literal("admin"),
		"permissions": Array(String()),
	})

	schema := Union(userSchema, adminSchema)

	// User object should pass
	userData := map[string]interface{}{
		"type": "user",
		"name": "John",
	}
	result := schema.Parse(userData)
	if !result.Ok {
		t.Errorf("Expected user object to pass. Errors: %v", result.Errors)
	}

	// Admin object should pass
	adminData := map[string]interface{}{
		"type":        "admin",
		"permissions": []interface{}{"read", "write"},
	}
	result = schema.Parse(adminData)
	if !result.Ok {
		t.Errorf("Expected admin object to pass. Errors: %v", result.Errors)
	}

	// Invalid object should fail
	invalidData := map[string]interface{}{
		"type": "guest",
	}
	result = schema.Parse(invalidData)
	if result.Ok {
		t.Error("Expected invalid object to fail")
	}
}

// Test union with arrays
func TestUnionArrays(t *testing.T) {
	schema := Union(
		Array(String()),
		Array(Number()),
	)

	// Array of strings should pass
	result := schema.Parse([]interface{}{"a", "b", "c"})
	if !result.Ok {
		t.Error("Expected string array to pass")
	}

	// Array of numbers should pass
	result = schema.Parse([]interface{}{1, 2, 3})
	if !result.Ok {
		t.Error("Expected number array to pass")
	}

	// Mixed array should fail both
	result = schema.Parse([]interface{}{"a", 1, "b"})
	if result.Ok {
		t.Error("Expected mixed array to fail")
	}
}

// Test nested unions
func TestUnionNested(t *testing.T) {
	schema := Union(
		String(),
		Union(Number(), Boolean()),
	)

	// String should pass
	result := schema.Parse("hello")
	if !result.Ok {
		t.Error("Expected string to pass")
	}

	// Number should pass (via nested union)
	result = schema.Parse(42)
	if !result.Ok {
		t.Error("Expected number to pass")
	}

	// Boolean should pass (via nested union)
	result = schema.Parse(true)
	if !result.Ok {
		t.Error("Expected boolean to pass")
	}
}

// Test nil value
func TestUnionNil(t *testing.T) {
	schema := Union(String(), Number())

	result := schema.Parse(nil)
	if result.Ok {
		t.Error("Expected nil to fail")
	}
}

// Test Optional
func TestUnionOptional(t *testing.T) {
	schema := Union(String(), Number()).Optional()

	// nil should pass
	result := schema.Parse(nil)
	if !result.Ok {
		t.Error("Expected nil to pass with Optional()")
	}

	// String should pass
	result = schema.Parse("hello")
	if !result.Ok {
		t.Error("Expected string to pass")
	}

	// Number should pass
	result = schema.Parse(42)
	if !result.Ok {
		t.Error("Expected number to pass")
	}
}

// Test Nullable
func TestUnionNullable(t *testing.T) {
	schema := Union(String(), Number()).Nullable()

	// nil should pass
	result := schema.Parse(nil)
	if !result.Ok {
		t.Error("Expected nil to pass with Nullable()")
	}
}

// Test Required
func TestUnionRequired(t *testing.T) {
	schema := Union(String(), Number()).Required()

	// nil should fail
	result := schema.Parse(nil)
	if result.Ok {
		t.Error("Expected nil to fail with Required()")
	}

	// Valid values should pass
	result = schema.Parse("hello")
	if !result.Ok {
		t.Error("Expected string to pass")
	}
}

// Test union in object
func TestUnionInObject(t *testing.T) {
	schema := Object(Schema{
		"name":  String(),
		"value": Union(String(), Number()),
	})

	// String value
	data := map[string]interface{}{
		"name":  "test",
		"value": "hello",
	}
	result := schema.Parse(data)
	if !result.Ok {
		t.Errorf("Expected object with string value to pass. Errors: %v", result.Errors)
	}

	// Number value
	data["value"] = 42
	result = schema.Parse(data)
	if !result.Ok {
		t.Errorf("Expected object with number value to pass. Errors: %v", result.Errors)
	}

	// Boolean value (invalid)
	data["value"] = true
	result = schema.Parse(data)
	if result.Ok {
		t.Error("Expected object with boolean value to fail")
	}
}

// Test union in array
func TestUnionInArray(t *testing.T) {
	schema := Array(Union(String(), Number()))

	// Mixed string and number array should pass
	data := []interface{}{"hello", 42, "world", 100}
	result := schema.Parse(data)
	if !result.Ok {
		t.Errorf("Expected mixed string/number array to pass. Errors: %v", result.Errors)
	}

	// Array with invalid element should fail
	data = []interface{}{"hello", 42, true}
	result = schema.Parse(data)
	if result.Ok {
		t.Error("Expected array with boolean to fail")
	}
}

// Test error message contains all failures
func TestUnionErrorMessage(t *testing.T) {
	schema := Union(String().Email(), Number().Min(100))

	result := schema.Parse("invalid")
	if result.Ok {
		t.Error("Expected invalid value to fail")
	}

	if len(result.Errors) == 0 {
		t.Error("Expected error message")
	}

	// Error message should mention both options failed
	errorMsg := result.Errors[0].Message
	if errorMsg == "" {
		t.Error("Expected non-empty error message")
	}
}

// Test discriminated union (common pattern)
func TestUnionDiscriminated(t *testing.T) {
	successSchema := Object(Schema{
		"status": Literal("success"),
		"data":   String(),
	})

	errorSchema := Object(Schema{
		"status":  Literal("error"),
		"message": String(),
	})

	schema := Union(successSchema, errorSchema)

	// Success response
	successData := map[string]interface{}{
		"status": "success",
		"data":   "result",
	}
	result := schema.Parse(successData)
	if !result.Ok {
		t.Errorf("Expected success response to pass. Errors: %v", result.Errors)
	}

	// Error response
	errorData := map[string]interface{}{
		"status":  "error",
		"message": "something went wrong",
	}
	result = schema.Parse(errorData)
	if !result.Ok {
		t.Errorf("Expected error response to pass. Errors: %v", result.Errors)
	}

	// Invalid status
	invalidData := map[string]interface{}{
		"status": "pending",
	}
	result = schema.Parse(invalidData)
	if result.Ok {
		t.Error("Expected invalid status to fail")
	}
}

// Test union with single validator (edge case)
func TestUnionSingle(t *testing.T) {
	schema := Union(String())

	result := schema.Parse("hello")
	if !result.Ok {
		t.Error("Expected string to pass single-validator union")
	}

	result = schema.Parse(42)
	if result.Ok {
		t.Error("Expected number to fail single-validator union")
	}
}

// Test empty union (edge case)
func TestUnionEmpty(t *testing.T) {
	schema := Union()

	// Should fail - no validators to match
	result := schema.Parse("anything")
	if result.Ok {
		t.Error("Expected empty union to fail")
	}
}

// Test union returns first successful result
func TestUnionFirstSuccess(t *testing.T) {
	schema := Union(
		String().Trim(),
		String().ToUpperCase(),
	)

	// Should use the first validator that succeeds
	result := schema.Parse("  hello  ")
	if !result.Ok {
		t.Error("Expected string to pass")
	}

	// Result should be from first successful validator (Trim)
	if result.Value != "hello" {
		t.Errorf("Expected 'hello' (trimmed), got '%v'", result.Value)
	}
}
