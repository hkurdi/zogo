package zogo

import (
	"fmt"
)

// IntersectionValidator validates that a value matches ALL of the provided validators
type IntersectionValidator struct {
	validators []Validator

	// Modifiers
	isRequired bool
	isOptional bool
	isNullable bool
}

// Intersection creates a new intersection validator with the given validators
func Intersection(validators ...Validator) *IntersectionValidator {
	return &IntersectionValidator{
		validators: validators,
	}
}

// Required marks the field as required
func (v *IntersectionValidator) Required() *IntersectionValidator {
	v.isRequired = true
	v.isOptional = false
	return v
}

// Optional allows nil values
func (v *IntersectionValidator) Optional() *IntersectionValidator {
	v.isOptional = true
	v.isRequired = false
	return v
}

// Nullable allows null values
func (v *IntersectionValidator) Nullable() *IntersectionValidator {
	v.isNullable = true
	return v
}

// Parse validates the input value against all intersection members
func (v *IntersectionValidator) Parse(value any) ParseResult {
	// Handle nil values based on modifiers
	if value == nil {
		if v.isOptional || v.isNullable {
			return Success(nil)
		}
		if v.isRequired {
			return FailureMessage("Expected value, received null")
		}
	}

	var allErrors ValidationErrors

	// Start with the original value
	currentValue := value

	for i, validator := range v.validators {
		// Validate against the current value (which may have been transformed by previous steps)
		result := validator.Parse(currentValue)

		if !result.Ok {
			// If validation fails, collect errors
			for _, err := range result.Errors {
				allErrors = append(allErrors, ValidationError{
					Path:    err.Path,
					Message: fmt.Sprintf("Intersection validator %d: %s", i+1, err.Message),
					Value:   err.Value,
				})
			}
		} else {
			// If validation succeeds, update currentValue to the transformed result
			// This allows chaining: String().Trim() -> passes "trimmed" to next validator
			currentValue = result.Value
		}
	}

	// If any errors occurred, return failure
	if len(allErrors) > 0 {
		return Failure(allErrors...)
	}

	// Return the final transformed value
	return Success(currentValue)
}
