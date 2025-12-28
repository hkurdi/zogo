package zogo

// AnyValidator accepts any value without validation
type AnyValidator struct {
	// Modifiers (though less meaningful for Any)
	isRequired bool
	isOptional bool
	isNullable bool
}

// Any creates a new any validator that accepts any value
func Any() *AnyValidator {
	return &AnyValidator{}
}

// Required marks the field as required (rejects only nil)
func (v *AnyValidator) Required() *AnyValidator {
	v.isRequired = true
	v.isOptional = false
	return v
}

// Optional allows nil values (default behavior for Any)
func (v *AnyValidator) Optional() *AnyValidator {
	v.isOptional = true
	v.isRequired = false
	return v
}

// Nullable allows null values (default behavior for Any)
func (v *AnyValidator) Nullable() *AnyValidator {
	v.isNullable = true
	return v
}

// Parse accepts any value
func (v *AnyValidator) Parse(value any) ParseResult {
	// If Required is explicitly set and value is nil, reject
	if v.isRequired && value == nil {
		return FailureMessage("Expected value, received null")
	}

	// Accept everything else
	return Success(value)
}
