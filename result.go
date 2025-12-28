package zogo

// ParseResult represents the result of a validation
type ParseResult struct {
	Ok     bool
	Value  any
	Errors ValidationErrors
}

// Success creates a successful parse result
func Success(value any) ParseResult {
	return ParseResult{
		Ok:    true,
		Value: value,
	}
}

// Failure creates a failed parse result with errors
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

// FailureWithCode creates a failed parse result with a message and code
func FailureWithCode(message string, code string) ParseResult {
	return Failure(ValidationError{
		Message: message,
		Code:    code,
	})
}

// FailureTypeMismatch creates a type mismatch error
func FailureTypeMismatch(expected string, received any) ParseResult {
	return Failure(ValidationError{
		Message: "Expected " + expected + ", received " + typeof(received),
		Code:    "invalid_type",
		Value:   received,
	})
}
