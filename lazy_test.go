package zogo

import (
	"testing"
)

// Test basic lazy validator
func TestLazyBasic(t *testing.T) {
	schema := Lazy(func() Validator {
		return String()
	})

	result := schema.Parse("hello")
	if !result.Ok {
		t.Error("Expected lazy string to pass")
	}
	if result.Value != "hello" {
		t.Errorf("Expected 'hello', got %v", result.Value)
	}
}

// Test lazy with number
func TestLazyNumber(t *testing.T) {
	schema := Lazy(func() Validator {
		return Number().Min(10)
	})

	result := schema.Parse(15)
	if !result.Ok {
		t.Error("Expected valid number to pass")
	}

	result = schema.Parse(5)
	if result.Ok {
		t.Error("Expected number < 10 to fail")
	}
}

// Test self-referential schema - simple tree
func TestLazySimpleTree(t *testing.T) {
	var treeSchema Validator
	treeSchema = Lazy(func() Validator {
		return Object(Schema{
			"value":    Number(),
			"children": Array(treeSchema).Optional(),
		})
	})

	// Single node (no children)
	data := map[string]interface{}{
		"value": 1,
	}
	result := treeSchema.Parse(data)
	if !result.Ok {
		t.Errorf("Expected single node to pass. Errors: %v", result.Errors)
	}

	// Node with children
	data = map[string]interface{}{
		"value": 1,
		"children": []interface{}{
			map[string]interface{}{
				"value": 2,
			},
			map[string]interface{}{
				"value": 3,
			},
		},
	}
	result = treeSchema.Parse(data)
	if !result.Ok {
		t.Errorf("Expected tree with children to pass. Errors: %v", result.Errors)
	}
}

// Test deeply nested recursion
func TestLazyDeeplyNested(t *testing.T) {
	var treeSchema Validator
	treeSchema = Lazy(func() Validator {
		return Object(Schema{
			"value": Number(),
			"child": treeSchema.(*LazyValidator).Optional(),
		})
	})

	// 5 levels deep
	data := map[string]interface{}{
		"value": 1,
		"child": map[string]interface{}{
			"value": 2,
			"child": map[string]interface{}{
				"value": 3,
				"child": map[string]interface{}{
					"value": 4,
					"child": map[string]interface{}{
						"value": 5,
					},
				},
			},
		},
	}

	result := treeSchema.Parse(data)
	if !result.Ok {
		t.Errorf("Expected deeply nested tree to pass. Errors: %v", result.Errors)
	}
}

// Test comment thread (realistic use case)
func TestLazyCommentThread(t *testing.T) {
	var commentSchema Validator
	commentSchema = Lazy(func() Validator {
		return Object(Schema{
			"id":      String(),
			"text":    String(),
			"author":  String(),
			"replies": Array(commentSchema).Optional(),
		})
	})

	data := map[string]interface{}{
		"id":     "1",
		"text":   "Great article!",
		"author": "John",
		"replies": []interface{}{
			map[string]interface{}{
				"id":     "2",
				"text":   "Thanks!",
				"author": "Author",
			},
			map[string]interface{}{
				"id":     "3",
				"text":   "I agree!",
				"author": "Jane",
				"replies": []interface{}{
					map[string]interface{}{
						"id":     "4",
						"text":   "Me too!",
						"author": "Bob",
					},
				},
			},
		},
	}

	result := commentSchema.Parse(data)
	if !result.Ok {
		t.Errorf("Expected comment thread to pass. Errors: %v", result.Errors)
	}
}

// Test file system structure
func TestLazyFileSystem(t *testing.T) {
	var nodeSchema Validator
	nodeSchema = Lazy(func() Validator {
		return Object(Schema{
			"name":     String(),
			"type":     Enum([]interface{}{"file", "folder"}),
			"children": Array(nodeSchema).Optional(),
		})
	})

	data := map[string]interface{}{
		"name": "root",
		"type": "folder",
		"children": []interface{}{
			map[string]interface{}{
				"name": "documents",
				"type": "folder",
				"children": []interface{}{
					map[string]interface{}{
						"name": "report.pdf",
						"type": "file",
					},
				},
			},
			map[string]interface{}{
				"name": "readme.txt",
				"type": "file",
			},
		},
	}

	result := nodeSchema.Parse(data)
	if !result.Ok {
		t.Errorf("Expected file system to pass. Errors: %v", result.Errors)
	}
}

// Test organization chart
func TestLazyOrgChart(t *testing.T) {
	var employeeSchema Validator
	employeeSchema = Lazy(func() Validator {
		return Object(Schema{
			"name":          String(),
			"title":         String(),
			"directReports": Array(employeeSchema).Optional(),
		})
	})

	data := map[string]interface{}{
		"name":  "CEO",
		"title": "Chief Executive Officer",
		"directReports": []interface{}{
			map[string]interface{}{
				"name":  "CTO",
				"title": "Chief Technology Officer",
				"directReports": []interface{}{
					map[string]interface{}{
						"name":  "John",
						"title": "Senior Developer",
					},
					map[string]interface{}{
						"name":  "Jane",
						"title": "DevOps Engineer",
					},
				},
			},
			map[string]interface{}{
				"name":  "CFO",
				"title": "Chief Financial Officer",
			},
		},
	}

	result := employeeSchema.Parse(data)
	if !result.Ok {
		t.Errorf("Expected org chart to pass. Errors: %v", result.Errors)
	}
}

// Test category tree
func TestLazyCategoryTree(t *testing.T) {
	var categorySchema Validator
	categorySchema = Lazy(func() Validator {
		return Object(Schema{
			"name":          String(),
			"description":   String().Optional(),
			"subcategories": Array(categorySchema).Optional(),
		})
	})

	data := map[string]interface{}{
		"name":        "Electronics",
		"description": "Electronic devices",
		"subcategories": []interface{}{
			map[string]interface{}{
				"name": "Computers",
				"subcategories": []interface{}{
					map[string]interface{}{
						"name": "Laptops",
					},
					map[string]interface{}{
						"name": "Desktops",
					},
				},
			},
			map[string]interface{}{
				"name": "Phones",
			},
		},
	}

	result := categorySchema.Parse(data)
	if !result.Ok {
		t.Errorf("Expected category tree to pass. Errors: %v", result.Errors)
	}
}

// Test validation error in nested structure
func TestLazyNestedError(t *testing.T) {
	var treeSchema Validator
	treeSchema = Lazy(func() Validator {
		return Object(Schema{
			"value":    Number().Min(0),
			"children": Array(treeSchema).Optional(),
		})
	})

	// Invalid value deep in tree
	data := map[string]interface{}{
		"value": 1,
		"children": []interface{}{
			map[string]interface{}{
				"value": 2,
				"children": []interface{}{
					map[string]interface{}{
						"value": -5, // invalid!
					},
				},
			},
		},
	}

	result := treeSchema.Parse(data)
	if result.Ok {
		t.Error("Expected invalid nested value to fail")
	}

	// Check error path
	if len(result.Errors) == 0 {
		t.Error("Expected errors")
	}
}

// Test nil value
func TestLazyNil(t *testing.T) {
	schema := Lazy(func() Validator {
		return String()
	})

	result := schema.Parse(nil)
	if result.Ok {
		t.Error("Expected nil to fail")
	}
}

// Test Optional
func TestLazyOptional(t *testing.T) {
	schema := Lazy(func() Validator {
		return String()
	}).Optional()

	// nil should pass
	result := schema.Parse(nil)
	if !result.Ok {
		t.Error("Expected nil to pass with Optional()")
	}

	// Valid string should pass
	result = schema.Parse("hello")
	if !result.Ok {
		t.Error("Expected valid string to pass with Optional()")
	}
}

// Test Nullable
func TestLazyNullable(t *testing.T) {
	schema := Lazy(func() Validator {
		return String()
	}).Nullable()

	// nil should pass
	result := schema.Parse(nil)
	if !result.Ok {
		t.Error("Expected nil to pass with Nullable()")
	}
}

// Test Required
func TestLazyRequired(t *testing.T) {
	schema := Lazy(func() Validator {
		return String()
	}).Required()

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

// Test lazy in array
func TestLazyInArray(t *testing.T) {
	var nodeSchema Validator
	nodeSchema = Lazy(func() Validator {
		return Object(Schema{
			"value": Number(),
		})
	})

	schema := Array(nodeSchema)

	data := []interface{}{
		map[string]interface{}{"value": 1},
		map[string]interface{}{"value": 2},
		map[string]interface{}{"value": 3},
	}

	result := schema.Parse(data)
	if !result.Ok {
		t.Error("Expected array of lazy validators to pass")
	}
}

// Test lazy in object
func TestLazyInObject(t *testing.T) {
	var treeSchema Validator
	treeSchema = Lazy(func() Validator {
		return Object(Schema{
			"value": Number(),
		})
	})

	schema := Object(Schema{
		"name": String(),
		"tree": treeSchema,
	})

	data := map[string]interface{}{
		"name": "test",
		"tree": map[string]interface{}{
			"value": 42,
		},
	}

	result := schema.Parse(data)
	if !result.Ok {
		t.Error("Expected object with lazy field to pass")
	}
}

// Test multiple recursive references
func TestLazyMultipleReferences(t *testing.T) {
	var nodeSchema Validator
	nodeSchema = Lazy(func() Validator {
		return Object(Schema{
			"value": Number(),
			"left":  nodeSchema.(*LazyValidator).Optional(),
			"right": nodeSchema.(*LazyValidator).Optional(),
		})
	})

	// Binary tree
	data := map[string]interface{}{
		"value": 1,
		"left": map[string]interface{}{
			"value": 2,
		},
		"right": map[string]interface{}{
			"value": 3,
			"left": map[string]interface{}{
				"value": 4,
			},
		},
	}

	result := nodeSchema.Parse(data)
	if !result.Ok {
		t.Errorf("Expected binary tree to pass. Errors: %v", result.Errors)
	}
}

// Test menu structure (realistic nested menu)
func TestLazyMenuStructure(t *testing.T) {
	var menuItemSchema Validator
	menuItemSchema = Lazy(func() Validator {
		return Object(Schema{
			"label":    String(),
			"url":      String().Optional(),
			"children": Array(menuItemSchema).Optional(),
		})
	})

	data := map[string]interface{}{
		"label": "Main Menu",
		"children": []interface{}{
			map[string]interface{}{
				"label": "Products",
				"children": []interface{}{
					map[string]interface{}{
						"label": "Electronics",
						"url":   "/products/electronics",
					},
					map[string]interface{}{
						"label": "Clothing",
						"url":   "/products/clothing",
					},
				},
			},
			map[string]interface{}{
				"label": "About",
				"url":   "/about",
			},
		},
	}

	result := menuItemSchema.Parse(data)
	if !result.Ok {
		t.Errorf("Expected menu structure to pass. Errors: %v", result.Errors)
	}
}

// Test empty children array
func TestLazyEmptyChildren(t *testing.T) {
	var treeSchema Validator
	treeSchema = Lazy(func() Validator {
		return Object(Schema{
			"value":    Number(),
			"children": Array(treeSchema).Optional(),
		})
	})

	data := map[string]interface{}{
		"value":    1,
		"children": []interface{}{}, // empty array
	}

	result := treeSchema.Parse(data)
	if !result.Ok {
		t.Error("Expected node with empty children to pass")
	}
}

// Test lazy with union (discriminated recursive union)
func TestLazyWithUnion(t *testing.T) {
	var exprSchema Validator
	exprSchema = Lazy(func() Validator {
		return Union(
			Object(Schema{
				"type":  Literal("number"),
				"value": Number(),
			}),
			Object(Schema{
				"type":  Literal("add"),
				"left":  exprSchema,
				"right": exprSchema,
			}),
		)
	})

	// Expression: 5 + 3
	data := map[string]interface{}{
		"type": "add",
		"left": map[string]interface{}{
			"type":  "number",
			"value": 5,
		},
		"right": map[string]interface{}{
			"type":  "number",
			"value": 3,
		},
	}

	result := exprSchema.Parse(data)
	if !result.Ok {
		t.Errorf("Expected expression tree to pass. Errors: %v", result.Errors)
	}
}

// Test JSON schema (schema that validates schemas)
func TestLazyJSONSchema(t *testing.T) {
	var schemaSchema Validator
	schemaSchema = Lazy(func() Validator {
		return Object(Schema{
			"type":       Enum([]interface{}{"string", "number", "object", "array"}),
			"properties": Record(String(), schemaSchema).Optional(),
			"items":      schemaSchema.(*LazyValidator).Optional(),
		})
	})

	// A schema for a user object
	data := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"name": map[string]interface{}{
				"type": "string",
			},
			"age": map[string]interface{}{
				"type": "number",
			},
			"tags": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "string",
				},
			},
		},
	}

	result := schemaSchema.Parse(data)
	if !result.Ok {
		t.Errorf("Expected JSON schema to pass. Errors: %v", result.Errors)
	}
}
