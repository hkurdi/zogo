package zogo

import (
	"fmt"
)

// RecordValidator validates map[string]T where all values are of the same type
type RecordValidator struct {
	keyValidator   Validator
	valueValidator Validator

	// Modifiers
	isRequired bool
	isOptional bool
	isNullable bool
}

// Record creates a new record validator with key and value validators
func Record(keyValidator, valueValidator Validator) *RecordValidator {
	return &RecordValidator{
		keyValidator:   keyValidator,
		valueValidator: valueValidator,
	}
}

// Required marks the field as required
func (v *RecordValidator) Required() *RecordValidator {
	v.isRequired = true
	v.isOptional = false
	return v
}

// Optional allows nil values
func (v *RecordValidator) Optional() *RecordValidator {
	v.isOptional = true
	v.isRequired = false
	return v
}

// Nullable allows null values
func (v *RecordValidator) Nullable() *RecordValidator {
	v.isNullable = true
	return v
}

// Parse validates the input value
func (v *RecordValidator) Parse(value any) ParseResult {
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
		return FailureMessage("Expected record (object), received null")
	}

	// Check if value is a map
	objMap, ok := value.(map[string]interface{})
	if !ok {
		return FailureMessage("Expected record (object), received " + typeof(value))
	}

	// Result map to build
	result := make(map[string]interface{})

	// Track all errors
	var errors ValidationErrors

	// Validate each key-value pair
	for key, val := range objMap {
		// Validate key
		keyResult := v.keyValidator.Parse(key)
		if !keyResult.Ok {
			for _, err := range keyResult.Errors {
				errors = append(errors, ValidationError{
					Path:    fmt.Sprintf("key(%s)%s", key, prependPath(err.Path)),
					Message: err.Message,
					Value:   err.Value,
				})
			}
			continue // Skip this entry if key is invalid
		}

		// Validate value
		valResult := v.valueValidator.Parse(val)
		if !valResult.Ok {
			for _, err := range valResult.Errors {
				errors = append(errors, ValidationError{
					Path:    fmt.Sprintf("%s%s", key, prependPath(err.Path)),
					Message: err.Message,
					Value:   err.Value,
				})
			}
		} else {
			// Use the validated key and value
			validatedKey, ok := keyResult.Value.(string)
			if !ok {
				// Key must be a string for map[string]interface{}
				errors = append(errors, ValidationError{
					Path:    fmt.Sprintf("key(%s)", key),
					Message: "Record key must be a string",
					Value:   keyResult.Value,
				})
			} else {
				result[validatedKey] = valResult.Value
			}
		}
	}

	// Return errors if any
	if len(errors) > 0 {
		return Failure(errors...)
	}

	return Success(result)
}
