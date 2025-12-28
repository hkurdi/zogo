package zogo

// ObjectValidator validates object/map values with nested schemas
type ObjectValidator struct {
	schema        Schema
	unknownFields string // "strict", "passthrough", or "strip"

	// Modifiers
	isRequired bool
	isOptional bool
	isNullable bool
}

// Object creates a new object validator with the given schema
func Object(schema Schema) *ObjectValidator {
	return &ObjectValidator{
		schema:        schema,
		unknownFields: "strip", // default: remove unknown fields
	}
}

// Strict makes the validator error on unknown fields
func (v *ObjectValidator) Strict() *ObjectValidator {
	v.unknownFields = "strict"
	return v
}

// Passthrough keeps unknown fields in the result
func (v *ObjectValidator) Passthrough() *ObjectValidator {
	v.unknownFields = "passthrough"
	return v
}

// Strip removes unknown fields from the result (default)
func (v *ObjectValidator) Strip() *ObjectValidator {
	v.unknownFields = "strip"
	return v
}

// Required marks the field as required
func (v *ObjectValidator) Required() *ObjectValidator {
	v.isRequired = true
	v.isOptional = false
	return v
}

// Optional allows nil values
func (v *ObjectValidator) Optional() *ObjectValidator {
	v.isOptional = true
	v.isRequired = false
	return v
}

// Nullable allows null values
func (v *ObjectValidator) Nullable() *ObjectValidator {
	v.isNullable = true
	return v
}

// Parse validates the input value
func (v *ObjectValidator) Parse(value any) ParseResult {
	// Handle nil values based on modifiers
	if value == nil {
		// If optional, nil is OK
		if v.isOptional {
			return Success(nil)
		}

		// If nullable, nil is OK
		if v.isNullable {
			return Success(nil)
		}

		// Otherwise, nil is not allowed
		return FailureMessage("Expected object, received null")
	}

	// Check if value is a map
	objMap, ok := value.(map[string]interface{})
	if !ok {
		return FailureMessage("Expected object, received " + typeof(value))
	}

	// Result object to build
	result := make(map[string]interface{})

	// Track all errors
	var errors ValidationErrors

	// Validate each field in the schema
	for fieldName, fieldValidator := range v.schema {
		fieldValue, exists := objMap[fieldName]

		// If field doesn't exist, pass nil to validator
		// The field validator will decide if that's OK based on its Optional/Required status
		if !exists {
			fieldValue = nil
		}

		// Validate the field
		fieldResult := fieldValidator.Parse(fieldValue)

		if !fieldResult.Ok {
			// Add field path to errors
			for _, err := range fieldResult.Errors {
				errors = append(errors, ValidationError{
					Path:    fieldName + prependPath(err.Path),
					Message: err.Message,
					Value:   err.Value,
				})
			}
		} else {
			// Only add to result if value is not nil
			// This prevents nil optional fields from appearing in output
			if fieldResult.Value != nil {
				result[fieldName] = fieldResult.Value
			}
		}
	}

	// Handle unknown fields (fields in objMap but not in schema)
	for fieldName, fieldValue := range objMap {
		// Check if field is in schema
		if _, inSchema := v.schema[fieldName]; !inSchema {
			switch v.unknownFields {
			case "strict":
				errors = append(errors, ValidationError{
					Path:    fieldName,
					Message: "Unknown field",
					Value:   fieldValue,
				})
			case "passthrough":
				result[fieldName] = fieldValue
			case "strip":
				// Do nothing - field is stripped
			}
		}
	}

	// Return errors if any
	if len(errors) > 0 {
		return Failure(errors...)
	}

	return Success(result)
}

// Helper function to prepend path separator
func prependPath(path string) string {
	if path == "" {
		return ""
	}
	return "." + path
}
