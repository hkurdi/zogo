package zogo

import (
	"testing"
)

// Test basic record - string keys to number values
func TestRecordBasic(t *testing.T) {
	schema := Record(String(), Number())

	data := map[string]interface{}{
		"score1": 100,
		"score2": 85,
		"score3": 92,
	}

	result := schema.Parse(data)
	if !result.Ok {
		t.Errorf("Expected valid record to pass. Errors: %v", result.Errors)
	}

	resultMap := result.Value.(map[string]interface{})
	if len(resultMap) != 3 {
		t.Errorf("Expected 3 entries, got %d", len(resultMap))
	}
	if resultMap["score1"] != float64(100) {
		t.Error("Expected score1 to be preserved")
	}
}

// Test record with string values
func TestRecordStringValues(t *testing.T) {
	schema := Record(String(), String())

	data := map[string]interface{}{
		"en": "Hello",
		"es": "Hola",
		"fr": "Bonjour",
	}

	result := schema.Parse(data)
	if !result.Ok {
		t.Error("Expected string record to pass")
	}
}

// Test record with object values
func TestRecordObjectValues(t *testing.T) {
	schema := Record(String(), Object(Schema{
		"name": String(),
		"age":  Number(),
	}))

	data := map[string]interface{}{
		"user1": map[string]interface{}{
			"name": "John",
			"age":  30,
		},
		"user2": map[string]interface{}{
			"name": "Jane",
			"age":  25,
		},
	}

	result := schema.Parse(data)
	if !result.Ok {
		t.Errorf("Expected record of objects to pass. Errors: %v", result.Errors)
	}
}

// Test record with Any values
func TestRecordAnyValues(t *testing.T) {
	schema := Record(String(), Any())

	data := map[string]interface{}{
		"string": "value",
		"number": 42,
		"bool":   true,
		"array":  []interface{}{1, 2, 3},
	}

	result := schema.Parse(data)
	if !result.Ok {
		t.Error("Expected record with Any values to pass")
	}
}

// Test invalid value type
func TestRecordInvalidValue(t *testing.T) {
	schema := Record(String(), Number())

	data := map[string]interface{}{
		"score1": 100,
		"score2": "not a number", // invalid
		"score3": 92,
	}

	result := schema.Parse(data)
	if result.Ok {
		t.Error("Expected record with invalid value to fail")
	}

	// Check error path
	if len(result.Errors) == 0 || result.Errors[0].Path != "score2" {
		t.Errorf("Expected error path 'score2', got '%s'", result.Errors[0].Path)
	}
}

// Test value constraints
func TestRecordValueConstraints(t *testing.T) {
	schema := Record(String(), Number().Min(0).Max(100))

	// All valid
	data := map[string]interface{}{
		"a": 50,
		"b": 75,
		"c": 100,
	}
	result := schema.Parse(data)
	if !result.Ok {
		t.Error("Expected valid scores to pass")
	}

	// One invalid - too high
	data["d"] = 150
	result = schema.Parse(data)
	if result.Ok {
		t.Error("Expected score > 100 to fail")
	}

	// One invalid - too low
	data = map[string]interface{}{
		"a": 50,
		"b": -10,
	}
	result = schema.Parse(data)
	if result.Ok {
		t.Error("Expected negative score to fail")
	}
}

// Test key constraints
func TestRecordKeyConstraints(t *testing.T) {
	// Keys must be valid emails
	schema := Record(String().Email(), String())

	data := map[string]interface{}{
		"user@example.com": "John",
		"admin@test.com":   "Jane",
	}

	result := schema.Parse(data)
	if !result.Ok {
		t.Errorf("Expected valid email keys to pass. Errors: %v", result.Errors)
	}

	// Invalid key
	data["notanemail"] = "Bob"
	result = schema.Parse(data)
	if result.Ok {
		t.Error("Expected invalid email key to fail")
	}
}

// Test empty record
func TestRecordEmpty(t *testing.T) {
	schema := Record(String(), Number())

	data := map[string]interface{}{}

	result := schema.Parse(data)
	if !result.Ok {
		t.Error("Expected empty record to pass")
	}

	resultMap := result.Value.(map[string]interface{})
	if len(resultMap) != 0 {
		t.Error("Expected empty result map")
	}
}

// Test nil value
func TestRecordNil(t *testing.T) {
	schema := Record(String(), Number())

	result := schema.Parse(nil)
	if result.Ok {
		t.Error("Expected nil to fail")
	}
}

// Test Optional
func TestRecordOptional(t *testing.T) {
	schema := Record(String(), Number()).Optional()

	// nil should pass
	result := schema.Parse(nil)
	if !result.Ok {
		t.Error("Expected nil to pass with Optional()")
	}

	// Valid record should pass
	data := map[string]interface{}{"a": 1}
	result = schema.Parse(data)
	if !result.Ok {
		t.Error("Expected valid record to pass with Optional()")
	}
}

// Test Nullable
func TestRecordNullable(t *testing.T) {
	schema := Record(String(), Number()).Nullable()

	// nil should pass
	result := schema.Parse(nil)
	if !result.Ok {
		t.Error("Expected nil to pass with Nullable()")
	}
}

// Test Required
func TestRecordRequired(t *testing.T) {
	schema := Record(String(), Number()).Required()

	// nil should fail
	result := schema.Parse(nil)
	if result.Ok {
		t.Error("Expected nil to fail with Required()")
	}

	// Valid record should pass
	data := map[string]interface{}{"a": 1}
	result = schema.Parse(data)
	if !result.Ok {
		t.Error("Expected valid record to pass with Required()")
	}
}

// Test invalid type (not object)
func TestRecordInvalidType(t *testing.T) {
	schema := Record(String(), Number())

	result := schema.Parse("not an object")
	if result.Ok {
		t.Error("Expected string to fail record validation")
	}

	result = schema.Parse(42)
	if result.Ok {
		t.Error("Expected number to fail record validation")
	}

	result = schema.Parse([]interface{}{1, 2, 3})
	if result.Ok {
		t.Error("Expected array to fail record validation")
	}
}

// Test record in object
func TestRecordInObject(t *testing.T) {
	schema := Object(Schema{
		"name":   String(),
		"scores": Record(String(), Number()),
	})

	data := map[string]interface{}{
		"name": "Student",
		"scores": map[string]interface{}{
			"math":    95,
			"english": 87,
			"science": 92,
		},
	}

	result := schema.Parse(data)
	if !result.Ok {
		t.Errorf("Expected object with record to pass. Errors: %v", result.Errors)
	}
}

// Test record in array
func TestRecordInArray(t *testing.T) {
	schema := Array(Record(String(), Number()))

	data := []interface{}{
		map[string]interface{}{"a": 1, "b": 2},
		map[string]interface{}{"c": 3, "d": 4},
	}

	result := schema.Parse(data)
	if !result.Ok {
		t.Errorf("Expected array of records to pass. Errors: %v", result.Errors)
	}
}

// Test multiple errors
func TestRecordMultipleErrors(t *testing.T) {
	schema := Record(String(), Number().Min(0))

	data := map[string]interface{}{
		"a": 10,
		"b": -5,  // invalid
		"c": "x", // invalid
		"d": -10, // invalid
	}

	result := schema.Parse(data)
	if result.Ok {
		t.Error("Expected record with multiple invalid values to fail")
	}

	// Should have 3 errors
	if len(result.Errors) != 3 {
		t.Errorf("Expected 3 errors, got %d", len(result.Errors))
	}
}

// Test nested record values
func TestRecordNestedValues(t *testing.T) {
	schema := Record(String(), Record(String(), Number()))

	data := map[string]interface{}{
		"group1": map[string]interface{}{
			"a": 1,
			"b": 2,
		},
		"group2": map[string]interface{}{
			"c": 3,
			"d": 4,
		},
	}

	result := schema.Parse(data)
	if !result.Ok {
		t.Errorf("Expected nested record to pass. Errors: %v", result.Errors)
	}
}

// Test common use case: scores/grades
func TestRecordScores(t *testing.T) {
	scoresSchema := Record(String(), Number().Min(0).Max(100))

	data := map[string]interface{}{
		"Math":    95,
		"Science": 88,
		"English": 92,
	}

	result := scoresSchema.Parse(data)
	if !result.Ok {
		t.Error("Expected valid scores to pass")
	}
}

// Test common use case: configuration/settings
func TestRecordConfig(t *testing.T) {
	configSchema := Record(String(), Union(String(), Number(), Boolean()))

	data := map[string]interface{}{
		"host":    "localhost",
		"port":    8080,
		"enabled": true,
		"timeout": 30,
	}

	result := configSchema.Parse(data)
	if !result.Ok {
		t.Error("Expected config record to pass")
	}
}

// Test common use case: user metadata
func TestRecordMetadata(t *testing.T) {
	metadataSchema := Record(String(), Any())

	data := map[string]interface{}{
		"theme":         "dark",
		"notifications": true,
		"fontSize":      14,
		"tags":          []interface{}{"developer", "admin"},
	}

	result := metadataSchema.Parse(data)
	if !result.Ok {
		t.Error("Expected metadata record to pass")
	}
}

// Test common use case: translations
func TestRecordTranslations(t *testing.T) {
	translationSchema := Record(String(), String())

	data := map[string]interface{}{
		"en": "Hello",
		"es": "Hola",
		"fr": "Bonjour",
		"de": "Hallo",
		"ja": "こんにちは",
	}

	result := translationSchema.Parse(data)
	if !result.Ok {
		t.Error("Expected translation record to pass")
	}
}

// Test with enum values
func TestRecordEnumValues(t *testing.T) {
	schema := Record(String(), Enum([]interface{}{"active", "inactive", "pending"}))

	data := map[string]interface{}{
		"user1": "active",
		"user2": "pending",
		"user3": "inactive",
	}

	result := schema.Parse(data)
	if !result.Ok {
		t.Error("Expected record with enum values to pass")
	}

	// Invalid enum value
	data["user4"] = "deleted"
	result = schema.Parse(data)
	if result.Ok {
		t.Error("Expected invalid enum value to fail")
	}
}
