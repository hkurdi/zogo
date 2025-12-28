package zogo

import (
	"math"
	"testing"
)

// Test basic number validation
func TestNumberBasic(t *testing.T) {
	schema := Number()

	result := schema.Parse(42)
	if !result.Ok {
		t.Error("Expected valid number to pass")
	}

	if result.Value != float64(42) {
		t.Errorf("Expected 42.0, got %v", result.Value)
	}
}

// Test different number types
func TestNumberTypes(t *testing.T) {
	schema := Number()

	// Test int
	result := schema.Parse(int(42))
	if !result.Ok || result.Value != float64(42) {
		t.Error("Expected int to pass")
	}

	// Test float64
	result = schema.Parse(float64(42.5))
	if !result.Ok || result.Value != 42.5 {
		t.Error("Expected float64 to pass")
	}

	// Test float32
	result = schema.Parse(float32(42.5))
	if !result.Ok {
		t.Error("Expected float32 to pass")
	}
}

// Test invalid type
func TestNumberInvalidType(t *testing.T) {
	schema := Number()

	result := schema.Parse("not a number")
	if result.Ok {
		t.Error("Expected string to fail number validation")
	}
}

// Test nil value
func TestNumberNil(t *testing.T) {
	schema := Number()

	result := schema.Parse(nil)
	if result.Ok {
		t.Error("Expected nil to fail")
	}
}

// Test Min
func TestNumberMin(t *testing.T) {
	schema := Number().Min(10)

	// Should pass
	result := schema.Parse(10)
	if !result.Ok {
		t.Error("Expected 10 to pass Min(10)")
	}

	result = schema.Parse(15)
	if !result.Ok {
		t.Error("Expected 15 to pass Min(10)")
	}

	// Should fail
	result = schema.Parse(5)
	if result.Ok {
		t.Error("Expected 5 to fail Min(10)")
	}
}

// Test Max
func TestNumberMax(t *testing.T) {
	schema := Number().Max(100)

	// Should pass
	result := schema.Parse(100)
	if !result.Ok {
		t.Error("Expected 100 to pass Max(100)")
	}

	result = schema.Parse(50)
	if !result.Ok {
		t.Error("Expected 50 to pass Max(100)")
	}

	// Should fail
	result = schema.Parse(150)
	if result.Ok {
		t.Error("Expected 150 to fail Max(100)")
	}
}

// Test Min and Max together
func TestNumberMinMax(t *testing.T) {
	schema := Number().Min(10).Max(100)

	// Should pass
	result := schema.Parse(50)
	if !result.Ok {
		t.Error("Expected 50 to pass Min(10).Max(100)")
	}

	// Should fail - too low
	result = schema.Parse(5)
	if result.Ok {
		t.Error("Expected 5 to fail Min(10)")
	}

	// Should fail - too high
	result = schema.Parse(150)
	if result.Ok {
		t.Error("Expected 150 to fail Max(100)")
	}
}

// Test Int
func TestNumberInt(t *testing.T) {
	schema := Number().Int()

	// Should pass
	result := schema.Parse(42)
	if !result.Ok {
		t.Error("Expected integer to pass Int()")
	}

	result = schema.Parse(float64(42))
	if !result.Ok {
		t.Error("Expected 42.0 to pass Int()")
	}

	// Should fail
	result = schema.Parse(42.5)
	if result.Ok {
		t.Error("Expected 42.5 to fail Int()")
	}
}

// Test Positive
func TestNumberPositive(t *testing.T) {
	schema := Number().Positive()

	// Should pass
	result := schema.Parse(1)
	if !result.Ok {
		t.Error("Expected 1 to pass Positive()")
	}

	// Should fail
	result = schema.Parse(0)
	if result.Ok {
		t.Error("Expected 0 to fail Positive()")
	}

	result = schema.Parse(-1)
	if result.Ok {
		t.Error("Expected -1 to fail Positive()")
	}
}

// Test Negative
func TestNumberNegative(t *testing.T) {
	schema := Number().Negative()

	// Should pass
	result := schema.Parse(-1)
	if !result.Ok {
		t.Error("Expected -1 to pass Negative()")
	}

	// Should fail
	result = schema.Parse(0)
	if result.Ok {
		t.Error("Expected 0 to fail Negative()")
	}

	result = schema.Parse(1)
	if result.Ok {
		t.Error("Expected 1 to fail Negative()")
	}
}

// Test NonNegative
func TestNumberNonNegative(t *testing.T) {
	schema := Number().NonNegative()

	// Should pass
	result := schema.Parse(0)
	if !result.Ok {
		t.Error("Expected 0 to pass NonNegative()")
	}

	result = schema.Parse(1)
	if !result.Ok {
		t.Error("Expected 1 to pass NonNegative()")
	}

	// Should fail
	result = schema.Parse(-1)
	if result.Ok {
		t.Error("Expected -1 to fail NonNegative()")
	}
}

// Test NonPositive
func TestNumberNonPositive(t *testing.T) {
	schema := Number().NonPositive()

	// Should pass
	result := schema.Parse(0)
	if !result.Ok {
		t.Error("Expected 0 to pass NonPositive()")
	}

	result = schema.Parse(-1)
	if !result.Ok {
		t.Error("Expected -1 to pass NonPositive()")
	}

	// Should fail
	result = schema.Parse(1)
	if result.Ok {
		t.Error("Expected 1 to fail NonPositive()")
	}
}

// Test Finite
func TestNumberFinite(t *testing.T) {
	schema := Number().Finite()

	// Should pass
	result := schema.Parse(42)
	if !result.Ok {
		t.Error("Expected 42 to pass Finite()")
	}

	// Should fail - Infinity
	result = schema.Parse(math.Inf(1))
	if result.Ok {
		t.Error("Expected Infinity to fail Finite()")
	}

	// Should fail - NaN
	result = schema.Parse(math.NaN())
	if result.Ok {
		t.Error("Expected NaN to fail Finite()")
	}
}

// Test Safe
func TestNumberSafe(t *testing.T) {
	schema := Number().Safe()

	// Should pass
	result := schema.Parse(42)
	if !result.Ok {
		t.Error("Expected 42 to pass Safe()")
	}

	// Should pass - max safe int
	result = schema.Parse(float64(9007199254740991))
	if !result.Ok {
		t.Error("Expected max safe int to pass Safe()")
	}

	// Should fail - beyond safe range
	result = schema.Parse(float64(9007199254740992))
	if result.Ok {
		t.Error("Expected number beyond safe range to fail Safe()")
	}
}

// Test MultipleOf
func TestNumberMultipleOf(t *testing.T) {
	schema := Number().MultipleOf(5)

	// Should pass
	result := schema.Parse(10)
	if !result.Ok {
		t.Error("Expected 10 to pass MultipleOf(5)")
	}

	result = schema.Parse(15)
	if !result.Ok {
		t.Error("Expected 15 to pass MultipleOf(5)")
	}

	// Should fail
	result = schema.Parse(7)
	if result.Ok {
		t.Error("Expected 7 to fail MultipleOf(5)")
	}
}

// Test Optional
func TestNumberOptional(t *testing.T) {
	schema := Number().Optional()

	// nil should pass
	result := schema.Parse(nil)
	if !result.Ok {
		t.Error("Expected nil to pass with Optional()")
	}

	// Valid number should still pass
	result = schema.Parse(42)
	if !result.Ok {
		t.Error("Expected valid number to pass with Optional()")
	}
}

// Test Nullable
func TestNumberNullable(t *testing.T) {
	schema := Number().Nullable()

	// nil should pass
	result := schema.Parse(nil)
	if !result.Ok {
		t.Error("Expected nil to pass with Nullable()")
	}
}

// Test Default
func TestNumberDefault(t *testing.T) {
	schema := Number().Default(99)

	// nil should return default
	result := schema.Parse(nil)
	if !result.Ok {
		t.Error("Expected nil to pass with Default()")
	}
	if result.Value != float64(99) {
		t.Errorf("Expected 99, got %v", result.Value)
	}

	// Provided value should override default
	result = schema.Parse(42)
	if !result.Ok {
		t.Error("Expected valid number to pass")
	}
	if result.Value != float64(42) {
		t.Errorf("Expected 42, got %v", result.Value)
	}
}

// Test Refine
func TestNumberRefine(t *testing.T) {
	// Must be even
	schema := Number().Refine(func(n float64) bool {
		return int(n)%2 == 0
	}, "Number must be even")

	// Should pass
	result := schema.Parse(42)
	if !result.Ok {
		t.Error("Expected even number to pass")
	}

	// Should fail
	result = schema.Parse(43)
	if result.Ok {
		t.Error("Expected odd number to fail")
	}
	if len(result.Errors) == 0 || result.Errors[0].Message != "Number must be even" {
		t.Error("Expected custom error message")
	}
}

// Test chaining validators
func TestNumberChained(t *testing.T) {
	schema := Number().Min(1).Max(100).Int().Positive()

	// Should pass
	result := schema.Parse(50)
	if !result.Ok {
		t.Error("Expected 50 to pass all validators")
	}

	// Should fail - not int
	result = schema.Parse(50.5)
	if result.Ok {
		t.Error("Expected 50.5 to fail Int()")
	}

	// Should fail - too low
	result = schema.Parse(0)
	if result.Ok {
		t.Error("Expected 0 to fail Positive()")
	}

	// Should fail - too high
	result = schema.Parse(150)
	if result.Ok {
		t.Error("Expected 150 to fail Max(100)")
	}
}
