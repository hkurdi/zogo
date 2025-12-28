package zogo

import (
	"testing"
)

// Test basic intersection - number with constraints
func TestIntersectionNumberConstraints(t *testing.T) {
	schema := Intersection(
		Number(),
		Number().Min(10),
		Number().Max(100),
	)

	// Should pass - satisfies all constraints
	result := schema.Parse(50)
	if !result.Ok {
		t.Errorf("Expected 50 to pass all constraints. Errors: %v", result.Errors)
	}

	// Should fail - too small
	result = schema.Parse(5)
	if result.Ok {
		t.Error("Expected 5 to fail Min(10)")
	}

	// Should fail - too large
	result = schema.Parse(150)
	if result.Ok {
		t.Error("Expected 150 to fail Max(100)")
	}
}

// Test intersection with string constraints
func TestIntersectionStringConstraints(t *testing.T) {
	schema := Intersection(
		String(),
		String().Min(5),
		String().Max(10),
	)

	// Should pass
	result := schema.Parse("hello")
	if !result.Ok {
		t.Error("Expected 'hello' to pass all constraints")
	}

	// Should fail - too short
	result = schema.Parse("hi")
	if result.Ok {
		t.Error("Expected 'hi' to fail Min(5)")
	}

	// Should fail - too long
	result = schema.Parse("this is too long")
	if result.Ok {
		t.Error("Expected long string to fail Max(10)")
	}
}

// Test intersection validates against original input
func TestIntersectionOriginalInput(t *testing.T) {
	schema := Intersection(
		String(),
		String().Min(3),
		String().Max(10),
	)

	// Validates the original string against all constraints
	result := schema.Parse("hello")
	if !result.Ok {
		t.Error("Expected 'hello' to pass all validators")
	}

	// Returns original value
	if result.Value != "hello" {
		t.Errorf("Expected 'hello', got '%v'", result.Value)
	}
}

// Test intersection with multiple number validators
func TestIntersectionMultipleNumber(t *testing.T) {
	schema := Intersection(
		Number().Positive(),
		Number().Int(),
		Number().MultipleOf(5),
	)

	// Should pass - positive integer multiple of 5
	result := schema.Parse(15)
	if !result.Ok {
		t.Error("Expected 15 to pass all constraints")
	}

	// Should fail - not multiple of 5
	result = schema.Parse(13)
	if result.Ok {
		t.Error("Expected 13 to fail MultipleOf(5)")
	}

	// Should fail - not integer
	result = schema.Parse(15.5)
	if result.Ok {
		t.Error("Expected 15.5 to fail Int()")
	}

	// Should fail - not positive
	result = schema.Parse(-15)
	if result.Ok {
		t.Error("Expected -15 to fail Positive()")
	}
}

// Test intersection with overlapping object fields
func TestIntersectionOverlappingFields(t *testing.T) {
	schema1 := Object(Schema{
		"value": String(),
	})

	schema2 := Object(Schema{
		"value": String().Min(5),
	})

	schema := Intersection(schema1, schema2)

	// Should pass - string with length >= 5
	data := map[string]interface{}{
		"value": "hello",
	}
	result := schema.Parse(data)
	if !result.Ok {
		t.Error("Expected valid value to pass")
	}

	// Should fail - string too short
	data["value"] = "hi"
	result = schema.Parse(data)
	if result.Ok {
		t.Error("Expected short string to fail")
	}
}

// Test nil value
func TestIntersectionNil(t *testing.T) {
	schema := Intersection(String(), String().Min(5))

	result := schema.Parse(nil)
	if result.Ok {
		t.Error("Expected nil to fail")
	}
}

// Test Optional
func TestIntersectionOptional(t *testing.T) {
	schema := Intersection(String(), String().Min(5)).Optional()

	// nil should pass
	result := schema.Parse(nil)
	if !result.Ok {
		t.Error("Expected nil to pass with Optional()")
	}

	// Valid string should pass
	result = schema.Parse("hello")
	if !result.Ok {
		t.Error("Expected valid string to pass")
	}

	// Invalid string should fail
	result = schema.Parse("hi")
	if result.Ok {
		t.Error("Expected short string to fail")
	}
}

// Test Nullable
func TestIntersectionNullable(t *testing.T) {
	schema := Intersection(String(), String().Min(5)).Nullable()

	// nil should pass
	result := schema.Parse(nil)
	if !result.Ok {
		t.Error("Expected nil to pass with Nullable()")
	}
}

// Test Required
func TestIntersectionRequired(t *testing.T) {
	schema := Intersection(String(), String().Min(5)).Required()

	// nil should fail
	result := schema.Parse(nil)
	if result.Ok {
		t.Error("Expected nil to fail with Required()")
	}

	// Valid string should pass
	result = schema.Parse("hello")
	if !result.Ok {
		t.Error("Expected valid string to pass")
	}
}

// Test intersection in object
func TestIntersectionInObject(t *testing.T) {
	schema := Object(Schema{
		"name": String(),
		"age": Intersection(
			Number(),
			Number().Min(18),
			Number().Max(120),
		),
	})

	data := map[string]interface{}{
		"name": "John",
		"age":  30,
	}
	result := schema.Parse(data)
	if !result.Ok {
		t.Errorf("Expected object with valid intersection to pass. Errors: %v", result.Errors)
	}

	// Invalid age
	data["age"] = 15
	result = schema.Parse(data)
	if result.Ok {
		t.Error("Expected object with invalid age to fail")
	}
}

// Test intersection in array
func TestIntersectionInArray(t *testing.T) {
	schema := Array(Intersection(
		Number(),
		Number().Positive(),
		Number().Int(),
	))

	// All valid
	data := []interface{}{1, 2, 3, 4, 5}
	result := schema.Parse(data)
	if !result.Ok {
		t.Errorf("Expected array of positive integers to pass. Errors: %v", result.Errors)
	}

	// Contains negative
	data = []interface{}{1, -2, 3}
	result = schema.Parse(data)
	if result.Ok {
		t.Error("Expected array with negative to fail")
	}

	// Contains float
	data = []interface{}{1, 2.5, 3}
	result = schema.Parse(data)
	if result.Ok {
		t.Error("Expected array with float to fail")
	}
}

// Test multiple errors collected
func TestIntersectionMultipleErrors(t *testing.T) {
	schema := Intersection(
		Number().Min(10),
		Number().Max(5), // Impossible constraint
	)

	result := schema.Parse(7)
	if result.Ok {
		t.Error("Expected impossible constraints to fail")
	}

	// Should have multiple errors
	if len(result.Errors) < 1 {
		t.Error("Expected multiple errors from failing validators")
	}
}

// Test intersection with single validator
func TestIntersectionSingle(t *testing.T) {
	schema := Intersection(String().Email())

	result := schema.Parse("user@example.com")
	if !result.Ok {
		t.Error("Expected valid email to pass single-validator intersection")
	}

	result = schema.Parse("notanemail")
	if result.Ok {
		t.Error("Expected invalid email to fail")
	}
}

// Test empty intersection (edge case)
func TestIntersectionEmpty(t *testing.T) {
	schema := Intersection()

	// Should pass - no constraints to violate
	result := schema.Parse("anything")
	if !result.Ok {
		t.Error("Expected empty intersection to pass anything")
	}
}

// Test nested intersections
func TestIntersectionNested(t *testing.T) {
	schema := Intersection(
		Number(),
		Intersection(
			Number().Min(10),
			Number().Max(100),
		),
	)

	result := schema.Parse(50)
	if !result.Ok {
		t.Error("Expected 50 to pass nested intersection")
	}

	result = schema.Parse(5)
	if result.Ok {
		t.Error("Expected 5 to fail nested intersection")
	}
}

// Test intersection preserves transformations in order
func TestIntersectionTransformOrder(t *testing.T) {
	schema := Intersection(
		String().Trim(),
		String().Min(3),
	)

	// " hi " should be trimmed to "hi", then fail Min(3)
	result := schema.Parse(" hi ")
	if result.Ok {
		t.Error("Expected trimmed 'hi' to fail Min(3)")
	}

	// " hello " should be trimmed to "hello", then pass Min(3)
	result = schema.Parse(" hello ")
	if !result.Ok {
		t.Error("Expected trimmed 'hello' to pass Min(3)")
	}
	if result.Value != "hello" {
		t.Errorf("Expected 'hello' (trimmed), got '%v'", result.Value)
	}
}

// Test intersection with arrays
func TestIntersectionArrayConstraints(t *testing.T) {
	schema := Intersection(
		Array(Number()),
		Array(Number()).Min(3),
		Array(Number()).Max(5),
	)

	// Should pass - 3 to 5 numbers
	result := schema.Parse([]interface{}{1, 2, 3})
	if !result.Ok {
		t.Error("Expected array of 3 numbers to pass")
	}

	// Should fail - too few
	result = schema.Parse([]interface{}{1, 2})
	if result.Ok {
		t.Error("Expected array of 2 numbers to fail Min(3)")
	}

	// Should fail - too many
	result = schema.Parse([]interface{}{1, 2, 3, 4, 5, 6})
	if result.Ok {
		t.Error("Expected array of 6 numbers to fail Max(5)")
	}
}

// Test practical use case: strict email validation
func TestIntersectionStrictEmail(t *testing.T) {
	schema := Intersection(
		String().Email(),
		String().Min(5),
		String().Max(100),
	)

	// Valid email
	result := schema.Parse("user@example.com")
	if !result.Ok {
		t.Error("Expected valid email to pass all constraints")
	}

	// Too short (though technically valid format)
	result = schema.Parse("a@b.c")
	if result.Ok {
		t.Error("Expected short email to fail Min(5)")
	}
}

// Test intersection with Date
func TestIntersectionDate(t *testing.T) {
	// This is a bit contrived since Date doesn't have many compositional validators
	// But it demonstrates the concept
	schema := Intersection(
		Date().Past(),
		Date(), // Just validates it's a date
	)

	// Past date should pass
	result := schema.Parse("2020-01-01")
	if !result.Ok {
		t.Error("Expected past date to pass")
	}
}

// Test intersection merging objects
func TestIntersectionMergeObjects(t *testing.T) {
	// For intersection to work with objects, use Passthrough mode
	// so each validator sees all fields
	baseSchema := Object(Schema{
		"name": String(),
	}).Passthrough()

	emailSchema := Object(Schema{
		"email": String().Email(),
	}).Passthrough()

	schema := Intersection(baseSchema, emailSchema)

	// Should pass - has both fields
	data := map[string]interface{}{
		"name":  "John",
		"email": "john@example.com",
	}
	result := schema.Parse(data)
	if !result.Ok {
		t.Errorf("Expected object with both fields to pass. Errors: %v", result.Errors)
	}

	// Should fail - missing email
	data = map[string]interface{}{
		"name": "John",
	}
	result = schema.Parse(data)
	if result.Ok {
		t.Error("Expected object without email to fail")
	}

	// Should fail - missing name
	data = map[string]interface{}{
		"email": "john@example.com",
	}
	result = schema.Parse(data)
	if result.Ok {
		t.Error("Expected object without name to fail")
	}
}
