package zogo

// BooleanValidator validates boolean values
type BooleanValidator struct {
	// Modifiers
	isRequired bool
	isOptional bool
	isNullable bool
	defaultVal *bool
}

// Boolean creates a new boolean validator
func Boolean() *BooleanValidator {
	return &BooleanValidator{}
}

// Required marks the field as required
func (v *BooleanValidator) Required() *BooleanValidator {
	v.isRequired = true
	v.isOptional = false
	return v
}

// Optional allows nil values
func (v *BooleanValidator) Optional() *BooleanValidator {
	v.isOptional = true
	v.isRequired = false
	return v
}

// Nullable allows null values
func (v *BooleanValidator) Nullable() *BooleanValidator {
	v.isNullable = true
	return v
}

// Default sets a default value if input is nil
func (v *BooleanValidator) Default(val bool) *BooleanValidator {
	v.defaultVal = &val
	return v
}

// Parse validates the input value
func (v *BooleanValidator) Parse(value any) ParseResult {
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
		return FailureMessage("Expected boolean, received null")
	}

	// Check if value is a boolean
	boolVal, ok := value.(bool)
	if !ok {
		return FailureMessage("Expected boolean, received " + typeof(value))
	}

	return Success(boolVal)
}
