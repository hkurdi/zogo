package zogo

import (
	"testing"
)

// Test basic string literal
func TestLiteralString(t *testing.T) {
	schema := Literal("user")

	// Should pass - exact match
	result := schema.Parse("user")
	if !result.Ok {
		t.Error("Expected 'user' to pass")
	}
	if result.Value != "user" {
		t.Errorf("Expected 'user', got %v", result.Value)
	}

	// Should fail - different value
	result = schema.Parse("admin")
	if result.Ok {
		t.Error("Expected 'admin' to fail")
	}

	result = schema.Parse("User") // case sensitive
	if result.Ok {
		t.Error("Expected 'User' to fail (case sensitive)")
	}
}

// Test number literal
func TestLiteralNumber(t *testing.T) {
	schema := Literal(42)

	// Should pass
	result := schema.Parse(42)
	if !result.Ok {
		t.Error("Expected 42 to pass")
	}

	// Should fail
	result = schema.Parse(43)
	if result.Ok {
		t.Error("Expected 43 to fail")
	}

	result = schema.Parse(0)
	if result.Ok {
		t.Error("Expected 0 to fail")
	}
}

// Test float literal
func TestLiteralFloat(t *testing.T) {
	schema := Literal(3.14)

	// Should pass
	result := schema.Parse(3.14)
	if !result.Ok {
		t.Error("Expected 3.14 to pass")
	}

	// Should fail
	result = schema.Parse(3.15)
	if result.Ok {
		t.Error("Expected 3.15 to fail")
	}
}

// Test boolean literal true
func TestLiteralBooleanTrue(t *testing.T) {
	schema := Literal(true)

	// Should pass
	result := schema.Parse(true)
	if !result.Ok {
		t.Error("Expected true to pass")
	}

	// Should fail
	result = schema.Parse(false)
	if result.Ok {
		t.Error("Expected false to fail")
	}
}

// Test boolean literal false
func TestLiteralBooleanFalse(t *testing.T) {
	schema := Literal(false)

	// Should pass
	result := schema.Parse(false)
	if !result.Ok {
		t.Error("Expected false to pass")
	}

	// Should fail
	result = schema.Parse(true)
	if result.Ok {
		t.Error("Expected true to fail")
	}
}

// Test numeric type flexibility
func TestLiteralNumericFlexibility(t *testing.T) {
	schema := Literal(1)

	// Should pass - float64 that equals int
	result := schema.Parse(float64(1))
	if !result.Ok {
		t.Error("Expected float64(1) to pass when literal is int(1)")
	}
}

// Test zero values
func TestLiteralZeroValues(t *testing.T) {
	// Zero int
	schema := Literal(0)
	result := schema.Parse(0)
	if !result.Ok {
		t.Error("Expected 0 to pass")
	}

	// Empty string
	schema = Literal("")
	result = schema.Parse("")
	if !result.Ok {
		t.Error("Expected empty string to pass")
	}

	// False
	schema = Literal(false)
	result = schema.Parse(false)
	if !result.Ok {
		t.Error("Expected false to pass")
	}
}

// Test nil value
func TestLiteralNil(t *testing.T) {
	schema := Literal("value")

	result := schema.Parse(nil)
	if result.Ok {
		t.Error("Expected nil to fail")
	}
}

// Test Optional
func TestLiteralOptional(t *testing.T) {
	schema := Literal("value").Optional()

	// nil should pass
	result := schema.Parse(nil)
	if !result.Ok {
		t.Error("Expected nil to pass with Optional()")
	}

	// Matching value should pass
	result = schema.Parse("value")
	if !result.Ok {
		t.Error("Expected 'value' to pass with Optional()")
	}

	// Non-matching value should fail
	result = schema.Parse("other")
	if result.Ok {
		t.Error("Expected 'other' to fail even with Optional()")
	}
}

// Test Nullable
func TestLiteralNullable(t *testing.T) {
	schema := Literal("value").Nullable()

	// nil should pass
	result := schema.Parse(nil)
	if !result.Ok {
		t.Error("Expected nil to pass with Nullable()")
	}

	// Matching value should pass
	result = schema.Parse("value")
	if !result.Ok {
		t.Error("Expected 'value' to pass with Nullable()")
	}
}

// Test Required
func TestLiteralRequired(t *testing.T) {
	schema := Literal("value").Required()

	// nil should fail
	result := schema.Parse(nil)
	if result.Ok {
		t.Error("Expected nil to fail with Required()")
	}

	// Matching value should pass
	result = schema.Parse("value")
	if !result.Ok {
		t.Error("Expected 'value' to pass with Required()")
	}
}

// Test literal in object (discriminated union pattern)
func TestLiteralInObject(t *testing.T) {
	schema := Object(Schema{
		"type": Literal("user"),
		"name": String(),
	})

	data := map[string]interface{}{
		"type": "user",
		"name": "John",
	}

	result := schema.Parse(data)
	if !result.Ok {
		t.Errorf("Expected object with literal to pass. Errors: %v", result.Errors)
	}

	// Test with wrong type
	data["type"] = "admin"
	result = schema.Parse(data)
	if result.Ok {
		t.Error("Expected object with wrong literal to fail")
	}
}

// Test literal in array
func TestLiteralInArray(t *testing.T) {
	schema := Array(Literal("active"))

	// All matching
	data := []interface{}{"active", "active", "active"}
	result := schema.Parse(data)
	if !result.Ok {
		t.Errorf("Expected array of matching literals to pass. Errors: %v", result.Errors)
	}

	// One non-matching
	data = []interface{}{"active", "inactive", "active"}
	result = schema.Parse(data)
	if result.Ok {
		t.Error("Expected array with non-matching literal to fail")
	}

	// Check error path
	if len(result.Errors) == 0 {
		t.Error("Expected errors")
	} else if result.Errors[0].Path != "[1]" {
		t.Errorf("Expected error path '[1]', got '%s'", result.Errors[0].Path)
	}
}

// Test discriminated union pattern (common use case)
func TestLiteralDiscriminatedUnion(t *testing.T) {
	userSchema := Object(Schema{
		"type":  Literal("user"),
		"name":  String(),
		"email": String().Email(),
	})

	adminSchema := Object(Schema{
		"type":        Literal("admin"),
		"name":        String(),
		"permissions": Array(String()),
	})

	// Validate user
	userData := map[string]interface{}{
		"type":  "user",
		"name":  "John",
		"email": "john@example.com",
	}

	result := userSchema.Parse(userData)
	if !result.Ok {
		t.Errorf("Expected user data to pass user schema. Errors: %v", result.Errors)
	}

	// User data should fail admin schema
	result = adminSchema.Parse(userData)
	if result.Ok {
		t.Error("Expected user data to fail admin schema")
	}

	// Validate admin
	adminData := map[string]interface{}{
		"type":        "admin",
		"name":        "Jane",
		"permissions": []interface{}{"read", "write", "delete"},
	}

	result = adminSchema.Parse(adminData)
	if !result.Ok {
		t.Errorf("Expected admin data to pass admin schema. Errors: %v", result.Errors)
	}

	// Admin data should fail user schema
	result = userSchema.Parse(adminData)
	if result.Ok {
		t.Error("Expected admin data to fail user schema")
	}
}

// Test error message
func TestLiteralErrorMessage(t *testing.T) {
	schema := Literal("expected")

	result := schema.Parse("actual")
	if result.Ok {
		t.Error("Expected non-matching value to fail")
	}

	if len(result.Errors) == 0 {
		t.Error("Expected error message")
	}

	// Error message should mention both expected and received
	errorMsg := result.Errors[0].Message
	if errorMsg == "" {
		t.Error("Expected non-empty error message")
	}
}

// Test literal with version numbers (common use case)
func TestLiteralVersion(t *testing.T) {
	schema := Object(Schema{
		"version": Literal(1),
		"data":    String(),
	})

	// v1 data
	data := map[string]interface{}{
		"version": 1,
		"data":    "some data",
	}

	result := schema.Parse(data)
	if !result.Ok {
		t.Error("Expected v1 data to pass")
	}

	// Wrong version
	data["version"] = 2
	result = schema.Parse(data)
	if result.Ok {
		t.Error("Expected v2 data to fail v1 schema")
	}
}

// Test literal with status codes
func TestLiteralStatusCode(t *testing.T) {
	successSchema := Object(Schema{
		"status": Literal(200),
		"data":   String(),
	})

	data := map[string]interface{}{
		"status": 200,
		"data":   "success",
	}

	result := successSchema.Parse(data)
	if !result.Ok {
		t.Error("Expected 200 status to pass")
	}

	data["status"] = 404
	result = successSchema.Parse(data)
	if result.Ok {
		t.Error("Expected 404 status to fail 200 schema")
	}
}

// Test type mismatch
func TestLiteralTypeMismatch(t *testing.T) {
	schema := Literal("string")

	// Different types should fail
	result := schema.Parse(123)
	if result.Ok {
		t.Error("Expected number to fail string literal")
	}

	result = schema.Parse(true)
	if result.Ok {
		t.Error("Expected boolean to fail string literal")
	}
}
