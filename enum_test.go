package zogo

import (
	"testing"
)

// Test basic string enum
func TestEnumString(t *testing.T) {
	schema := Enum([]interface{}{"active", "inactive", "pending"})

	// Should pass - valid value
	result := schema.Parse("active")
	if !result.Ok {
		t.Error("Expected 'active' to pass")
	}
	if result.Value != "active" {
		t.Errorf("Expected 'active', got %v", result.Value)
	}

	result = schema.Parse("inactive")
	if !result.Ok {
		t.Error("Expected 'inactive' to pass")
	}

	result = schema.Parse("pending")
	if !result.Ok {
		t.Error("Expected 'pending' to pass")
	}

	// Should fail - invalid value
	result = schema.Parse("deleted")
	if result.Ok {
		t.Error("Expected 'deleted' to fail")
	}
}

// Test number enum
func TestEnumNumber(t *testing.T) {
	schema := Enum([]interface{}{1, 2, 3})

	// Should pass - valid values
	result := schema.Parse(1)
	if !result.Ok {
		t.Error("Expected 1 to pass")
	}

	result = schema.Parse(2)
	if !result.Ok {
		t.Error("Expected 2 to pass")
	}

	result = schema.Parse(3)
	if !result.Ok {
		t.Error("Expected 3 to pass")
	}

	// Should fail - invalid value
	result = schema.Parse(4)
	if result.Ok {
		t.Error("Expected 4 to fail")
	}
}

// Test float enum
func TestEnumFloat(t *testing.T) {
	schema := Enum([]interface{}{1.5, 2.5, 3.5})

	// Should pass
	result := schema.Parse(1.5)
	if !result.Ok {
		t.Error("Expected 1.5 to pass")
	}

	// Should fail
	result = schema.Parse(1.6)
	if result.Ok {
		t.Error("Expected 1.6 to fail")
	}
}

// Test boolean enum
func TestEnumBoolean(t *testing.T) {
	schema := Enum([]interface{}{true})

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

// Test mixed type enum (though not recommended)
func TestEnumMixed(t *testing.T) {
	schema := Enum([]interface{}{"active", 1, true})

	// Should pass - all values
	result := schema.Parse("active")
	if !result.Ok {
		t.Error("Expected 'active' to pass")
	}

	result = schema.Parse(1)
	if !result.Ok {
		t.Error("Expected 1 to pass")
	}

	result = schema.Parse(true)
	if !result.Ok {
		t.Error("Expected true to pass")
	}

	// Should fail
	result = schema.Parse("inactive")
	if result.Ok {
		t.Error("Expected 'inactive' to fail")
	}

	result = schema.Parse(2)
	if result.Ok {
		t.Error("Expected 2 to fail")
	}

	result = schema.Parse(false)
	if result.Ok {
		t.Error("Expected false to fail")
	}
}

// Test numeric type flexibility (int vs float64)
func TestEnumNumericFlexibility(t *testing.T) {
	schema := Enum([]interface{}{1, 2, 3})

	// Should pass - float64 that equals int
	result := schema.Parse(float64(1))
	if !result.Ok {
		t.Error("Expected float64(1) to pass when enum contains int(1)")
	}

	result = schema.Parse(float64(2))
	if !result.Ok {
		t.Error("Expected float64(2) to pass when enum contains int(2)")
	}
}

// Test empty enum
func TestEnumEmpty(t *testing.T) {
	schema := Enum([]interface{}{})

	// Should fail - no allowed values
	result := schema.Parse("anything")
	if result.Ok {
		t.Error("Expected any value to fail empty enum")
	}
}

// Test single value enum
func TestEnumSingle(t *testing.T) {
	schema := Enum([]interface{}{"only"})

	// Should pass
	result := schema.Parse("only")
	if !result.Ok {
		t.Error("Expected 'only' to pass")
	}

	// Should fail
	result = schema.Parse("other")
	if result.Ok {
		t.Error("Expected 'other' to fail")
	}
}

// Test nil value
func TestEnumNil(t *testing.T) {
	schema := Enum([]interface{}{"a", "b", "c"})

	result := schema.Parse(nil)
	if result.Ok {
		t.Error("Expected nil to fail")
	}
}

// Test Optional
func TestEnumOptional(t *testing.T) {
	schema := Enum([]interface{}{"a", "b", "c"}).Optional()

	// nil should pass
	result := schema.Parse(nil)
	if !result.Ok {
		t.Error("Expected nil to pass with Optional()")
	}

	// Valid enum value should still pass
	result = schema.Parse("a")
	if !result.Ok {
		t.Error("Expected 'a' to pass with Optional()")
	}
}

// Test Nullable
func TestEnumNullable(t *testing.T) {
	schema := Enum([]interface{}{"a", "b", "c"}).Nullable()

	// nil should pass
	result := schema.Parse(nil)
	if !result.Ok {
		t.Error("Expected nil to pass with Nullable()")
	}
}

// Test Default
func TestEnumDefault(t *testing.T) {
	schema := Enum([]interface{}{"a", "b", "c"}).Default("b")

	// nil should return default
	result := schema.Parse(nil)
	if !result.Ok {
		t.Error("Expected nil to pass with Default()")
	}
	if result.Value != "b" {
		t.Errorf("Expected 'b', got %v", result.Value)
	}

	// Provided value should override default
	result = schema.Parse("c")
	if !result.Ok {
		t.Error("Expected 'c' to pass")
	}
	if result.Value != "c" {
		t.Errorf("Expected 'c', got %v", result.Value)
	}
}

// Test Required
func TestEnumRequired(t *testing.T) {
	schema := Enum([]interface{}{"a", "b", "c"}).Required()

	// nil should fail
	result := schema.Parse(nil)
	if result.Ok {
		t.Error("Expected nil to fail with Required()")
	}

	// Valid value should pass
	result = schema.Parse("a")
	if !result.Ok {
		t.Error("Expected 'a' to pass with Required()")
	}
}

// Test enum in object
func TestEnumInObject(t *testing.T) {
	schema := Object(Schema{
		"name":   String(),
		"status": Enum([]interface{}{"active", "inactive", "pending"}),
	})

	data := map[string]interface{}{
		"name":   "John",
		"status": "active",
	}

	result := schema.Parse(data)
	if !result.Ok {
		t.Errorf("Expected object with enum to pass. Errors: %v", result.Errors)
	}

	// Test invalid enum value
	data["status"] = "deleted"
	result = schema.Parse(data)
	if result.Ok {
		t.Error("Expected object with invalid enum to fail")
	}
}

// Test enum in array
func TestEnumInArray(t *testing.T) {
	schema := Array(Enum([]interface{}{"red", "green", "blue"}))

	data := []interface{}{"red", "green", "blue", "red"}

	result := schema.Parse(data)
	if !result.Ok {
		t.Errorf("Expected array of valid enums to pass. Errors: %v", result.Errors)
	}

	// Test with invalid value
	data = []interface{}{"red", "yellow", "blue"}
	result = schema.Parse(data)
	if result.Ok {
		t.Error("Expected array with invalid enum to fail")
	}

	// Check error path
	if len(result.Errors) == 0 {
		t.Error("Expected errors")
	} else if result.Errors[0].Path != "[1]" {
		t.Errorf("Expected error path '[1]', got '%s'", result.Errors[0].Path)
	}
}

// Test error message contains expected values
func TestEnumErrorMessage(t *testing.T) {
	schema := Enum([]interface{}{"a", "b", "c"})

	result := schema.Parse("d")
	if result.Ok {
		t.Error("Expected invalid value to fail")
	}

	if len(result.Errors) == 0 {
		t.Error("Expected error message")
	}

	// Error message should mention allowed values
	errorMsg := result.Errors[0].Message
	if errorMsg == "" {
		t.Error("Expected non-empty error message")
	}
}

// Test case sensitivity for strings
func TestEnumCaseSensitive(t *testing.T) {
	schema := Enum([]interface{}{"Active", "Inactive"})

	// Should pass - exact match
	result := schema.Parse("Active")
	if !result.Ok {
		t.Error("Expected 'Active' to pass")
	}

	// Should fail - different case
	result = schema.Parse("active")
	if result.Ok {
		t.Error("Expected 'active' to fail (case sensitive)")
	}
}

// Test enum with HTTP status codes (common use case)
func TestEnumHTTPStatus(t *testing.T) {
	schema := Enum([]interface{}{200, 201, 400, 401, 403, 404, 500})

	result := schema.Parse(200)
	if !result.Ok {
		t.Error("Expected 200 to pass")
	}

	result = schema.Parse(404)
	if !result.Ok {
		t.Error("Expected 404 to pass")
	}

	result = schema.Parse(418) // I'm a teapot
	if result.Ok {
		t.Error("Expected 418 to fail")
	}
}
