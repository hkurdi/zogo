package zogo

import (
	"testing"
	"time"
)

// Test basic date validation with time.Time
func TestDateBasic(t *testing.T) {
	schema := Date()

	now := time.Now()
	result := schema.Parse(now)

	if !result.Ok {
		t.Errorf("Expected valid date to pass. Errors: %v", result.Errors)
	}

	resultDate, ok := result.Value.(time.Time)
	if !ok {
		t.Error("Expected result to be time.Time")
	}

	if !resultDate.Equal(now) {
		t.Error("Expected date to be preserved")
	}
}

// Test date parsing from string - RFC3339
func TestDateParseRFC3339(t *testing.T) {
	schema := Date()

	result := schema.Parse("2024-01-15T10:30:00Z")
	if !result.Ok {
		t.Errorf("Expected RFC3339 string to parse. Errors: %v", result.Errors)
	}

	resultDate := result.Value.(time.Time)
	if resultDate.Year() != 2024 || resultDate.Month() != 1 || resultDate.Day() != 15 {
		t.Error("Expected parsed date to match")
	}
}

// Test date parsing from string - YYYY-MM-DD
func TestDateParseYYYYMMDD(t *testing.T) {
	schema := Date()

	result := schema.Parse("2024-01-15")
	if !result.Ok {
		t.Errorf("Expected YYYY-MM-DD string to parse. Errors: %v", result.Errors)
	}

	resultDate := result.Value.(time.Time)
	if resultDate.Year() != 2024 || resultDate.Month() != 1 || resultDate.Day() != 15 {
		t.Error("Expected parsed date to match")
	}
}

// Test date parsing from string - MM/DD/YYYY
func TestDateParseMMDDYYYY(t *testing.T) {
	schema := Date()

	result := schema.Parse("01/15/2024")
	if !result.Ok {
		t.Errorf("Expected MM/DD/YYYY string to parse. Errors: %v", result.Errors)
	}

	resultDate := result.Value.(time.Time)
	if resultDate.Year() != 2024 || resultDate.Month() != 1 || resultDate.Day() != 15 {
		t.Error("Expected parsed date to match")
	}
}

// Test invalid date string
func TestDateInvalidString(t *testing.T) {
	schema := Date()

	result := schema.Parse("not a date")
	if result.Ok {
		t.Error("Expected invalid date string to fail")
	}
}

// Test invalid type
func TestDateInvalidType(t *testing.T) {
	schema := Date()

	result := schema.Parse(123)
	if result.Ok {
		t.Error("Expected number to fail date validation")
	}

	result = schema.Parse(true)
	if result.Ok {
		t.Error("Expected boolean to fail date validation")
	}
}

// Test nil value
func TestDateNil(t *testing.T) {
	schema := Date()

	result := schema.Parse(nil)
	if result.Ok {
		t.Error("Expected nil to fail")
	}
}

// Test Min date
func TestDateMin(t *testing.T) {
	minDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	schema := Date().Min(minDate)

	// Should pass - after min
	result := schema.Parse(time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC))
	if !result.Ok {
		t.Error("Expected date after min to pass")
	}

	// Should pass - equal to min
	result = schema.Parse(minDate)
	if !result.Ok {
		t.Error("Expected date equal to min to pass")
	}

	// Should fail - before min
	result = schema.Parse(time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC))
	if result.Ok {
		t.Error("Expected date before min to fail")
	}
}

// Test Max date
func TestDateMax(t *testing.T) {
	maxDate := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
	schema := Date().Max(maxDate)

	// Should pass - before max
	result := schema.Parse(time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC))
	if !result.Ok {
		t.Error("Expected date before max to pass")
	}

	// Should pass - equal to max
	result = schema.Parse(maxDate)
	if !result.Ok {
		t.Error("Expected date equal to max to pass")
	}

	// Should fail - after max
	result = schema.Parse(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC))
	if result.Ok {
		t.Error("Expected date after max to fail")
	}
}

// Test Min and Max together
func TestDateMinMax(t *testing.T) {
	minDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	maxDate := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
	schema := Date().Min(minDate).Max(maxDate)

	// Should pass - within range
	result := schema.Parse(time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC))
	if !result.Ok {
		t.Error("Expected date within range to pass")
	}

	// Should fail - before min
	result = schema.Parse(time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC))
	if result.Ok {
		t.Error("Expected date before min to fail")
	}

	// Should fail - after max
	result = schema.Parse(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC))
	if result.Ok {
		t.Error("Expected date after max to fail")
	}
}

// Test Future
func TestDateFuture(t *testing.T) {
	schema := Date().Future()

	// Should pass - 1 year in future
	futureDate := time.Now().Add(365 * 24 * time.Hour)
	result := schema.Parse(futureDate)
	if !result.Ok {
		t.Error("Expected future date to pass Future()")
	}

	// Should fail - in the past
	pastDate := time.Now().Add(-365 * 24 * time.Hour)
	result = schema.Parse(pastDate)
	if result.Ok {
		t.Error("Expected past date to fail Future()")
	}
}

// Test Past
func TestDatePast(t *testing.T) {
	schema := Date().Past()

	// Should pass - 1 year in past
	pastDate := time.Now().Add(-365 * 24 * time.Hour)
	result := schema.Parse(pastDate)
	if !result.Ok {
		t.Error("Expected past date to pass Past()")
	}

	// Should fail - in the future
	futureDate := time.Now().Add(365 * 24 * time.Hour)
	result = schema.Parse(futureDate)
	if result.Ok {
		t.Error("Expected future date to fail Past()")
	}
}

// Test Optional
func TestDateOptional(t *testing.T) {
	schema := Date().Optional()

	// nil should pass
	result := schema.Parse(nil)
	if !result.Ok {
		t.Error("Expected nil to pass with Optional()")
	}

	// Valid date should still pass
	result = schema.Parse(time.Now())
	if !result.Ok {
		t.Error("Expected valid date to pass with Optional()")
	}
}

// Test Nullable
func TestDateNullable(t *testing.T) {
	schema := Date().Nullable()

	// nil should pass
	result := schema.Parse(nil)
	if !result.Ok {
		t.Error("Expected nil to pass with Nullable()")
	}
}

// Test Default
func TestDateDefault(t *testing.T) {
	defaultDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	schema := Date().Default(defaultDate)

	// nil should return default
	result := schema.Parse(nil)
	if !result.Ok {
		t.Error("Expected nil to pass with Default()")
	}

	resultDate := result.Value.(time.Time)
	if !resultDate.Equal(defaultDate) {
		t.Errorf("Expected default date, got %v", resultDate)
	}

	// Provided value should override default
	customDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	result = schema.Parse(customDate)
	if !result.Ok {
		t.Error("Expected valid date to pass")
	}

	resultDate = result.Value.(time.Time)
	if !resultDate.Equal(customDate) {
		t.Errorf("Expected custom date, got %v", resultDate)
	}
}

// Test Refine
func TestDateRefine(t *testing.T) {
	// Must be a Monday
	schema := Date().Refine(func(d time.Time) bool {
		return d.Weekday() == time.Monday
	}, "Date must be a Monday")

	// Should pass - Monday
	monday := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC) // Jan 1, 2024 is a Monday
	result := schema.Parse(monday)
	if !result.Ok {
		t.Error("Expected Monday to pass")
	}

	// Should fail - Tuesday
	tuesday := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)
	result = schema.Parse(tuesday)
	if result.Ok {
		t.Error("Expected Tuesday to fail")
	}

	if len(result.Errors) == 0 || result.Errors[0].Message != "Date must be a Monday" {
		t.Error("Expected custom error message")
	}
}

// Test chaining validators
func TestDateChained(t *testing.T) {
	minDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	maxDate := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
	schema := Date().Min(minDate).Max(maxDate).Refine(func(d time.Time) bool {
		return d.Month() == time.June
	}, "Date must be in June")

	// Should pass
	juneDate := time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)
	result := schema.Parse(juneDate)
	if !result.Ok {
		t.Error("Expected June date within range to pass")
	}

	// Should fail - wrong month
	mayDate := time.Date(2024, 5, 15, 0, 0, 0, 0, time.UTC)
	result = schema.Parse(mayDate)
	if result.Ok {
		t.Error("Expected non-June date to fail")
	}

	// Should fail - before min
	result = schema.Parse(time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC))
	if result.Ok {
		t.Error("Expected date before min to fail")
	}
}

// Test date in object
func TestDateInObject(t *testing.T) {
	schema := Object(Schema{
		"name":      String(),
		"birthdate": Date().Past(),
	})

	data := map[string]interface{}{
		"name":      "John",
		"birthdate": time.Date(1990, 5, 15, 0, 0, 0, 0, time.UTC),
	}

	result := schema.Parse(data)
	if !result.Ok {
		t.Errorf("Expected object with date to pass. Errors: %v", result.Errors)
	}
}

// Test date parsing in object from string
func TestDateInObjectString(t *testing.T) {
	schema := Object(Schema{
		"name":      String(),
		"birthdate": Date(),
	})

	data := map[string]interface{}{
		"name":      "John",
		"birthdate": "1990-05-15",
	}

	result := schema.Parse(data)
	if !result.Ok {
		t.Errorf("Expected object with date string to pass. Errors: %v", result.Errors)
	}
}
