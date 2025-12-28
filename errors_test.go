package zogo

import (
	"strings"
	"testing"
)

// Test ValidationError.Error()
func TestValidationErrorError(t *testing.T) {
	err := ValidationError{
		Path:    "user.email",
		Message: "Invalid email",
		Code:    "invalid_email",
	}

	expected := "user.email: Invalid email"
	if err.Error() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, err.Error())
	}

	// Error without path
	err2 := ValidationError{
		Message: "Invalid value",
		Code:    "invalid_value",
	}

	if err2.Error() != "Invalid value" {
		t.Errorf("Expected 'Invalid value', got '%s'", err2.Error())
	}
}

// Test ValidationErrors.Error()
func TestValidationErrorsError(t *testing.T) {
	errors := ValidationErrors{
		ValidationError{Path: "name", Message: "Required"},
		ValidationError{Path: "email", Message: "Invalid email"},
	}

	result := errors.Error()

	// Should contain both errors
	if !strings.Contains(result, "name: Required") {
		t.Error("Expected error string to contain name error")
	}
	if !strings.Contains(result, "email: Invalid email") {
		t.Error("Expected error string to contain email error")
	}
	if !strings.Contains(result, "2 error(s)") {
		t.Error("Expected error count in message")
	}
}

// Test ValidationErrors.First()
func TestValidationErrorsFirst(t *testing.T) {
	errors := ValidationErrors{
		ValidationError{Path: "name", Message: "Required"},
		ValidationError{Path: "email", Message: "Invalid email"},
	}

	first := errors.First()
	if first == nil {
		t.Error("Expected first error to not be nil")
	}
	if first.Path != "name" {
		t.Errorf("Expected first error path 'name', got '%s'", first.Path)
	}

	// Empty errors
	emptyErrors := ValidationErrors{}
	if emptyErrors.First() != nil {
		t.Error("Expected First() on empty errors to return nil")
	}
}

// Test ValidationErrors.HasPath()
func TestValidationErrorsHasPath(t *testing.T) {
	errors := ValidationErrors{
		ValidationError{Path: "name", Message: "Required"},
		ValidationError{Path: "email", Message: "Invalid email"},
	}

	if !errors.HasPath("name") {
		t.Error("Expected HasPath('name') to return true")
	}

	if !errors.HasPath("email") {
		t.Error("Expected HasPath('email') to return true")
	}

	if errors.HasPath("age") {
		t.Error("Expected HasPath('age') to return false")
	}
}

// Test ValidationErrors.ByPath()
func TestValidationErrorsByPath(t *testing.T) {
	errors := ValidationErrors{
		ValidationError{Path: "name", Message: "Required"},
		ValidationError{Path: "name", Message: "Too short"},
		ValidationError{Path: "email", Message: "Invalid email"},
	}

	nameErrors := errors.ByPath("name")
	if len(nameErrors) != 2 {
		t.Errorf("Expected 2 errors for 'name', got %d", len(nameErrors))
	}

	emailErrors := errors.ByPath("email")
	if len(emailErrors) != 1 {
		t.Errorf("Expected 1 error for 'email', got %d", len(emailErrors))
	}

	ageErrors := errors.ByPath("age")
	if len(ageErrors) != 0 {
		t.Errorf("Expected 0 errors for 'age', got %d", len(ageErrors))
	}
}

// Test ValidationErrors.Issues()
func TestValidationErrorsIssues(t *testing.T) {
	errors := ValidationErrors{
		ValidationError{
			Path:    "name",
			Message: "Required",
			Code:    "required",
		},
		ValidationError{
			Path:    "email",
			Message: "Invalid email",
			Code:    "invalid_email",
			Value:   "not-an-email",
		},
	}

	issues := errors.Issues()

	if len(issues) != 2 {
		t.Errorf("Expected 2 issues, got %d", len(issues))
	}

	// Check first issue
	if issues[0]["path"] != "name" {
		t.Error("Expected first issue path to be 'name'")
	}
	if issues[0]["message"] != "Required" {
		t.Error("Expected first issue message to be 'Required'")
	}
	if issues[0]["code"] != "required" {
		t.Error("Expected first issue code to be 'required'")
	}

	// Check second issue
	if issues[1]["path"] != "email" {
		t.Error("Expected second issue path to be 'email'")
	}
	if issues[1]["received"] != "not-an-email" {
		t.Error("Expected second issue to have received value")
	}
}

// Test FailureWithCode
func TestFailureWithCode(t *testing.T) {
	result := FailureWithCode("Invalid value", "invalid_value")

	if result.Ok {
		t.Error("Expected result to not be Ok")
	}

	if len(result.Errors) != 1 {
		t.Error("Expected 1 error")
	}

	if result.Errors[0].Code != "invalid_value" {
		t.Errorf("Expected error code 'invalid_value', got '%s'", result.Errors[0].Code)
	}
}

// Test FailureTypeMismatch
func TestFailureTypeMismatch(t *testing.T) {
	result := FailureTypeMismatch("string", 42)

	if result.Ok {
		t.Error("Expected result to not be Ok")
	}

	if len(result.Errors) != 1 {
		t.Error("Expected 1 error")
	}

	if result.Errors[0].Code != "invalid_type" {
		t.Errorf("Expected error code 'invalid_type', got '%s'", result.Errors[0].Code)
	}

	if result.Errors[0].Value != 42 {
		t.Error("Expected error value to be 42")
	}
}

// Test error formatting in real validation
func TestErrorFormattingInValidation(t *testing.T) {
	schema := Object(Schema{
		"name":  String().Min(3),
		"email": String().Email(),
		"age":   Number().Min(18),
	})

	data := map[string]interface{}{
		"name":  "Jo",           // too short
		"email": "not-an-email", // invalid
		"age":   15,             // too young
	}

	result := schema.Parse(data)

	if result.Ok {
		t.Error("Expected validation to fail")
	}

	// Should have 3 errors
	if len(result.Errors) != 3 {
		t.Errorf("Expected 3 errors, got %d", len(result.Errors))
	}

	// Check error paths
	if !result.Errors.HasPath("name") {
		t.Error("Expected error at 'name' path")
	}

	if !result.Errors.HasPath("email") {
		t.Error("Expected error at 'email' path")
	}

	if !result.Errors.HasPath("age") {
		t.Error("Expected error at 'age' path")
	}

	// Test formatted error message
	errorMsg := result.Errors.Error()
	if !strings.Contains(errorMsg, "3 error(s)") {
		t.Error("Expected formatted error to mention error count")
	}
}

// Test single error formatting
func TestSingleErrorFormatting(t *testing.T) {
	schema := String().Email()
	result := schema.Parse("not-an-email")

	if result.Ok {
		t.Error("Expected validation to fail")
	}

	// Single error should not include count
	errorMsg := result.Errors.Error()
	if strings.Contains(errorMsg, "error(s)") {
		t.Error("Single error should not include error count")
	}
}

// Test empty errors
func TestEmptyErrors(t *testing.T) {
	errors := ValidationErrors{}

	if errors.Error() != "No errors" {
		t.Errorf("Expected 'No errors', got '%s'", errors.Error())
	}

	if errors.First() != nil {
		t.Error("Expected First() to return nil for empty errors")
	}

	if len(errors.Issues()) != 0 {
		t.Error("Expected Issues() to return empty array")
	}
}
