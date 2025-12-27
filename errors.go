package zogo

import "fmt"

// ValidationError represents a single validation failure
type ValidationError struct {
	Path    string // Field path, e.g., "user.email" or "items.0"
	Message string // Human-readable error message
	Value   any    // The value that failed validation
}

// Error implements the error interface
func (e ValidationError) Error() string {
	if e.Path != "" {
		return fmt.Sprintf("%s: %s", e.Path, e.Message)
	}
	return e.Message
}

// ValidationErrors is a collection of validation errors
type ValidationErrors []ValidationError

// Error implements the error interface
func (e ValidationErrors) Error() string {
	if len(e) == 0 {
		return "no errors"
	}
	if len(e) == 1 {
		return e[0].Error()
	}
	return fmt.Sprintf("%d validation errors", len(e))
}

// First returns the first error if any, otherwise nil
func (e ValidationErrors) First() *ValidationError {
	if len(e) > 0 {
		return &e[0]
	}
	return nil
}
