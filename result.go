package zogo

// ParseResult represents the result of a validation operation
type ParseResult struct {
	Ok     bool             // Whether validation succeeded
	Value  any              // The validated value (if Ok is true)
	Errors ValidationErrors // Validation errors (if Ok is false)
}

// Success creates a successful parse result
func Success(value any) ParseResult {
	return ParseResult{
		Ok:    true,
		Value: value,
	}
}

// Failure creates a failed parse result with one or more errors
func Failure(errors ...ValidationError) ParseResult {
	return ParseResult{
		Ok:     false,
		Errors: errors,
	}
}

// FailureMessage creates a failed parse result with a simple message
func FailureMessage(message string) ParseResult {
	return Failure(ValidationError{
		Message: message,
	})
}
