package zogo

import (
	"fmt"
	"strings"
)

// ValidationError represents a single validation error
type ValidationError struct {
	Path    string // Field path (e.g., "user.email" or "items[0].name")
	Message string // Human-readable error message
	Value   any    // The value that failed validation
	Code    string // Error code (e.g., "invalid_type", "too_small")
}

// Error returns the error message
func (e ValidationError) Error() string {
	if e.Path != "" {
		return fmt.Sprintf("%s: %s", e.Path, e.Message)
	}
	return e.Message
}

// ValidationErrors is a collection of validation errors
type ValidationErrors []ValidationError

// Error returns a formatted string of all errors
func (e ValidationErrors) Error() string {
	if len(e) == 0 {
		return "No errors"
	}

	if len(e) == 1 {
		return e[0].Error()
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Validation failed with %d error(s):\n", len(e)))
	for i, err := range e {
		sb.WriteString(fmt.Sprintf("  %d. %s\n", i+1, err.Error()))
	}
	return sb.String()
}

// First returns the first error, or nil if there are no errors
func (e ValidationErrors) First() *ValidationError {
	if len(e) == 0 {
		return nil
	}
	return &e[0]
}

// HasPath checks if there's an error at the given path
func (e ValidationErrors) HasPath(path string) bool {
	for _, err := range e {
		if err.Path == path {
			return true
		}
	}
	return false
}

// ByPath returns all errors for a given path
func (e ValidationErrors) ByPath(path string) ValidationErrors {
	var result ValidationErrors
	for _, err := range e {
		if err.Path == path {
			result = append(result, err)
		}
	}
	return result
}

// Issues returns a structured list of all issues (useful for JSON responses)
func (e ValidationErrors) Issues() []map[string]interface{} {
	issues := make([]map[string]interface{}, len(e))
	for i, err := range e {
		issues[i] = map[string]interface{}{
			"path":    err.Path,
			"message": err.Message,
			"code":    err.Code,
		}
		if err.Value != nil {
			issues[i]["received"] = err.Value
		}
	}
	return issues
}
