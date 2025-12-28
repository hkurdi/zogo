package zogo

import (
	"testing"
)

// Test basic array validation
func TestArrayBasic(t *testing.T) {
	schema := Array(String())

	data := []interface{}{"hello", "world"}

	result := schema.Parse(data)
	if !result.Ok {
		t.Errorf("Expected valid array to pass. Errors: %v", result.Errors)
	}

	// Check result values
	resultArr, ok := result.Value.([]interface{})
	if !ok {
		t.Error("Expected result to be a slice")
	}

	if len(resultArr) != 2 {
		t.Errorf("Expected 2 elements, got %d", len(resultArr))
	}

	if resultArr[0] != "hello" || resultArr[1] != "world" {
		t.Error("Expected array values to be preserved")
	}
}

// Test array of numbers
func TestArrayNumbers(t *testing.T) {
	schema := Array(Number())

	data := []interface{}{1, 2, 3, 4, 5}

	result := schema.Parse(data)
	if !result.Ok {
		t.Errorf("Expected valid number array to pass. Errors: %v", result.Errors)
	}

	resultArr := result.Value.([]interface{})
	if len(resultArr) != 5 {
		t.Errorf("Expected 5 elements, got %d", len(resultArr))
	}
}

// Test array of booleans
func TestArrayBooleans(t *testing.T) {
	schema := Array(Boolean())

	data := []interface{}{true, false, true}

	result := schema.Parse(data)
	if !result.Ok {
		t.Errorf("Expected valid boolean array to pass. Errors: %v", result.Errors)
	}
}

// Test array of objects
func TestArrayObjects(t *testing.T) {
	schema := Array(Object(Schema{
		"name": String(),
		"age":  Number(),
	}))

	data := []interface{}{
		map[string]interface{}{
			"name": "John",
			"age":  30,
		},
		map[string]interface{}{
			"name": "Jane",
			"age":  25,
		},
	}

	result := schema.Parse(data)
	if !result.Ok {
		t.Errorf("Expected valid object array to pass. Errors: %v", result.Errors)
	}
}

// Test nested arrays
func TestArrayNested(t *testing.T) {
	schema := Array(Array(Number()))

	data := []interface{}{
		[]interface{}{1, 2, 3},
		[]interface{}{4, 5, 6},
	}

	result := schema.Parse(data)
	if !result.Ok {
		t.Errorf("Expected nested array to pass. Errors: %v", result.Errors)
	}
}

// Test array element validation error
func TestArrayElementError(t *testing.T) {
	schema := Array(String().Email())

	data := []interface{}{
		"valid@example.com",
		"invalid",
		"another@example.com",
	}

	result := schema.Parse(data)
	if result.Ok {
		t.Error("Expected array with invalid element to fail")
	}

	// Should have error at index 1
	if len(result.Errors) == 0 {
		t.Error("Expected errors")
	} else if result.Errors[0].Path != "[1]" {
		t.Errorf("Expected error path '[1]', got '%s'", result.Errors[0].Path)
	}
}

// Test array element error with nested object
func TestArrayNestedObjectError(t *testing.T) {
	schema := Array(Object(Schema{
		"email": String().Email(),
	}))

	data := []interface{}{
		map[string]interface{}{
			"email": "valid@example.com",
		},
		map[string]interface{}{
			"email": "invalid",
		},
	}

	result := schema.Parse(data)
	if result.Ok {
		t.Error("Expected array with invalid nested object to fail")
	}

	// Check error path
	if len(result.Errors) == 0 {
		t.Error("Expected errors")
	} else if result.Errors[0].Path != "[1].email" {
		t.Errorf("Expected error path '[1].email', got '%s'", result.Errors[0].Path)
	}
}

// Test Min length
func TestArrayMin(t *testing.T) {
	schema := Array(String()).Min(2)

	// Should pass
	result := schema.Parse([]interface{}{"a", "b"})
	if !result.Ok {
		t.Error("Expected array with 2 elements to pass Min(2)")
	}

	result = schema.Parse([]interface{}{"a", "b", "c"})
	if !result.Ok {
		t.Error("Expected array with 3 elements to pass Min(2)")
	}

	// Should fail
	result = schema.Parse([]interface{}{"a"})
	if result.Ok {
		t.Error("Expected array with 1 element to fail Min(2)")
	}
}

// Test Max length
func TestArrayMax(t *testing.T) {
	schema := Array(String()).Max(3)

	// Should pass
	result := schema.Parse([]interface{}{"a", "b", "c"})
	if !result.Ok {
		t.Error("Expected array with 3 elements to pass Max(3)")
	}

	result = schema.Parse([]interface{}{"a", "b"})
	if !result.Ok {
		t.Error("Expected array with 2 elements to pass Max(3)")
	}

	// Should fail
	result = schema.Parse([]interface{}{"a", "b", "c", "d"})
	if result.Ok {
		t.Error("Expected array with 4 elements to fail Max(3)")
	}
}

// Test exact Length
func TestArrayLength(t *testing.T) {
	schema := Array(String()).Length(3)

	// Should pass
	result := schema.Parse([]interface{}{"a", "b", "c"})
	if !result.Ok {
		t.Error("Expected array with 3 elements to pass Length(3)")
	}

	// Should fail - too few
	result = schema.Parse([]interface{}{"a", "b"})
	if result.Ok {
		t.Error("Expected array with 2 elements to fail Length(3)")
	}

	// Should fail - too many
	result = schema.Parse([]interface{}{"a", "b", "c", "d"})
	if result.Ok {
		t.Error("Expected array with 4 elements to fail Length(3)")
	}
}

// Test NonEmpty
func TestArrayNonEmpty(t *testing.T) {
	schema := Array(String()).NonEmpty()

	// Should pass
	result := schema.Parse([]interface{}{"a"})
	if !result.Ok {
		t.Error("Expected non-empty array to pass NonEmpty()")
	}

	// Should fail
	result = schema.Parse([]interface{}{})
	if result.Ok {
		t.Error("Expected empty array to fail NonEmpty()")
	}
}

// Test empty array (without NonEmpty)
func TestArrayEmpty(t *testing.T) {
	schema := Array(String())

	// Empty array should pass without NonEmpty
	result := schema.Parse([]interface{}{})
	if !result.Ok {
		t.Error("Expected empty array to pass without NonEmpty()")
	}
}

// Test nil value
func TestArrayNil(t *testing.T) {
	schema := Array(String())

	result := schema.Parse(nil)
	if result.Ok {
		t.Error("Expected nil to fail")
	}
}

// Test Optional
func TestArrayOptional(t *testing.T) {
	schema := Array(String()).Optional()

	// nil should pass
	result := schema.Parse(nil)
	if !result.Ok {
		t.Error("Expected nil to pass with Optional()")
	}

	// Valid array should still pass
	result = schema.Parse([]interface{}{"hello"})
	if !result.Ok {
		t.Error("Expected valid array to pass with Optional()")
	}
}

// Test Nullable
func TestArrayNullable(t *testing.T) {
	schema := Array(String()).Nullable()

	// nil should pass
	result := schema.Parse(nil)
	if !result.Ok {
		t.Error("Expected nil to pass with Nullable()")
	}
}

// Test invalid type
func TestArrayInvalidType(t *testing.T) {
	schema := Array(String())

	result := schema.Parse("not an array")
	if result.Ok {
		t.Error("Expected string to fail array validation")
	}

	result = schema.Parse(123)
	if result.Ok {
		t.Error("Expected number to fail array validation")
	}

	result = schema.Parse(map[string]interface{}{"key": "value"})
	if result.Ok {
		t.Error("Expected object to fail array validation")
	}
}

// Test multiple element errors
func TestArrayMultipleErrors(t *testing.T) {
	schema := Array(Number().Min(10))

	data := []interface{}{5, 15, 3, 20, 7}

	result := schema.Parse(data)
	if result.Ok {
		t.Error("Expected array with multiple invalid elements to fail")
	}

	// Should have errors at indices 0, 2, 4
	if len(result.Errors) != 3 {
		t.Errorf("Expected 3 errors, got %d", len(result.Errors))
	}
}

// Test chaining validators
func TestArrayChained(t *testing.T) {
	schema := Array(String().Min(3)).Min(2).Max(5).NonEmpty()

	// Should pass
	result := schema.Parse([]interface{}{"hello", "world"})
	if !result.Ok {
		t.Error("Expected valid array to pass all validators")
	}

	// Should fail - array too short
	result = schema.Parse([]interface{}{"hello"})
	if result.Ok {
		t.Error("Expected array with 1 element to fail Min(2)")
	}

	// Should fail - array too long
	result = schema.Parse([]interface{}{"a", "b", "c", "d", "e", "f"})
	if result.Ok {
		t.Error("Expected array with 6 elements to fail Max(5)")
	}

	// Should fail - element too short
	result = schema.Parse([]interface{}{"hi", "world"})
	if result.Ok {
		t.Error("Expected array with short string to fail")
	}
}

// Test array with validation on elements
func TestArrayElementValidation(t *testing.T) {
	schema := Array(String().Email())

	data := []interface{}{
		"user1@example.com",
		"user2@example.com",
		"user3@example.com",
	}

	result := schema.Parse(data)
	if !result.Ok {
		t.Errorf("Expected all valid emails to pass. Errors: %v", result.Errors)
	}
}

// Test array in object
func TestArrayInObject(t *testing.T) {
	schema := Object(Schema{
		"name": String(),
		"tags": Array(String()),
	})

	data := map[string]interface{}{
		"name": "John",
		"tags": []interface{}{"developer", "golang", "zogo"},
	}

	result := schema.Parse(data)
	if !result.Ok {
		t.Errorf("Expected object with array to pass. Errors: %v", result.Errors)
	}
}

// Test array in object with error path
func TestArrayInObjectErrorPath(t *testing.T) {
	schema := Object(Schema{
		"users": Array(Object(Schema{
			"email": String().Email(),
		})),
	})

	data := map[string]interface{}{
		"users": []interface{}{
			map[string]interface{}{
				"email": "valid@example.com",
			},
			map[string]interface{}{
				"email": "invalid",
			},
		},
	}

	result := schema.Parse(data)
	if result.Ok {
		t.Error("Expected invalid nested array element to fail")
	}

	// Check error path
	if len(result.Errors) == 0 {
		t.Error("Expected errors")
	} else if result.Errors[0].Path != "users[1].email" {
		t.Errorf("Expected error path 'users[1].email', got '%s'", result.Errors[0].Path)
	}
}
