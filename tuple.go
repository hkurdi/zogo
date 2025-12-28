package zogo

import (
	"fmt"
)

// TupleValidator validates fixed-length arrays with typed positions
type TupleValidator struct {
	validators []Validator
	rest       Validator // Optional validator for remaining elements

	// Modifiers
	isRequired bool
	isOptional bool
	isNullable bool
}

// Tuple creates a new tuple validator with the given position validators
func Tuple(validators ...Validator) *TupleValidator {
	return &TupleValidator{
		validators: validators,
	}
}

// Rest sets a validator for additional elements beyond the tuple positions
func (v *TupleValidator) Rest(validator Validator) *TupleValidator {
	v.rest = validator
	return v
}

// Required marks the field as required
func (v *TupleValidator) Required() *TupleValidator {
	v.isRequired = true
	v.isOptional = false
	return v
}

// Optional allows nil values
func (v *TupleValidator) Optional() *TupleValidator {
	v.isOptional = true
	v.isRequired = false
	return v
}

// Nullable allows null values
func (v *TupleValidator) Nullable() *TupleValidator {
	v.isNullable = true
	return v
}

// Parse validates the input value
func (v *TupleValidator) Parse(value any) ParseResult {
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
		return FailureMessage("Expected tuple, received null")
	}

	// Check if value is an array
	arr, ok := value.([]interface{})
	if !ok {
		return FailureMessage("Expected tuple (array), received " + typeof(value))
	}

	// Check length
	expectedLen := len(v.validators)
	actualLen := len(arr)

	// If no rest validator, array must be exact length
	if v.rest == nil && actualLen != expectedLen {
		return FailureMessage(fmt.Sprintf("Expected tuple of length %d, received length %d", expectedLen, actualLen))
	}

	// If rest validator, array must be at least the required length
	if v.rest != nil && actualLen < expectedLen {
		return FailureMessage(fmt.Sprintf("Expected tuple of at least length %d, received length %d", expectedLen, actualLen))
	}

	// Validate each position
	result := make([]interface{}, 0, len(arr))
	var errors ValidationErrors

	// Validate fixed positions
	for i, validator := range v.validators {
		elemResult := validator.Parse(arr[i])

		if !elemResult.Ok {
			// Add tuple index to error path
			for _, err := range elemResult.Errors {
				errors = append(errors, ValidationError{
					Path:    fmt.Sprintf("[%d]%s", i, prependPath(err.Path)),
					Message: err.Message,
					Value:   err.Value,
				})
			}
		} else {
			result = append(result, elemResult.Value)
		}
	}

	// Validate rest elements if rest validator is set
	if v.rest != nil {
		for i := expectedLen; i < actualLen; i++ {
			elemResult := v.rest.Parse(arr[i])

			if !elemResult.Ok {
				// Add tuple index to error path
				for _, err := range elemResult.Errors {
					errors = append(errors, ValidationError{
						Path:    fmt.Sprintf("[%d]%s", i, prependPath(err.Path)),
						Message: err.Message,
						Value:   err.Value,
					})
				}
			} else {
				result = append(result, elemResult.Value)
			}
		}
	}

	// Return errors if any
	if len(errors) > 0 {
		return Failure(errors...)
	}

	return Success(result)
}
