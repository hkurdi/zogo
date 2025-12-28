package zogo

// UnknownValidator accepts any value without validation (safer alternative to Any)
type UnknownValidator struct {
	// Modifiers
	isRequired bool
	isOptional bool
	isNullable bool
}

// Unknown creates a new unknown validator that accepts any value
func Unknown() *UnknownValidator {
	return &UnknownValidator{}
}

// Required marks the field as required (rejects only nil)
func (v *UnknownValidator) Required() *UnknownValidator {
	v.isRequired = true
	v.isOptional = false
	return v
}

// Optional allows nil values (default behavior for Unknown)
func (v *UnknownValidator) Optional() *UnknownValidator {
	v.isOptional = true
	v.isRequired = false
	return v
}

// Nullable allows null values (default behavior for Unknown)
func (v *UnknownValidator) Nullable() *UnknownValidator {
	v.isNullable = true
	return v
}

// Parse accepts any value
func (v *UnknownValidator) Parse(value any) ParseResult {
	// If Required is explicitly set and value is nil, reject
	if v.isRequired && value == nil {
		return FailureMessage("Expected value, received null")
	}

	// Accept everything else
	return Success(value)
}
