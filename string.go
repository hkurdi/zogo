package zogo

import (
	"fmt"
	"regexp"
	"strings"
)

type StringValidator struct {
	// Validation rules
	minLen   *int
	maxLen   *int
	exactLen *int
	pattern  *regexp.Regexp

	// Format validators
	isEmail    bool
	isURL      bool
	isUUID     bool
	startsWith *string
	endsWith   *string
	contains   *string

	// Transformations
	shouldTrim      bool
	shouldLowercase bool
	shouldUppercase bool

	// Modifiers
	isRequired bool
	isOptional bool
	isNullable bool
	defaultVal *string

	// Custom validators
	refinements []Refinement
}

type Refinement struct {
	Check   func(string) bool
	Message string
}

// String creates a new string validator
func String() *StringValidator {
	return &StringValidator{}
}

// Min sets the minimum string length
func (v *StringValidator) Min(length int) *StringValidator {
	v.minLen = &length
	return v
}

// Max sets the maximum string length
func (v *StringValidator) Max(length int) *StringValidator {
	v.maxLen = &length
	return v
}

// Length sets the exact string length required
func (v *StringValidator) Length(length int) *StringValidator {
	v.exactLen = &length
	return v
}

// Email validates email format
func (v *StringValidator) Email() *StringValidator {
	v.isEmail = true
	return v
}

// URL validates URL format
func (v *StringValidator) URL() *StringValidator {
	v.isURL = true
	return v
}

// UUID validates UUID format
func (v *StringValidator) UUID() *StringValidator {
	v.isUUID = true
	return v
}

// Regex validates against a regular expression pattern
func (v *StringValidator) Regex(pattern string) *StringValidator {
	v.pattern = regexp.MustCompile(pattern)
	return v
}

// StartsWith checks if string starts with the given prefix
func (v *StringValidator) StartsWith(prefix string) *StringValidator {
	v.startsWith = &prefix
	return v
}

// EndsWith checks if string ends with the given suffix
func (v *StringValidator) EndsWith(suffix string) *StringValidator {
	v.endsWith = &suffix
	return v
}

// Contains checks if string contains the given substring
func (v *StringValidator) Contains(substring string) *StringValidator {
	v.contains = &substring
	return v
}

// Trim removes leading and trailing whitespace
func (v *StringValidator) Trim() *StringValidator {
	v.shouldTrim = true
	return v
}

// ToLowerCase converts string to lowercase
func (v *StringValidator) ToLowerCase() *StringValidator {
	v.shouldLowercase = true
	return v
}

// ToUpperCase converts string to uppercase
func (v *StringValidator) ToUpperCase() *StringValidator {
	v.shouldUppercase = true
	return v
}

// Required marks the field as required (this is the default behavior)
func (v *StringValidator) Required() *StringValidator {
	v.isRequired = true
	v.isOptional = false
	return v
}

// Optional allows nil/undefined values
func (v *StringValidator) Optional() *StringValidator {
	v.isOptional = true
	v.isRequired = false
	return v
}

// Nullable allows null values
func (v *StringValidator) Nullable() *StringValidator {
	v.isNullable = true
	return v
}

// Default sets a default value if input is nil or empty string
func (v *StringValidator) Default(val string) *StringValidator {
	v.defaultVal = &val
	return v
}

// Refine adds custom validation logic
func (v *StringValidator) Refine(check func(string) bool, message string) *StringValidator {
	v.refinements = append(v.refinements, Refinement{
		Check:   check,
		Message: message,
	})
	return v
}

// Parse validates the input value
func (v *StringValidator) Parse(value any) ParseResult {
	// Check if value is nil
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
		return FailureMessage("Expected string, received null")
	}

	// Check if value is a string
	str, ok := value.(string)
	if !ok {
		return FailureMessage("Expected string, received " + typeof(value))
	}

	// Apply transformations first
	if v.shouldTrim {
		str = strings.TrimSpace(str)
	}

	if v.shouldLowercase {
		str = strings.ToLower(str)
	}

	if v.shouldUppercase {
		str = strings.ToUpper(str)
	}

	// Check exact length if specified
	if v.exactLen != nil && len(str) != *v.exactLen {
		return FailureMessage(fmt.Sprintf("String must be exactly %d characters", *v.exactLen))
	}

	// Check minimum length
	if v.minLen != nil && len(str) < *v.minLen {
		return FailureMessage(fmt.Sprintf("String must be at least %d characters", *v.minLen))
	}

	// Check maximum length
	if v.maxLen != nil && len(str) > *v.maxLen {
		return FailureMessage(fmt.Sprintf("String must be at most %d characters", *v.maxLen))
	}

	// Check email format
	if v.isEmail && !isValidEmail(str) {
		return FailureMessage("Invalid email format")
	}

	// Check URL format
	if v.isURL && !isValidURL(str) {
		return FailureMessage("Invalid URL format")
	}

	// Check UUID format
	if v.isUUID && !isValidUUID(str) {
		return FailureMessage("Invalid UUID format")
	}

	// Check regex pattern
	if v.pattern != nil && !v.pattern.MatchString(str) {
		return FailureMessage("String does not match required pattern")
	}

	// Check startsWith
	if v.startsWith != nil && !strings.HasPrefix(str, *v.startsWith) {
		return FailureMessage(fmt.Sprintf("String must start with '%s'", *v.startsWith))
	}

	// Check endsWith
	if v.endsWith != nil && !strings.HasSuffix(str, *v.endsWith) {
		return FailureMessage(fmt.Sprintf("String must end with '%s'", *v.endsWith))
	}

	// Check contains
	if v.contains != nil && !strings.Contains(str, *v.contains) {
		return FailureMessage(fmt.Sprintf("String must contain '%s'", *v.contains))
	}

	// Run custom refinements
	for _, refinement := range v.refinements {
		if !refinement.Check(str) {
			return FailureMessage(refinement.Message)
		}
	}

	return Success(str)
}

// Helper function to get type name
func typeof(value any) string {
	if value == nil {
		return "null"
	}
	switch value.(type) {
	case string:
		return "string"
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return "number"
	case float32, float64:
		return "number"
	case bool:
		return "boolean"
	default:
		return "unknown"
	}
}

// isValidEmail checks if string is a valid email
func isValidEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(email)
}

// isValidURL checks if string is a valid URL
func isValidURL(str string) bool {
	pattern := `^https?://[a-zA-Z0-9\-._~:/?#[\]@!$&'()*+,;=%]+$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(str)
}

// isValidUUID checks if string is a valid UUID
func isValidUUID(str string) bool {
	pattern := `^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(strings.ToLower(str))
}
