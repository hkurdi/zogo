package zogo

import (
	"fmt"
	"time"
)

// DateValidator validates date/time values
type DateValidator struct {
	// Validation rules
	minDate *time.Time
	maxDate *time.Time

	// Type checks
	isFuture bool
	isPast   bool

	// Modifiers
	isRequired bool
	isOptional bool
	isNullable bool
	defaultVal *time.Time

	// Custom validators
	refinements []DateRefinement
}

// DateRefinement holds custom validation logic for dates
type DateRefinement struct {
	Check   func(time.Time) bool
	Message string
}

// Date creates a new date validator
func Date() *DateValidator {
	return &DateValidator{}
}

// Min sets the minimum date
func (v *DateValidator) Min(date time.Time) *DateValidator {
	v.minDate = &date
	return v
}

// Max sets the maximum date
func (v *DateValidator) Max(date time.Time) *DateValidator {
	v.maxDate = &date
	return v
}

// Future requires the date to be in the future
func (v *DateValidator) Future() *DateValidator {
	v.isFuture = true
	return v
}

// Past requires the date to be in the past
func (v *DateValidator) Past() *DateValidator {
	v.isPast = true
	return v
}

// Required marks the field as required
func (v *DateValidator) Required() *DateValidator {
	v.isRequired = true
	v.isOptional = false
	return v
}

// Optional allows nil values
func (v *DateValidator) Optional() *DateValidator {
	v.isOptional = true
	v.isRequired = false
	return v
}

// Nullable allows null values
func (v *DateValidator) Nullable() *DateValidator {
	v.isNullable = true
	return v
}

// Default sets a default value if input is nil
func (v *DateValidator) Default(val time.Time) *DateValidator {
	v.defaultVal = &val
	return v
}

// Refine adds custom validation logic
func (v *DateValidator) Refine(check func(time.Time) bool, message string) *DateValidator {
	v.refinements = append(v.refinements, DateRefinement{
		Check:   check,
		Message: message,
	})
	return v
}

// Parse validates the input value
func (v *DateValidator) Parse(value any) ParseResult {
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
		return FailureMessage("Expected date, received null")
	}

	// Try to convert to time.Time
	var dateVal time.Time

	switch v := value.(type) {
	case time.Time:
		dateVal = v
	case string:
		// Try parsing string as date
		parsed, err := parseDate(v)
		if err != nil {
			return FailureMessage("Invalid date string: " + err.Error())
		}
		dateVal = parsed
	default:
		return FailureMessage("Expected date, received " + typeof(value))
	}

	// Get current time for future/past checks
	now := time.Now()

	// Check if future
	if v.isFuture && !dateVal.After(now) {
		return FailureMessage("Date must be in the future")
	}

	// Check if past
	if v.isPast && !dateVal.Before(now) {
		return FailureMessage("Date must be in the past")
	}

	// Check minimum date
	if v.minDate != nil && dateVal.Before(*v.minDate) {
		return FailureMessage(fmt.Sprintf("Date must be at or after %s", v.minDate.Format(time.RFC3339)))
	}

	// Check maximum date
	if v.maxDate != nil && dateVal.After(*v.maxDate) {
		return FailureMessage(fmt.Sprintf("Date must be at or before %s", v.maxDate.Format(time.RFC3339)))
	}

	// Run custom refinements
	for _, refinement := range v.refinements {
		if !refinement.Check(dateVal) {
			return FailureMessage(refinement.Message)
		}
	}

	return Success(dateVal)
}

// parseDate tries to parse a string as a date using multiple common formats
func parseDate(s string) (time.Time, error) {
	// List of common date formats to try
	formats := []string{
		time.RFC3339,          // "2006-01-02T15:04:05Z07:00"
		time.RFC3339Nano,      // "2006-01-02T15:04:05.999999999Z07:00"
		"2006-01-02",          // "YYYY-MM-DD"
		"2006-01-02 15:04:05", // "YYYY-MM-DD HH:MM:SS"
		"2006-01-02T15:04:05", // "YYYY-MM-DDTHH:MM:SS"
		time.RFC1123,          // "Mon, 02 Jan 2006 15:04:05 MST"
		time.RFC1123Z,         // "Mon, 02 Jan 2006 15:04:05 -0700"
		time.RFC822,           // "02 Jan 06 15:04 MST"
		time.RFC822Z,          // "02 Jan 06 15:04 -0700"
		time.RFC850,           // "Monday, 02-Jan-06 15:04:05 MST"
		"01/02/2006",          // "MM/DD/YYYY"
		"01/02/2006 15:04:05", // "MM/DD/YYYY HH:MM:SS"
		"02-01-2006",          // "DD-MM-YYYY"
		"02-01-2006 15:04:05", // "DD-MM-YYYY HH:MM:SS"
	}

	var lastErr error
	for _, format := range formats {
		parsed, err := time.Parse(format, s)
		if err == nil {
			return parsed, nil
		}
		lastErr = err
	}

	return time.Time{}, lastErr
}
