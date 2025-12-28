package zogo

import (
	"fmt"
	"reflect"
)

// EnumValidator validates that a value is one of the allowed values
type EnumValidator struct {
	allowedValues []interface{}

	// Modifiers
	isRequired bool
	isOptional bool
	isNullable bool
	defaultVal *interface{}
}

// Enum creates a new enum validator with the given allowed values
func Enum(allowedValues []interface{}) *EnumValidator {
	return &EnumValidator{
		allowedValues: allowedValues,
	}
}

// Required marks the field as required
func (v *EnumValidator) Required() *EnumValidator {
	v.isRequired = true
	v.isOptional = false
	return v
}

// Optional allows nil values
func (v *EnumValidator) Optional() *EnumValidator {
	v.isOptional = true
	v.isRequired = false
	return v
}

// Nullable allows null values
func (v *EnumValidator) Nullable() *EnumValidator {
	v.isNullable = true
	return v
}

// Default sets a default value if input is nil
func (v *EnumValidator) Default(val interface{}) *EnumValidator {
	v.defaultVal = &val
	return v
}

// Parse validates the input value
func (v *EnumValidator) Parse(value any) ParseResult {
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
		return FailureMessage("Expected enum value, received null")
	}

	// Check if value is in allowed values
	for _, allowed := range v.allowedValues {
		if deepEqual(value, allowed) {
			return Success(value)
		}
	}

	// Value not found in allowed values
	return FailureMessage(fmt.Sprintf("Invalid enum value. Expected one of: %v, received: %v", v.allowedValues, value))
}

// deepEqual compares two values for equality, handling different numeric types
func deepEqual(a, b interface{}) bool {
	// Use reflect.DeepEqual for most cases
	if reflect.DeepEqual(a, b) {
		return true
	}

	// Handle numeric type conversions
	// This allows comparing int(1) with float64(1), etc.
	aVal := reflect.ValueOf(a)
	bVal := reflect.ValueOf(b)

	// Check if both are numeric
	if isNumeric(aVal.Kind()) && isNumeric(bVal.Kind()) {
		// Convert both to float64 for comparison
		aFloat := toFloat64(a)
		bFloat := toFloat64(b)
		return aFloat == bFloat
	}

	return false
}

// isNumeric checks if a reflect.Kind is a numeric type
func isNumeric(k reflect.Kind) bool {
	switch k {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return true
	}
	return false
}

// toFloat64 converts numeric values to float64
func toFloat64(val interface{}) float64 {
	v := reflect.ValueOf(val)
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(v.Uint())
	case reflect.Float32, reflect.Float64:
		return v.Float()
	}
	return 0
}
