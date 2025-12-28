package zogo

import "testing"

// Test basic boolean validation
func TestBooleanBasic(t *testing.T) {
	schema := Boolean()

	// Test true
	result := schema.Parse(true)
	if !result.Ok {
		t.Error("Expected true to pass")
	}
	if result.Value != true {
		t.Errorf("Expected true, got %v", result.Value)
	}

	// Test false
	result = schema.Parse(false)
	if !result.Ok {
		t.Error("Expected false to pass")
	}
	if result.Value != false {
		t.Errorf("Expected false, got %v", result.Value)
	}
}

// Test invalid type
func TestBooleanInvalidType(t *testing.T) {
	schema := Boolean()

	// String should fail
	result := schema.Parse("true")
	if result.Ok {
		t.Error("Expected string to fail boolean validation")
	}

	// Number should fail
	result = schema.Parse(1)
	if result.Ok {
		t.Error("Expected number to fail boolean validation")
	}
}

// Test nil value
func TestBooleanNil(t *testing.T) {
	schema := Boolean()

	result := schema.Parse(nil)
	if result.Ok {
		t.Error("Expected nil to fail")
	}
}

// Test Optional
func TestBooleanOptional(t *testing.T) {
	schema := Boolean().Optional()

	// nil should pass
	result := schema.Parse(nil)
	if !result.Ok {
		t.Error("Expected nil to pass with Optional()")
	}

	// Valid boolean should still pass
	result = schema.Parse(true)
	if !result.Ok {
		t.Error("Expected true to pass with Optional()")
	}
}

// Test Nullable
func TestBooleanNullable(t *testing.T) {
	schema := Boolean().Nullable()

	// nil should pass
	result := schema.Parse(nil)
	if !result.Ok {
		t.Error("Expected nil to pass with Nullable()")
	}

	// Valid boolean should still pass
	result = schema.Parse(false)
	if !result.Ok {
		t.Error("Expected false to pass with Nullable()")
	}
}

// Test Default with true
func TestBooleanDefaultTrue(t *testing.T) {
	schema := Boolean().Default(true)

	// nil should return default
	result := schema.Parse(nil)
	if !result.Ok {
		t.Error("Expected nil to pass with Default()")
	}
	if result.Value != true {
		t.Errorf("Expected true, got %v", result.Value)
	}

	// Provided value should override default
	result = schema.Parse(false)
	if !result.Ok {
		t.Error("Expected false to pass")
	}
	if result.Value != false {
		t.Errorf("Expected false, got %v", result.Value)
	}
}

// Test Default with false
func TestBooleanDefaultFalse(t *testing.T) {
	schema := Boolean().Default(false)

	// nil should return default
	result := schema.Parse(nil)
	if !result.Ok {
		t.Error("Expected nil to pass with Default()")
	}
	if result.Value != false {
		t.Errorf("Expected false, got %v", result.Value)
	}

	// Provided value should override default
	result = schema.Parse(true)
	if !result.Ok {
		t.Error("Expected true to pass")
	}
	if result.Value != true {
		t.Errorf("Expected true, got %v", result.Value)
	}
}

// Test Required
func TestBooleanRequired(t *testing.T) {
	schema := Boolean().Required()

	// nil should fail
	result := schema.Parse(nil)
	if result.Ok {
		t.Error("Expected nil to fail with Required()")
	}

	// Valid boolean should pass
	result = schema.Parse(true)
	if !result.Ok {
		t.Error("Expected true to pass with Required()")
	}
}
