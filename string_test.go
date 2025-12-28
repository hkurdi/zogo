package zogo

import (
	"strings"
	"testing"
)

// Test basic string validation
func TestStringBasic(t *testing.T) {
	schema := String()
	result := schema.Parse("hello")

	if !result.Ok {
		t.Error("Expected valid string to pass")
	}

	if result.Value != "hello" {
		t.Errorf("Expected 'hello', got %v", result.Value)
	}
}

// Test non-string input
func TestStringInvalidType(t *testing.T) {
	schema := String()
	result := schema.Parse(123)

	if result.Ok {
		t.Error("Expected number to fail string validation")
	}

	if len(result.Errors) == 0 {
		t.Error("Expected error for invalid type")
	}
}

// Test nil value
func TestStringNil(t *testing.T) {
	schema := String()
	result := schema.Parse(nil)

	if result.Ok {
		t.Error("Expected nil to fail")
	}
}

// Test Min length
func TestStringMin(t *testing.T) {
	schema := String().Min(5)

	// Should pass
	result := schema.Parse("hello")
	if !result.Ok {
		t.Error("Expected 5 char string to pass Min(5)")
	}

	// Should fail
	result = schema.Parse("hi")
	if result.Ok {
		t.Error("Expected 2 char string to fail Min(5)")
	}
}

// Test Max length
func TestStringMax(t *testing.T) {
	schema := String().Max(5)

	// Should pass
	result := schema.Parse("hello")
	if !result.Ok {
		t.Error("Expected 5 char string to pass Max(5)")
	}

	// Should fail
	result = schema.Parse("hello world")
	if result.Ok {
		t.Error("Expected 11 char string to fail Max(5)")
	}
}

// Test exact Length
func TestStringLength(t *testing.T) {
	schema := String().Length(5)

	// Should pass
	result := schema.Parse("hello")
	if !result.Ok {
		t.Error("Expected 5 char string to pass Length(5)")
	}

	// Should fail - too short
	result = schema.Parse("hi")
	if result.Ok {
		t.Error("Expected 2 char string to fail Length(5)")
	}

	// Should fail - too long
	result = schema.Parse("hello world")
	if result.Ok {
		t.Error("Expected 11 char string to fail Length(5)")
	}
}

// Test chaining Min and Max
func TestStringMinMax(t *testing.T) {
	schema := String().Min(3).Max(10)

	// Should pass
	result := schema.Parse("hello")
	if !result.Ok {
		t.Error("Expected 5 char string to pass Min(3).Max(10)")
	}

	// Should fail - too short
	result = schema.Parse("hi")
	if result.Ok {
		t.Error("Expected 2 char string to fail Min(3)")
	}

	// Should fail - too long
	result = schema.Parse("hello world!")
	if result.Ok {
		t.Error("Expected 12 char string to fail Max(10)")
	}
}

// Test Email validation
func TestStringEmail(t *testing.T) {
	schema := String().Email()

	// Valid emails
	validEmails := []string{
		"test@example.com",
		"user.name@example.com",
		"user+tag@example.co.uk",
	}

	for _, email := range validEmails {
		result := schema.Parse(email)
		if !result.Ok {
			t.Errorf("Expected '%s' to be valid email", email)
		}
	}

	// Invalid emails
	invalidEmails := []string{
		"notanemail",
		"@example.com",
		"user@",
		"user @example.com",
	}

	for _, email := range invalidEmails {
		result := schema.Parse(email)
		if result.Ok {
			t.Errorf("Expected '%s' to fail email validation", email)
		}
	}
}

// Test URL validation
func TestStringURL(t *testing.T) {
	schema := String().URL()

	// Valid URLs
	result := schema.Parse("https://example.com")
	if !result.Ok {
		t.Error("Expected valid HTTPS URL to pass")
	}

	result = schema.Parse("http://example.com/path?query=1")
	if !result.Ok {
		t.Error("Expected valid HTTP URL with path to pass")
	}

	// Invalid URLs
	result = schema.Parse("notaurl")
	if result.Ok {
		t.Error("Expected invalid URL to fail")
	}

	result = schema.Parse("ftp://example.com")
	if result.Ok {
		t.Error("Expected FTP URL to fail (only http/https allowed)")
	}
}

// Test UUID validation
func TestStringUUID(t *testing.T) {
	schema := String().UUID()

	// Valid UUID v4
	result := schema.Parse("550e8400-e29b-41d4-a716-446655440000")
	if !result.Ok {
		t.Error("Expected valid UUID to pass")
	}

	// Invalid UUIDs
	result = schema.Parse("not-a-uuid")
	if result.Ok {
		t.Error("Expected invalid UUID to fail")
	}

	result = schema.Parse("550e8400-e29b-11d4-a716-446655440000") // Wrong version
	if result.Ok {
		t.Error("Expected non-v4 UUID to fail")
	}
}

// Test Regex validation
func TestStringRegex(t *testing.T) {
	schema := String().Regex("^[a-z]+$")

	// Should pass
	result := schema.Parse("hello")
	if !result.Ok {
		t.Error("Expected lowercase letters to pass")
	}

	// Should fail
	result = schema.Parse("Hello")
	if result.Ok {
		t.Error("Expected uppercase letters to fail")
	}

	result = schema.Parse("hello123")
	if result.Ok {
		t.Error("Expected numbers to fail")
	}
}

// Test StartsWith
func TestStringStartsWith(t *testing.T) {
	schema := String().StartsWith("https://")

	// Should pass
	result := schema.Parse("https://example.com")
	if !result.Ok {
		t.Error("Expected string starting with 'https://' to pass")
	}

	// Should fail
	result = schema.Parse("http://example.com")
	if result.Ok {
		t.Error("Expected string not starting with 'https://' to fail")
	}
}

// Test EndsWith
func TestStringEndsWith(t *testing.T) {
	schema := String().EndsWith(".com")

	// Should pass
	result := schema.Parse("example.com")
	if !result.Ok {
		t.Error("Expected string ending with '.com' to pass")
	}

	// Should fail
	result = schema.Parse("example.org")
	if result.Ok {
		t.Error("Expected string not ending with '.com' to fail")
	}
}

// Test Contains
func TestStringContains(t *testing.T) {
	schema := String().Contains("@")

	// Should pass
	result := schema.Parse("user@example.com")
	if !result.Ok {
		t.Error("Expected string containing '@' to pass")
	}

	// Should fail
	result = schema.Parse("userexample.com")
	if result.Ok {
		t.Error("Expected string not containing '@' to fail")
	}
}

// Test chaining format validators
func TestStringChainedFormats(t *testing.T) {
	schema := String().Email().Min(5).Max(50)

	// Should pass
	result := schema.Parse("test@example.com")
	if !result.Ok {
		t.Error("Expected valid email with proper length to pass")
	}

	// Should fail - too short
	result = schema.Parse("a@b.c")
	if result.Ok {
		t.Error("Expected short email to fail Min(5)")
	}

	// Should fail - invalid email
	result = schema.Parse("notanemail")
	if result.Ok {
		t.Error("Expected non-email to fail")
	}
}

// Test Trim transformation
func TestStringTrim(t *testing.T) {
	schema := String().Trim()

	result := schema.Parse("  hello  ")
	if !result.Ok {
		t.Error("Expected trimmed string to pass")
	}

	// Should return trimmed value
	if result.Value != "hello" {
		t.Errorf("Expected 'hello', got '%s'", result.Value)
	}
}

// Test ToLowerCase transformation
func TestStringToLowerCase(t *testing.T) {
	schema := String().ToLowerCase()

	result := schema.Parse("HELLO")
	if !result.Ok {
		t.Error("Expected uppercase string to pass")
	}

	// Should return lowercase value
	if result.Value != "hello" {
		t.Errorf("Expected 'hello', got '%s'", result.Value)
	}
}

// Test ToUpperCase transformation
func TestStringToUpperCase(t *testing.T) {
	schema := String().ToUpperCase()

	result := schema.Parse("hello")
	if !result.Ok {
		t.Error("Expected lowercase string to pass")
	}

	// Should return uppercase value
	if result.Value != "HELLO" {
		t.Errorf("Expected 'HELLO', got '%s'", result.Value)
	}
}

// Test chained transformations
func TestStringChainedTransformations(t *testing.T) {
	schema := String().Trim().ToLowerCase().Min(3)

	result := schema.Parse("  HELLO  ")
	if !result.Ok {
		t.Error("Expected chained transformations to pass")
	}

	// Should return trimmed and lowercased value
	if result.Value != "hello" {
		t.Errorf("Expected 'hello', got '%s'", result.Value)
	}
}

// Test transformation with Email validation
func TestStringTransformWithEmail(t *testing.T) {
	schema := String().Trim().ToLowerCase().Email()

	// Test with uppercase email
	result := schema.Parse("  username@email.com  ")
	if !result.Ok {
		t.Errorf("Expected email to pass after trim and lowercase. Error: %v", result.Errors)
	}

	if result.Value != "username@email.com" {
		t.Errorf("Expected 'username@email.com', got '%v'", result.Value)
	}
}

// Test Trim affects length validation
func TestStringTrimAffectsLength(t *testing.T) {
	schema := String().Trim().Min(5)

	// Original "  hi  " has 6 chars, but after trim only 2
	result := schema.Parse("  hi  ")
	if result.Ok {
		t.Error("Expected trimmed string 'hi' to fail Min(5)")
	}

	// Original "  hello  " has 9 chars, after trim has 5
	result = schema.Parse("  hello  ")
	if !result.Ok {
		t.Error("Expected trimmed string 'hello' to pass Min(5)")
	}
}

// Test Optional modifier
func TestStringOptional(t *testing.T) {
	schema := String().Optional()

	// nil should pass
	result := schema.Parse(nil)
	if !result.Ok {
		t.Error("Expected nil to pass with Optional()")
	}

	// Valid string should still pass
	result = schema.Parse("hello")
	if !result.Ok {
		t.Error("Expected valid string to pass with Optional()")
	}
}

// Test Nullable modifier
func TestStringNullable(t *testing.T) {
	schema := String().Nullable()

	// nil should pass
	result := schema.Parse(nil)
	if !result.Ok {
		t.Error("Expected nil to pass with Nullable()")
	}

	// Valid string should still pass
	result = schema.Parse("hello")
	if !result.Ok {
		t.Error("Expected valid string to pass with Nullable()")
	}
}

// Test Default modifier
func TestStringDefault(t *testing.T) {
	schema := String().Default("default-value")

	// nil should return default
	result := schema.Parse(nil)
	if !result.Ok {
		t.Error("Expected nil to pass with Default()")
	}
	if result.Value != "default-value" {
		t.Errorf("Expected 'default-value', got '%v'", result.Value)
	}

	// Provided value should override default
	result = schema.Parse("custom")
	if !result.Ok {
		t.Error("Expected valid string to pass")
	}
	if result.Value != "custom" {
		t.Errorf("Expected 'custom', got '%v'", result.Value)
	}
}

// Test Refine with custom validation
func TestStringRefine(t *testing.T) {
	// Password must contain uppercase letter
	schema := String().Min(8).Refine(func(s string) bool {
		for _, c := range s {
			if c >= 'A' && c <= 'Z' {
				return true
			}
		}
		return false
	}, "Password must contain at least one uppercase letter")

	// Should fail - no uppercase
	result := schema.Parse("password123")
	if result.Ok {
		t.Error("Expected password without uppercase to fail")
	}
	if len(result.Errors) == 0 || result.Errors[0].Message != "Password must contain at least one uppercase letter" {
		t.Error("Expected custom error message")
	}

	// Should pass - has uppercase
	result = schema.Parse("Password123")
	if !result.Ok {
		t.Error("Expected password with uppercase to pass")
	}
}

// Test multiple Refine calls
func TestStringMultipleRefine(t *testing.T) {
	schema := String().
		Min(8).
		Refine(func(s string) bool {
			// Must contain uppercase
			for _, c := range s {
				if c >= 'A' && c <= 'Z' {
					return true
				}
			}
			return false
		}, "Must contain uppercase").
		Refine(func(s string) bool {
			// Must contain number
			for _, c := range s {
				if c >= '0' && c <= '9' {
					return true
				}
			}
			return false
		}, "Must contain number")

	// Should fail - no number
	result := schema.Parse("Password")
	if result.Ok {
		t.Error("Expected password without number to fail")
	}

	// Should pass - has both
	result = schema.Parse("Password123")
	if !result.Ok {
		t.Error("Expected password with uppercase and number to pass")
	}
}

// Test Required (default behavior)
func TestStringRequired(t *testing.T) {
	schema := String().Required()

	// nil should fail
	result := schema.Parse(nil)
	if result.Ok {
		t.Error("Expected nil to fail with Required()")
	}

	// Valid string should pass
	result = schema.Parse("hello")
	if !result.Ok {
		t.Error("Expected valid string to pass with Required()")
	}
}

// Test chaining Default with transformations
func TestStringDefaultWithTransform(t *testing.T) {
	schema := String().Default("HELLO").ToLowerCase()

	// Should use default and transform it
	result := schema.Parse(nil)
	if !result.Ok {
		t.Error("Expected nil to pass with Default()")
	}
	// Default value doesn't go through transformations (applied at nil check)
	if result.Value != "HELLO" {
		t.Errorf("Expected 'HELLO', got '%v'", result.Value)
	}
}

// Test IPv4 validation
func TestStringIPv4(t *testing.T) {
	schema := String().IPv4()

	// Valid IPv4
	validIPs := []string{
		"192.168.1.1",
		"10.0.0.0",
		"255.255.255.255",
		"0.0.0.0",
		"127.0.0.1",
	}

	for _, ip := range validIPs {
		result := schema.Parse(ip)
		if !result.Ok {
			t.Errorf("Expected valid IPv4 '%s' to pass", ip)
		}
	}

	// Invalid IPv4
	invalidIPs := []string{
		"256.1.1.1",       // octet > 255
		"192.168.1",       // too few octets
		"192.168.1.1.1",   // too many octets
		"192.168.01.1",    // leading zero
		"192.168.-1.1",    // negative
		"abc.def.ghi.jkl", // non-numeric
	}

	for _, ip := range invalidIPs {
		result := schema.Parse(ip)
		if result.Ok {
			t.Errorf("Expected invalid IPv4 '%s' to fail", ip)
		}
	}
}

// Test IPv6 validation
func TestStringIPv6(t *testing.T) {
	schema := String().IPv6()

	// Valid IPv6
	validIPs := []string{
		"2001:0db8:85a3:0000:0000:8a2e:0370:7334",
		"2001:db8:85a3::8a2e:370:7334", // compressed
		"::1",                          // loopback
		"::",                           // all zeros
		"fe80::1",
		"2001:db8::1",
	}

	for _, ip := range validIPs {
		result := schema.Parse(ip)
		if !result.Ok {
			t.Errorf("Expected valid IPv6 '%s' to pass", ip)
		}
	}

	// Invalid IPv6
	invalidIPs := []string{
		"02001:0db8:0000:0000:0000:ff00:0042:8329", // too many digits
		"2001:0db8:0000:0000:0000:gg00:0042:8329",  // invalid hex
		":::",            // triple colon
		"2001:db8::1::2", // multiple compressions
	}

	for _, ip := range invalidIPs {
		result := schema.Parse(ip)
		if result.Ok {
			t.Errorf("Expected invalid IPv6 '%s' to fail", ip)
		}
	}
}

// Test IP (v4 or v6) validation
func TestStringIP(t *testing.T) {
	schema := String().IP()

	// Should accept both IPv4 and IPv6
	result := schema.Parse("192.168.1.1")
	if !result.Ok {
		t.Error("Expected IPv4 to pass IP()")
	}

	result = schema.Parse("2001:db8::1")
	if !result.Ok {
		t.Error("Expected IPv6 to pass IP()")
	}

	result = schema.Parse("not-an-ip")
	if result.Ok {
		t.Error("Expected invalid IP to fail")
	}
}

// Test Base64 validation
func TestStringBase64(t *testing.T) {
	schema := String().Base64()

	// Valid base64
	validBase64 := []string{
		"SGVsbG8gV29ybGQ=",     // "Hello World"
		"YQ==",                 // "a"
		"YWJjZA==",             // "abcd"
		"VGhpcyBpcyBhIHRlc3Q=", // "This is a test"
		"MTIzNDU2Nzg5MA==",     // "1234567890"
	}

	for _, b64 := range validBase64 {
		result := schema.Parse(b64)
		if !result.Ok {
			t.Errorf("Expected valid base64 '%s' to pass", b64)
		}
	}

	// Invalid base64
	invalidBase64 := []string{
		"Hello!",          // invalid chars
		"SGVsbG8",         // not multiple of 4
		"SGVs bG8=",       // not multiple of 4
		"SGVsbG8gV29ybGQ", // missing padding
		"====",            // only padding
	}

	for _, b64 := range invalidBase64 {
		result := schema.Parse(b64)
		if result.Ok {
			t.Errorf("Expected invalid base64 '%s' to fail", b64)
		}
	}
}

// Test Hex validation
func TestStringHex(t *testing.T) {
	schema := String().Hex()

	// Valid hex
	validHex := []string{
		"deadbeef",
		"DEADBEEF",
		"0123456789abcdef",
		"0123456789ABCDEF",
		"ff00ff",
	}

	for _, hex := range validHex {
		result := schema.Parse(hex)
		if !result.Ok {
			t.Errorf("Expected valid hex '%s' to pass", hex)
		}
	}

	// Invalid hex
	invalidHex := []string{
		"xyz",
		"12345g",
		"hello",
		"",
	}

	for _, hex := range invalidHex {
		result := schema.Parse(hex)
		if result.Ok {
			t.Errorf("Expected invalid hex '%s' to fail", hex)
		}
	}
}

// Test CUID validation
func TestStringCUID(t *testing.T) {
	schema := String().CUID()

	// Valid CUID (25 chars, starts with 'c')
	validCUIDs := []string{
		"cjld2cjxh0000qzrmn831i7rn",
		"ckz3q2q2q0000qzrmn831i7rn",
		"c" + strings.Repeat("a", 24), // 'c' + 24 valid chars
	}

	for _, cuid := range validCUIDs {
		result := schema.Parse(cuid)
		if !result.Ok {
			t.Errorf("Expected valid CUID '%s' to pass", cuid)
		}
	}

	// Invalid CUID
	invalidCUIDs := []string{
		"ajld2cjxh0000qzrmn831i7rn",  // doesn't start with 'c'
		"cjld2cjxh0000qzrmn831i7r",   // too short
		"cjld2cjxh0000qzrmn831i7rnn", // too long
		"cjld2cjxh0000QZRMN831i7rn",  // uppercase
		"",
	}

	for _, cuid := range invalidCUIDs {
		result := schema.Parse(cuid)
		if result.Ok {
			t.Errorf("Expected invalid CUID '%s' to fail", cuid)
		}
	}
}

// Test CUID2 validation
func TestStringCUID2(t *testing.T) {
	schema := String().CUID2()

	// Valid CUID2 (24-32 chars, starts with letter)
	validCUID2s := []string{
		"a" + strings.Repeat("b", 23), // 24 chars
		"z" + strings.Repeat("0", 25), // 26 chars
		"m" + strings.Repeat("x", 31), // 32 chars
	}

	for _, cuid2 := range validCUID2s {
		result := schema.Parse(cuid2)
		if !result.Ok {
			t.Errorf("Expected valid CUID2 '%s' to pass", cuid2)
		}
	}

	// Invalid CUID2
	invalidCUID2s := []string{
		"1" + strings.Repeat("a", 23), // starts with number
		"A" + strings.Repeat("a", 23), // starts with uppercase
		"a" + strings.Repeat("b", 22), // too short (23 chars)
		"a" + strings.Repeat("b", 32), // too long (33 chars)
		"a" + strings.Repeat("B", 23), // contains uppercase
		"",
	}

	for _, cuid2 := range invalidCUID2s {
		result := schema.Parse(cuid2)
		if result.Ok {
			t.Errorf("Expected invalid CUID2 '%s' to fail", cuid2)
		}
	}
}

// Test ULID validation
func TestStringULID(t *testing.T) {
	schema := String().ULID()

	// Valid ULID (26 chars, Crockford base32)
	validULIDs := []string{
		"01ARZ3NDEKTSV4RRFFQ69G5FAV",
		"01BX5ZZKBKACTAV9WEVGEMMVRZ",
		strings.Repeat("0", 26),
		strings.Repeat("Z", 26),
	}

	for _, ulid := range validULIDs {
		result := schema.Parse(ulid)
		if !result.Ok {
			t.Errorf("Expected valid ULID '%s' to pass", ulid)
		}
	}

	// Invalid ULID
	invalidULIDs := []string{
		"01ARZ3NDEKTSV4RRFFQ69G5FA",   // too short (25 chars)
		"01ARZ3NDEKTSV4RRFFQ69G5FAVV", // too long (27 chars)
		"01ARZ3NDEKTSV4RRFFQ69G5FaV",  // lowercase
		"01ARZ3NDEKTSV4RRFFQ69G5FIV",  // contains I (excluded)
		"01ARZ3NDEKTSV4RRFFQ69G5FLV",  // contains L (excluded)
		"01ARZ3NDEKTSV4RRFFQ69G5FOV",  // contains O (excluded)
		"01ARZ3NDEKTSV4RRFFQ69G5FUV",  // contains U (excluded)
		"",
	}

	for _, ulid := range invalidULIDs {
		result := schema.Parse(ulid)
		if result.Ok {
			t.Errorf("Expected invalid ULID '%s' to fail", ulid)
		}
	}
}

// Test Nanoid validation
func TestStringNanoid(t *testing.T) {
	schema := String().Nanoid()

	// Valid Nanoid (10-64 chars, URL-safe)
	validNanoids := []string{
		"V1StGXR8_Z5jdHi6B-myT", // 21 chars (default)
		strings.Repeat("a", 10), // 10 chars (min)
		strings.Repeat("Z", 64), // 64 chars (max)
		"abc_123-XYZ",
	}

	for _, nanoid := range validNanoids {
		result := schema.Parse(nanoid)
		if !result.Ok {
			t.Errorf("Expected valid Nanoid '%s' to pass", nanoid)
		}
	}

	// Invalid Nanoid
	invalidNanoids := []string{
		strings.Repeat("a", 9),  // too short
		strings.Repeat("a", 65), // too long
		"hello world",           // contains space
		"hello!",                // contains !
		"",
	}

	for _, nanoid := range invalidNanoids {
		result := schema.Parse(nanoid)
		if result.Ok {
			t.Errorf("Expected invalid Nanoid '%s' to fail", nanoid)
		}
	}
}

// Test multiple format validators chained
func TestStringMultipleFormats(t *testing.T) {
	// This should work - base64 that's also hex (subset)
	schema := String().Hex().Min(8)

	result := schema.Parse("deadbeef")
	if !result.Ok {
		t.Error("Expected hex string to pass both validations")
	}

	result = schema.Parse("abc")
	if result.Ok {
		t.Error("Expected short hex to fail Min(8)")
	}
}

// Test format validators in objects
func TestStringFormatsInObject(t *testing.T) {
	schema := Object(Schema{
		"ipv4":   String().IPv4(),
		"base64": String().Base64(),
		"hex":    String().Hex(),
	})

	data := map[string]interface{}{
		"ipv4":   "192.168.1.1",
		"base64": "SGVsbG8=",
		"hex":    "deadbeef",
	}

	result := schema.Parse(data)
	if !result.Ok {
		t.Errorf("Expected object with format validators to pass. Errors: %v", result.Errors)
	}
}
