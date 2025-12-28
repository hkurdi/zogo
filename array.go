package zogo

import (
	"fmt"
)

// ArrayValidator validates array/slice values with typed elements
type ArrayValidator struct {
	elementValidator Validator
	minLen           *int
	maxLen           *int
	isNonEmpty       bool

	// Modifiers
	isRequired bool
	isOptional bool
	isNullable bool
}

// Array creates a new array validator with the given element validator
func Array(elementValidator Validator) *ArrayValidator {
	return &ArrayValidator{
		elementValidator: elementValidator,
	}
}

// Min sets the minimum array length
func (v *ArrayValidator) Min(length int) *ArrayValidator {
	v.minLen = &length
	return v
}

// Max sets the maximum array length
func (v *ArrayValidator) Max(length int) *ArrayValidator {
	v.maxLen = &length
	return v
}

// Length sets exact array length (same as Min(n).Max(n))
func (v *ArrayValidator) Length(length int) *ArrayValidator {
	v.minLen = &length
	v.maxLen = &length
	return v
}

// NonEmpty requires array to have at least one element
func (v *ArrayValidator) NonEmpty() *ArrayValidator {
	v.isNonEmpty = true
	return v
}

// Required marks the field as required
func (v *ArrayValidator) Required() *ArrayValidator {
	v.isRequired = true
	v.isOptional = false
	return v
}

// Optional allows nil values
func (v *ArrayValidator) Optional() *ArrayValidator {
	v.isOptional = true
	v.isRequired = false
	return v
}

// Nullable allows null values
func (v *ArrayValidator) Nullable() *ArrayValidator {
	v.isNullable = true
	return v
}

// Parse validates the input value
func (v *ArrayValidator) Parse(value any) ParseResult {
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
		return FailureMessage("Expected array, received null")
	}

	// Check if value is a slice
	arr, ok := value.([]interface{})
	if !ok {
		return FailureMessage("Expected array, received " + typeof(value))
	}

	// Check length constraints
	arrLen := len(arr)

	if v.isNonEmpty && arrLen == 0 {
		return FailureMessage("Array must not be empty")
	}

	if v.minLen != nil && arrLen < *v.minLen {
		return FailureMessage(fmt.Sprintf("Array must contain at least %d element(s)", *v.minLen))
	}

	if v.maxLen != nil && arrLen > *v.maxLen {
		return FailureMessage(fmt.Sprintf("Array must contain at most %d element(s)", *v.maxLen))
	}

	// Validate each element
	result := make([]interface{}, 0, len(arr))
	var errors ValidationErrors

	for i, elem := range arr {
		elemResult := v.elementValidator.Parse(elem)

		if !elemResult.Ok {
			// Add array index to error path
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

	// Return errors if any
	if len(errors) > 0 {
		return Failure(errors...)
	}

	return Success(result)
}
