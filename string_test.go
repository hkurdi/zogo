package zogo

import "testing"

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
