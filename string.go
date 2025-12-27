package zogo

import (
	"fmt"
	"regexp"
)

type StringValidator struct {
	// Validation rules
	minLen   *int
	maxLen   *int
	exactLen *int
	pattern  *regexp.Regexp

	// Format validators
	isEmail    bool
	isURL      bool
	isUUID     bool
	startsWith *string
	endsWith   *string
	contains   *string

	// Transformations
	shouldTrim      bool
	shouldLowercase bool
	shouldUppercase bool

	// Modifiers
	isRequired bool
	isOptional bool
	isNullable bool
	defaultVal *string

	// Custom validators
	refinements []Refinement
}

type Refinement struct {
	Check   func(string) bool
	Message string
}

// String creates a new string validator
func String() *StringValidator {
	return &StringValidator{}
}

// Min sets the minimum string length
func (v *StringValidator) Min(length int) *StringValidator {
	v.minLen = &length
	return v
}

// Max sets the maximum string length
func (v *StringValidator) Max(length int) *StringValidator {
	v.maxLen = &length
	return v
}

// Length sets the exact string length required
func (v *StringValidator) Length(length int) *StringValidator {
	v.exactLen = &length
	return v
}

// Parse validates the input value
func (v *StringValidator) Parse(value any) ParseResult {
	// Check if value is nil
	if value == nil {
		return FailureMessage("Expected string, received null")
	}

	// Check if value is a string
	str, ok := value.(string)
	if !ok {
		return FailureMessage("Expected string, received " + typeof(value))
	}

	// Check exact length if specified
	if v.exactLen != nil && len(str) != *v.exactLen {
		return FailureMessage(fmt.Sprintf("String must be exactly %d characters", *v.exactLen))
	}

	// Check minimum length
	if v.minLen != nil && len(str) < *v.minLen {
		return FailureMessage(fmt.Sprintf("String must be at least %d characters", *v.minLen))
	}

	// Check maximum length
	if v.maxLen != nil && len(str) > *v.maxLen {
		return FailureMessage(fmt.Sprintf("String must be at most %d characters", *v.maxLen))
	}

	return Success(str)
}

// Helper function to get type name
func typeof(value any) string {
	if value == nil {
		return "null"
	}
	switch value.(type) {
	case string:
		return "string"
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return "number"
	case float32, float64:
		return "number"
	case bool:
		return "boolean"
	default:
		return "unknown"
	}
}
