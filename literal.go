package zogo

import (
	"fmt"
)

// LiteralValidator validates that a value exactly matches the expected literal value
type LiteralValidator struct {
	expectedValue interface{}

	// Modifiers
	isRequired bool
	isOptional bool
	isNullable bool
}

// Literal creates a new literal validator with the expected value
func Literal(expectedValue interface{}) *LiteralValidator {
	return &LiteralValidator{
		expectedValue: expectedValue,
	}
}

// Required marks the field as required
func (v *LiteralValidator) Required() *LiteralValidator {
	v.isRequired = true
	v.isOptional = false
	return v
}

// Optional allows nil values
func (v *LiteralValidator) Optional() *LiteralValidator {
	v.isOptional = true
	v.isRequired = false
	return v
}

// Nullable allows null values
func (v *LiteralValidator) Nullable() *LiteralValidator {
	v.isNullable = true
	return v
}

// Parse validates the input value
func (v *LiteralValidator) Parse(value any) ParseResult {
	// Handle nil values based on modifiers
	if value == nil {
		// If optional, nil is OK
		if v.isOptional {
			return Success(nil)
		}

		// If nullable, nil is OK
		if v.isNullable {
			return Success(nil)
		}

		// Otherwise, nil is not allowed
		return FailureMessage(fmt.Sprintf("Expected literal value %v, received null", v.expectedValue))
	}

	// Check if value matches expected literal
	if deepEqual(value, v.expectedValue) {
		return Success(value)
	}

	// Value doesn't match
	return FailureMessage(fmt.Sprintf("Invalid literal value. Expected %v, received %v", v.expectedValue, value))
}
