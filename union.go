package zogo

import (
	"fmt"
	"strings"
)

// UnionValidator validates that a value matches at least one of the provided validators
type UnionValidator struct {
	validators []Validator

	// Modifiers
	isRequired bool
	isOptional bool
	isNullable bool
}

// Union creates a new union validator with the given validators
func Union(validators ...Validator) *UnionValidator {
	return &UnionValidator{
		validators: validators,
	}
}

// Required marks the field as required
func (v *UnionValidator) Required() *UnionValidator {
	v.isRequired = true
	v.isOptional = false
	return v
}

// Optional allows nil values
func (v *UnionValidator) Optional() *UnionValidator {
	v.isOptional = true
	v.isRequired = false
	return v
}

// Nullable allows null values
func (v *UnionValidator) Nullable() *UnionValidator {
	v.isNullable = true
	return v
}

// Parse validates the input value against all union members
func (v *UnionValidator) Parse(value any) ParseResult {
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

		// If explicitly required, reject
		if v.isRequired {
			return FailureMessage("Expected value, received null")
		}

		// Otherwise, try validating nil through each validator
		// (some validators like String().Optional() might accept it)
	}

	// Try each validator in the union
	var allErrors []string

	for i, validator := range v.validators {
		result := validator.Parse(value)

		// If any validator passes, return success immediately
		if result.Ok {
			return Success(result.Value)
		}

		// Collect error messages for reporting
		if len(result.Errors) > 0 {
			// Format error for this union member
			errorMsgs := make([]string, len(result.Errors))
			for j, err := range result.Errors {
				errorMsgs[j] = err.Message
			}
			allErrors = append(allErrors, fmt.Sprintf("Option %d: %s", i+1, strings.Join(errorMsgs, ", ")))
		} else {
			allErrors = append(allErrors, fmt.Sprintf("Option %d: validation failed", i+1))
		}
	}

	// None of the validators passed
	errorMsg := fmt.Sprintf("Value did not match any union type. Errors: %s", strings.Join(allErrors, "; "))
	return FailureMessage(errorMsg)
}
