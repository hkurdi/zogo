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
	isIP       bool
	isIPv4     bool
	isIPv6     bool
	isBase64   bool
	isHex      bool
	isCUID     bool
	isCUID2    bool
	isULID     bool
	isNanoid   bool
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

// IP validates IPv4 or IPv6 address
func (v *StringValidator) IP() *StringValidator {
	v.isIP = true
	return v
}

// IPv4 validates IPv4 address
func (v *StringValidator) IPv4() *StringValidator {
	v.isIPv4 = true
	return v
}

// IPv6 validates IPv6 address
func (v *StringValidator) IPv6() *StringValidator {
	v.isIPv6 = true
	return v
}

// Base64 validates base64 encoded string
func (v *StringValidator) Base64() *StringValidator {
	v.isBase64 = true
	return v
}

// Hex validates hexadecimal string
func (v *StringValidator) Hex() *StringValidator {
	v.isHex = true
	return v
}

// CUID validates CUID (Collision-resistant Unique Identifier)
func (v *StringValidator) CUID() *StringValidator {
	v.isCUID = true
	return v
}

// CUID2 validates CUID2 format
func (v *StringValidator) CUID2() *StringValidator {
	v.isCUID2 = true
	return v
}

// ULID validates ULID (Universally Unique Lexicographically Sortable Identifier)
func (v *StringValidator) ULID() *StringValidator {
	v.isULID = true
	return v
}

// Nanoid validates Nanoid format
func (v *StringValidator) Nanoid() *StringValidator {
	v.isNanoid = true
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

	// Check IP address
	if v.isIP && !isValidIP(str) {
		return FailureMessage("Invalid IP address")
	}

	// Check IPv4
	if v.isIPv4 && !isValidIPv4(str) {
		return FailureMessage("Invalid IPv4 address")
	}

	// Check IPv6
	if v.isIPv6 && !isValidIPv6(str) {
		return FailureMessage("Invalid IPv6 address")
	}

	// Check base64
	if v.isBase64 && !isValidBase64(str) {
		return FailureMessage("Invalid base64 string")
	}

	// Check hex
	if v.isHex && !isValidHex(str) {
		return FailureMessage("Invalid hexadecimal string")
	}

	// Check CUID
	if v.isCUID && !isValidCUID(str) {
		return FailureMessage("Invalid CUID format")
	}

	// Check CUID2
	if v.isCUID2 && !isValidCUID2(str) {
		return FailureMessage("Invalid CUID2 format")
	}

	// Check ULID
	if v.isULID && !isValidULID(str) {
		return FailureMessage("Invalid ULID format")
	}

	// Check Nanoid
	if v.isNanoid && !isValidNanoid(str) {
		return FailureMessage("Invalid Nanoid format")
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

// isValidIP checks if string is a valid IP address (v4 or v6)
func isValidIP(s string) bool {
	return isValidIPv4(s) || isValidIPv6(s)
}

// isValidIPv4 checks if string is a valid IPv4 address
func isValidIPv4(s string) bool {
	parts := strings.Split(s, ".")
	if len(parts) != 4 {
		return false
	}

	for _, part := range parts {
		if len(part) == 0 || len(part) > 3 {
			return false
		}

		// Check for leading zeros (except "0" itself)
		if len(part) > 1 && part[0] == '0' {
			return false
		}

		num := 0
		for _, ch := range part {
			if ch < '0' || ch > '9' {
				return false
			}
			num = num*10 + int(ch-'0')
		}

		if num > 255 {
			return false
		}
	}

	return true
}

// isValidIPv6 checks if string is a valid IPv6 address
func isValidIPv6(s string) bool {
	// Basic IPv6 validation
	// Supports standard format and :: compression
	if strings.Contains(s, ":::") {
		return false
	}

	// Split on ::
	parts := strings.Split(s, "::")
	if len(parts) > 2 {
		return false
	}

	var groups []string
	if len(parts) == 2 {
		// Has compression
		left := strings.Split(parts[0], ":")
		right := strings.Split(parts[1], ":")

		// Filter empty strings
		leftFiltered := make([]string, 0)
		for _, g := range left {
			if g != "" {
				leftFiltered = append(leftFiltered, g)
			}
		}

		rightFiltered := make([]string, 0)
		for _, g := range right {
			if g != "" {
				rightFiltered = append(rightFiltered, g)
			}
		}

		totalGroups := len(leftFiltered) + len(rightFiltered)
		if totalGroups > 7 {
			return false
		}

		groups = append(leftFiltered, rightFiltered...)
	} else {
		groups = strings.Split(s, ":")
		if len(groups) != 8 {
			return false
		}
	}

	// Validate each group
	for _, group := range groups {
		if len(group) == 0 || len(group) > 4 {
			return false
		}

		for _, ch := range group {
			if !((ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F')) {
				return false
			}
		}
	}

	return true
}

// isValidBase64 checks if string is valid base64
func isValidBase64(s string) bool {
	if len(s) == 0 {
		return false
	}

	// Base64 length must be multiple of 4
	if len(s)%4 != 0 {
		return false
	}

	for i, ch := range s {
		valid := (ch >= 'A' && ch <= 'Z') ||
			(ch >= 'a' && ch <= 'z') ||
			(ch >= '0' && ch <= '9') ||
			ch == '+' || ch == '/' ||
			(ch == '=' && i >= len(s)-2) // = only at end

		if !valid {
			return false
		}
	}

	return true
}

// isValidHex checks if string is valid hexadecimal
func isValidHex(s string) bool {
	if len(s) == 0 {
		return false
	}

	for _, ch := range s {
		if !((ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F')) {
			return false
		}
	}

	return true
}

// isValidCUID checks if string is a valid CUID
// Format: c + timestamp (base36) + counter (base36) + fingerprint + random (base36)
// Example: cjld2cjxh0000qzrmn831i7rn
func isValidCUID(s string) bool {
	if len(s) != 25 {
		return false
	}

	if s[0] != 'c' {
		return false
	}

	// Rest should be base36 (0-9, a-z)
	for i := 1; i < len(s); i++ {
		ch := s[i]
		if !((ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'z')) {
			return false
		}
	}

	return true
}

// isValidCUID2 checks if string is a valid CUID2
// CUID2 is variable length (24-32 chars) and starts with a letter
func isValidCUID2(s string) bool {
	length := len(s)
	if length < 24 || length > 32 {
		return false
	}

	// Must start with a letter
	if s[0] < 'a' || s[0] > 'z' {
		return false
	}

	// Rest should be alphanumeric lowercase
	for i := 1; i < len(s); i++ {
		ch := s[i]
		if !((ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'z')) {
			return false
		}
	}

	return true
}

// isValidULID checks if string is a valid ULID
// Format: 26 characters, base32 encoded (0-9, A-Z excluding I, L, O, U)
// Example: 01ARZ3NDEKTSV4RRFFQ69G5FAV
func isValidULID(s string) bool {
	if len(s) != 26 {
		return false
	}

	// ULID uses Crockford's base32: 0-9, A-Z excluding I, L, O, U
	for _, ch := range s {
		if !((ch >= '0' && ch <= '9') || (ch >= 'A' && ch <= 'H') ||
			(ch >= 'J' && ch <= 'K') || (ch >= 'M' && ch <= 'N') ||
			(ch >= 'P' && ch <= 'T') || (ch >= 'V' && ch <= 'Z')) {
			return false
		}
	}

	return true
}

// isValidNanoid checks if string is a valid Nanoid
// Default Nanoid is 21 characters, URL-safe alphabet
func isValidNanoid(s string) bool {
	// Nanoid can be various lengths, but default is 21
	// We'll accept 10-64 as reasonable range
	length := len(s)
	if length < 10 || length > 64 {
		return false
	}

	// Nanoid uses URL-safe alphabet: A-Za-z0-9_-
	for _, ch := range s {
		if !((ch >= '0' && ch <= '9') || (ch >= 'A' && ch <= 'Z') ||
			(ch >= 'a' && ch <= 'z') || ch == '_' || ch == '-') {
			return false
		}
	}

	return true
}
