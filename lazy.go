package zogo

// LazyValidator defers schema construction until validation time
// This enables recursive/self-referential schemas
type LazyValidator struct {
	factory func() Validator

	// Modifiers
	isRequired bool
	isOptional bool
	isNullable bool
}

// Lazy creates a new lazy validator that constructs the actual validator at validation time
// This enables recursive schemas where a validator references itself
func Lazy(factory func() Validator) *LazyValidator {
	return &LazyValidator{
		factory: factory,
	}
}

// Required marks the field as required
func (v *LazyValidator) Required() *LazyValidator {
	v.isRequired = true
	v.isOptional = false
	return v
}

// Optional allows nil values
func (v *LazyValidator) Optional() *LazyValidator {
	v.isOptional = true
	v.isRequired = false
	return v
}

// Nullable allows null values
func (v *LazyValidator) Nullable() *LazyValidator {
	v.isNullable = true
	return v
}

// Parse validates the input value by constructing the actual validator at runtime
func (v *LazyValidator) Parse(value any) ParseResult {
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
	}

	// Construct the actual validator at validation time
	actualValidator := v.factory()

	// Delegate to the actual validator
	return actualValidator.Parse(value)
}
