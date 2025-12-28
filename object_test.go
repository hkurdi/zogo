package zogo

import (
	"testing"
)

// Test basic object validation
func TestObjectBasic(t *testing.T) {
	schema := Object(Schema{
		"name": String(),
		"age":  Number(),
	})

	data := map[string]interface{}{
		"name": "John",
		"age":  30,
	}

	result := schema.Parse(data)
	if !result.Ok {
		t.Errorf("Expected valid object to pass. Errors: %v", result.Errors)
	}

	// Check result values
	resultMap, ok := result.Value.(map[string]interface{})
	if !ok {
		t.Error("Expected result to be a map")
	}

	if resultMap["name"] != "John" {
		t.Errorf("Expected name 'John', got %v", resultMap["name"])
	}

	if resultMap["age"] != float64(30) {
		t.Errorf("Expected age 30, got %v", resultMap["age"])
	}
}

// Test nested objects
func TestObjectNested(t *testing.T) {
	schema := Object(Schema{
		"user": Object(Schema{
			"name":  String(),
			"email": String().Email(),
		}),
		"active": Boolean(),
	})

	data := map[string]interface{}{
		"user": map[string]interface{}{
			"name":  "John",
			"email": "john@example.com",
		},
		"active": true,
	}

	result := schema.Parse(data)
	if !result.Ok {
		t.Errorf("Expected nested object to pass. Errors: %v", result.Errors)
	}
}

// Test nested object with error - check path
func TestObjectNestedError(t *testing.T) {
	schema := Object(Schema{
		"user": Object(Schema{
			"email": String().Email(),
		}),
	})

	data := map[string]interface{}{
		"user": map[string]interface{}{
			"email": "notanemail",
		},
	}

	result := schema.Parse(data)
	if result.Ok {
		t.Error("Expected invalid nested email to fail")
	}

	// Check error path
	if len(result.Errors) == 0 {
		t.Error("Expected errors")
	} else if result.Errors[0].Path != "user.email" {
		t.Errorf("Expected error path 'user.email', got '%s'", result.Errors[0].Path)
	}
}

// Test required fields
func TestObjectRequiredFields(t *testing.T) {
	schema := Object(Schema{
		"name":  String().Required(),
		"email": String().Email().Required(),
	})

	// Missing required field
	data := map[string]interface{}{
		"name": "John",
		// email missing
	}

	result := schema.Parse(data)
	if result.Ok {
		t.Error("Expected missing required field to fail")
	}

	if len(result.Errors) == 0 {
		t.Error("Expected error for missing field")
	}
}

// Test optional fields
func TestObjectOptionalFields(t *testing.T) {
	schema := Object(Schema{
		"name": String().Required(),
		"age":  Number().Optional(),
	})

	// Optional field missing - should pass
	data := map[string]interface{}{
		"name": "John",
		// age missing but optional
	}

	result := schema.Parse(data)
	if !result.Ok {
		t.Errorf("Expected object with missing optional field to pass. Errors: %v", result.Errors)
	}

	resultMap := result.Value.(map[string]interface{})

	// Optional field shouldn't be in result if not provided
	if _, exists := resultMap["age"]; exists {
		t.Error("Expected missing optional field to not appear in result")
	}
}

// Test Strip mode (default)
func TestObjectStrip(t *testing.T) {
	schema := Object(Schema{
		"name": String(),
	}).Strip()

	data := map[string]interface{}{
		"name":    "John",
		"unknown": "field",
	}

	result := schema.Parse(data)
	if !result.Ok {
		t.Error("Expected object with unknown field to pass in Strip mode")
	}

	resultMap := result.Value.(map[string]interface{})

	// Unknown field should be stripped
	if _, exists := resultMap["unknown"]; exists {
		t.Error("Expected unknown field to be stripped")
	}

	if resultMap["name"] != "John" {
		t.Error("Expected known field to be preserved")
	}
}

// Test Passthrough mode
func TestObjectPassthrough(t *testing.T) {
	schema := Object(Schema{
		"name": String(),
	}).Passthrough()

	data := map[string]interface{}{
		"name":    "John",
		"unknown": "field",
	}

	result := schema.Parse(data)
	if !result.Ok {
		t.Error("Expected object with unknown field to pass in Passthrough mode")
	}

	resultMap := result.Value.(map[string]interface{})

	// Unknown field should be passed through
	if resultMap["unknown"] != "field" {
		t.Error("Expected unknown field to be passed through")
	}

	if resultMap["name"] != "John" {
		t.Error("Expected known field to be preserved")
	}
}

// Test Strict mode
func TestObjectStrict(t *testing.T) {
	schema := Object(Schema{
		"name": String(),
	}).Strict()

	data := map[string]interface{}{
		"name":    "John",
		"unknown": "field",
	}

	result := schema.Parse(data)
	if result.Ok {
		t.Error("Expected object with unknown field to fail in Strict mode")
	}

	if len(result.Errors) == 0 {
		t.Error("Expected error for unknown field")
	}
}

// Test nil value
func TestObjectNil(t *testing.T) {
	schema := Object(Schema{
		"name": String(),
	})

	result := schema.Parse(nil)
	if result.Ok {
		t.Error("Expected nil to fail")
	}
}

// Test Optional object
func TestObjectOptional(t *testing.T) {
	schema := Object(Schema{
		"name": String(),
	}).Optional()

	// nil should pass
	result := schema.Parse(nil)
	if !result.Ok {
		t.Error("Expected nil to pass with Optional()")
	}

	// Valid object should still pass
	data := map[string]interface{}{
		"name": "John",
	}
	result = schema.Parse(data)
	if !result.Ok {
		t.Error("Expected valid object to pass with Optional()")
	}
}

// Test Nullable object
func TestObjectNullable(t *testing.T) {
	schema := Object(Schema{
		"name": String(),
	}).Nullable()

	// nil should pass
	result := schema.Parse(nil)
	if !result.Ok {
		t.Error("Expected nil to pass with Nullable()")
	}
}

// Test invalid type
func TestObjectInvalidType(t *testing.T) {
	schema := Object(Schema{
		"name": String(),
	})

	result := schema.Parse("not an object")
	if result.Ok {
		t.Error("Expected string to fail object validation")
	}

	result = schema.Parse(123)
	if result.Ok {
		t.Error("Expected number to fail object validation")
	}
}

// Test multiple validation errors
func TestObjectMultipleErrors(t *testing.T) {
	schema := Object(Schema{
		"name":  String().Min(5),
		"email": String().Email(),
		"age":   Number().Min(18),
	})

	data := map[string]interface{}{
		"name":  "Jo",         // too short
		"email": "notanemail", // invalid email
		"age":   10,           // too young
	}

	result := schema.Parse(data)
	if result.Ok {
		t.Error("Expected object with multiple invalid fields to fail")
	}

	// Should have 3 errors
	if len(result.Errors) != 3 {
		t.Errorf("Expected 3 errors, got %d", len(result.Errors))
	}
}

// Test empty object
func TestObjectEmpty(t *testing.T) {
	schema := Object(Schema{})

	data := map[string]interface{}{}

	result := schema.Parse(data)
	if !result.Ok {
		t.Error("Expected empty object to pass empty schema")
	}
}

// Test deeply nested objects
func TestObjectDeeplyNested(t *testing.T) {
	schema := Object(Schema{
		"user": Object(Schema{
			"profile": Object(Schema{
				"name": String(),
			}),
		}),
	})

	data := map[string]interface{}{
		"user": map[string]interface{}{
			"profile": map[string]interface{}{
				"name": "John",
			},
		},
	}

	result := schema.Parse(data)
	if !result.Ok {
		t.Errorf("Expected deeply nested object to pass. Errors: %v", result.Errors)
	}
}

// Test deeply nested error path
func TestObjectDeeplyNestedErrorPath(t *testing.T) {
	schema := Object(Schema{
		"user": Object(Schema{
			"profile": Object(Schema{
				"email": String().Email(),
			}),
		}),
	})

	data := map[string]interface{}{
		"user": map[string]interface{}{
			"profile": map[string]interface{}{
				"email": "notanemail",
			},
		},
	}

	result := schema.Parse(data)
	if result.Ok {
		t.Error("Expected invalid deeply nested email to fail")
	}

	// Check error path
	if len(result.Errors) == 0 {
		t.Error("Expected errors")
	} else if result.Errors[0].Path != "user.profile.email" {
		t.Errorf("Expected error path 'user.profile.email', got '%s'", result.Errors[0].Path)
	}
}
