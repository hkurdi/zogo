package zogo

import (
	"fmt"
	"math"
)

// NumberValidator validates number values with chainable methods
type NumberValidator struct {
	// Validation rules
	minVal     *float64
	maxVal     *float64
	multipleOf *float64

	// Type checks
	isInt         bool
	isPositive    bool
	isNegative    bool
	isNonNegative bool
	isNonPositive bool
	isFinite      bool
	isSafe        bool

	// Modifiers
	isRequired bool
	isOptional bool
	isNullable bool
	defaultVal *float64

	// Custom validators
	refinements []NumberRefinement
}

// NumberRefinement holds custom validation logic for numbers
type NumberRefinement struct {
	Check   func(float64) bool
	Message string
}

// Number creates a new number validator
func Number() *NumberValidator {
	return &NumberValidator{}
}

// Min sets the minimum value
func (v *NumberValidator) Min(val float64) *NumberValidator {
	v.minVal = &val
	return v
}

// Max sets the maximum value
func (v *NumberValidator) Max(val float64) *NumberValidator {
	v.maxVal = &val
	return v
}

// Int requires the number to be an integer
func (v *NumberValidator) Int() *NumberValidator {
	v.isInt = true
	return v
}

// Positive requires number > 0
func (v *NumberValidator) Positive() *NumberValidator {
	v.isPositive = true
	return v
}

// Negative requires number < 0
func (v *NumberValidator) Negative() *NumberValidator {
	v.isNegative = true
	return v
}

// NonNegative requires number >= 0
func (v *NumberValidator) NonNegative() *NumberValidator {
	v.isNonNegative = true
	return v
}

// NonPositive requires number <= 0
func (v *NumberValidator) NonPositive() *NumberValidator {
	v.isNonPositive = true
	return v
}

// Finite disallows Infinity and NaN
func (v *NumberValidator) Finite() *NumberValidator {
	v.isFinite = true
	return v
}

// Safe requires number to be within safe integer range
func (v *NumberValidator) Safe() *NumberValidator {
	v.isSafe = true
	return v
}

// MultipleOf requires number to be a multiple of the given value
func (v *NumberValidator) MultipleOf(val float64) *NumberValidator {
	v.multipleOf = &val
	return v
}

// Required marks the field as required
func (v *NumberValidator) Required() *NumberValidator {
	v.isRequired = true
	v.isOptional = false
	return v
}

// Optional allows nil values
func (v *NumberValidator) Optional() *NumberValidator {
	v.isOptional = true
	v.isRequired = false
	return v
}

// Nullable allows null values
func (v *NumberValidator) Nullable() *NumberValidator {
	v.isNullable = true
	return v
}

// Default sets a default value if input is nil
func (v *NumberValidator) Default(val float64) *NumberValidator {
	v.defaultVal = &val
	return v
}

// Refine adds custom validation logic
func (v *NumberValidator) Refine(check func(float64) bool, message string) *NumberValidator {
	v.refinements = append(v.refinements, NumberRefinement{
		Check:   check,
		Message: message,
	})
	return v
}

// Parse validates the input value
func (v *NumberValidator) Parse(value any) ParseResult {
	// Handle nil values based on modifiers
	if value == nil {
		// If default is set, use it
		if v.defaultVal != nil {
			return Success(*v.defaultVal)
		}

		// If optional, nil is OK
		if v.isOptional {
			return Success(nil)
		}

		// If nullable, nil is OK
		if v.isNullable {
			return Success(nil)
		}

		// Otherwise, nil is not allowed
		return FailureMessage("Expected number, received null")
	}

	// Convert to float64
	var num float64
	switch v := value.(type) {
	case int:
		num = float64(v)
	case int8:
		num = float64(v)
	case int16:
		num = float64(v)
	case int32:
		num = float64(v)
	case int64:
		num = float64(v)
	case uint:
		num = float64(v)
	case uint8:
		num = float64(v)
	case uint16:
		num = float64(v)
	case uint32:
		num = float64(v)
	case uint64:
		num = float64(v)
	case float32:
		num = float64(v)
	case float64:
		num = v
	default:
		return FailureMessage("Expected number, received " + typeof(value))
	}

	// Check if finite (no Infinity or NaN)
	if v.isFinite && (math.IsInf(num, 0) || math.IsNaN(num)) {
		return FailureMessage("Number must be finite")
	}

	// Check if integer
	if v.isInt && num != math.Floor(num) {
		return FailureMessage("Number must be an integer")
	}

	// Check if safe integer
	if v.isSafe {
		const maxSafeInt = 9007199254740991  // 2^53 - 1
		const minSafeInt = -9007199254740991 // -(2^53 - 1)
		if num > maxSafeInt || num < minSafeInt {
			return FailureMessage("Number must be within safe integer range")
		}
	}

	// Check minimum value
	if v.minVal != nil && num < *v.minVal {
		return FailureMessage(fmt.Sprintf("Number must be at least %v", *v.minVal))
	}

	// Check maximum value
	if v.maxVal != nil && num > *v.maxVal {
		return FailureMessage(fmt.Sprintf("Number must be at most %v", *v.maxVal))
	}

	// Check positive
	if v.isPositive && num <= 0 {
		return FailureMessage("Number must be positive")
	}

	// Check negative
	if v.isNegative && num >= 0 {
		return FailureMessage("Number must be negative")
	}

	// Check non-negative
	if v.isNonNegative && num < 0 {
		return FailureMessage("Number must be non-negative")
	}

	// Check non-positive
	if v.isNonPositive && num > 0 {
		return FailureMessage("Number must be non-positive")
	}

	// Check multiple of
	if v.multipleOf != nil {
		remainder := math.Mod(num, *v.multipleOf)
		// Use small epsilon for floating point comparison
		if math.Abs(remainder) > 1e-10 && math.Abs(remainder-*v.multipleOf) > 1e-10 {
			return FailureMessage(fmt.Sprintf("Number must be a multiple of %v", *v.multipleOf))
		}
	}

	// Run custom refinements
	for _, refinement := range v.refinements {
		if !refinement.Check(num) {
			return FailureMessage(refinement.Message)
		}
	}

	return Success(num)
}
