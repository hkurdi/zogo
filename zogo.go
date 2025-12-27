// Package zogo provides schema validation for Go, inspired by Zod.
//
// Basic usage:
//
//	schema := zogo.String().Email().Min(5)
//	result := schema.Parse("[email protected]")
//	if !result.Ok {
//	    fmt.Println(result.Errors)
//	}
//
// Object validation:
//
//	userSchema := zogo.Object(zogo.Schema{
//	    "email": zogo.String().Email().Required(),
//	    "age":   zogo.Number().Min(18).Required(),
//	})
package zogo

// Schema represents a map of field names to validators for object validation
type Schema map[string]Validator
