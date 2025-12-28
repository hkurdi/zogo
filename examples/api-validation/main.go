package main

import (
	"encoding/json"
	"fmt"

	"github.com/hkurdi/zogo"
)

func main() {
	fmt.Println("=== Zogo Examples: API Validation ===\n")

	// Example 1: User Registration
	userRegistrationExample()

	// Example 2: Blog Post with Comments (Recursive)
	blogPostExample()

	// Example 3: E-commerce Order (Discriminated Union)
	ecommerceOrderExample()

	// Example 4: Configuration Validation
	configValidationExample()

	// Example 5: API Response (Union)
	apiResponseExample()
}

// Example 1: User Registration
func userRegistrationExample() {
	fmt.Println("--- Example 1: User Registration ---")

	// Define schema
	userSchema := zogo.Object(zogo.Schema{
		"username": zogo.String().Min(3).Max(20),
		"email":    zogo.String().Email(),
		"password": zogo.String().Min(8).Refine(func(s string) bool {
			// Password must contain uppercase
			for _, c := range s {
				if c >= 'A' && c <= 'Z' {
					return true
				}
			}
			return false
		}, "Password must contain at least one uppercase letter"),
		"age":       zogo.Number().Min(18).Optional(),
		"ipAddress": zogo.String().IPv4().Optional(),
		"metadata":  zogo.Record(zogo.String(), zogo.Any()).Optional(),
	})

	// Valid user
	validUser := map[string]interface{}{
		"username":  "john_doe",
		"email":     "john@example.com",
		"password":  "SecurePass123",
		"age":       25,
		"ipAddress": "192.168.1.1",
		"metadata": map[string]interface{}{
			"theme":    "dark",
			"language": "en",
		},
	}

	result := userSchema.Parse(validUser)
	if result.Ok {
		fmt.Println("✓ Valid user registration")
	} else {
		fmt.Printf("✗ Validation failed: %v\n", result.Errors)
	}

	// Invalid user
	invalidUser := map[string]interface{}{
		"username":  "jo",              // too short
		"email":     "not-an-email",    // invalid email
		"password":  "weak",            // too short, no uppercase
		"age":       15,                // too young
		"ipAddress": "999.999.999.999", // invalid IP
	}

	result = userSchema.Parse(invalidUser)
	if !result.Ok {
		fmt.Println("✓ Invalid user correctly rejected")
		fmt.Println("  Errors:")
		for _, err := range result.Errors {
			fmt.Printf("    - %s: %s\n", err.Path, err.Message)
		}
	}

	fmt.Println()
}

// Example 2: Blog Post with Recursive Comments
func blogPostExample() {
	fmt.Println("--- Example 2: Blog Post with Recursive Comments ---")

	// Recursive comment schema
	var commentSchema zogo.Validator
	commentSchema = zogo.Lazy(func() zogo.Validator {
		return zogo.Object(zogo.Schema{
			"id":      zogo.String().CUID(),
			"author":  zogo.String(),
			"text":    zogo.String().Min(1),
			"replies": zogo.Array(commentSchema).Optional(),
		})
	})

	// Blog post schema
	postSchema := zogo.Object(zogo.Schema{
		"id":        zogo.String().ULID(),
		"title":     zogo.String().Min(5).Max(200),
		"content":   zogo.String().Min(10),
		"author":    zogo.String(),
		"tags":      zogo.Array(zogo.String()).Min(1).Max(5),
		"published": zogo.Date().Past().Optional(),
		"comments":  zogo.Array(commentSchema).Optional(),
	})

	// Valid blog post with nested comments
	post := map[string]interface{}{
		"id":        "01ARZ3NDEKTSV4RRFFQ69G5FAV",
		"title":     "Introduction to Zogo",
		"content":   "Zogo is a powerful validation library for Go...",
		"author":    "Jane Doe",
		"tags":      []interface{}{"go", "validation", "zogo"},
		"published": "2024-01-15T10:00:00Z",
		"comments": []interface{}{
			map[string]interface{}{
				"id":     "cjld2cjxh0000qzrmn831i7rn",
				"author": "Reader1",
				"text":   "Great post!",
				"replies": []interface{}{
					map[string]interface{}{
						"id":     "cjld2cjxh0001qzrmn831i7ro",
						"author": "Jane Doe",
						"text":   "Thanks!",
					},
				},
			},
		},
	}

	result := postSchema.Parse(post)
	if result.Ok {
		fmt.Println("✓ Valid blog post with nested comments")
	} else {
		fmt.Printf("✗ Validation failed: %v\n", result.Errors)
	}

	fmt.Println()
}

// Example 3: E-commerce Order (Discriminated Union)
func ecommerceOrderExample() {
	fmt.Println("--- Example 3: E-commerce Order (Discriminated Union) ---")

	// Different payment methods
	creditCardPayment := zogo.Object(zogo.Schema{
		"method":     zogo.Literal("credit_card"),
		"cardNumber": zogo.String().Regex(`^\d{16}$`),
		"cvv":        zogo.String().Regex(`^\d{3,4}$`),
	})

	paypalPayment := zogo.Object(zogo.Schema{
		"method": zogo.Literal("paypal"),
		"email":  zogo.String().Email(),
	})

	cryptoPayment := zogo.Object(zogo.Schema{
		"method":  zogo.Literal("crypto"),
		"wallet":  zogo.String().Hex(),
		"network": zogo.Enum([]interface{}{"ethereum", "bitcoin", "polygon"}),
	})

	// Union of payment methods
	paymentSchema := zogo.Union(creditCardPayment, paypalPayment, cryptoPayment)

	// Order schema
	orderSchema := zogo.Object(zogo.Schema{
		"orderId": zogo.String().Nanoid(),
		"items": zogo.Array(zogo.Object(zogo.Schema{
			"productId": zogo.String(),
			"quantity":  zogo.Number().Int().Positive(),
			"price":     zogo.Number().Min(0),
		})).Min(1),
		"total":   zogo.Number().Min(0),
		"payment": paymentSchema,
		"shipping": zogo.Object(zogo.Schema{
			"address": zogo.String(),
			"city":    zogo.String(),
			"zipCode": zogo.String().Regex(`^\d{5}$`),
		}),
	})

	// Valid order with credit card
	order1 := map[string]interface{}{
		"orderId": "V1StGXR8_Z5jdHi6B-myT",
		"items": []interface{}{
			map[string]interface{}{
				"productId": "PROD-001",
				"quantity":  2,
				"price":     29.99,
			},
		},
		"total": 59.98,
		"payment": map[string]interface{}{
			"method":     "credit_card",
			"cardNumber": "1234567812345678",
			"cvv":        "123",
		},
		"shipping": map[string]interface{}{
			"address": "123 Main St",
			"city":    "Boston",
			"zipCode": "02101",
		},
	}

	result := orderSchema.Parse(order1)
	if result.Ok {
		fmt.Println("✓ Valid order with credit card payment")
	}

	// Valid order with PayPal
	order2 := map[string]interface{}{
		"orderId": "A2BtGYR9_Z6jdHi7C-nzU",
		"items": []interface{}{
			map[string]interface{}{
				"productId": "PROD-002",
				"quantity":  1,
				"price":     99.99,
			},
		},
		"total": 99.99,
		"payment": map[string]interface{}{
			"method": "paypal",
			"email":  "user@paypal.com",
		},
		"shipping": map[string]interface{}{
			"address": "456 Oak Ave",
			"city":    "Seattle",
			"zipCode": "98101",
		},
	}

	result = orderSchema.Parse(order2)
	if result.Ok {
		fmt.Println("✓ Valid order with PayPal payment")
	}

	fmt.Println()
}

// Example 4: Configuration Validation
func configValidationExample() {
	fmt.Println("--- Example 4: Configuration Validation ---")

	configSchema := zogo.Object(zogo.Schema{
		"server": zogo.Object(zogo.Schema{
			"host": zogo.String().IP(),
			"port": zogo.Number().Int().Min(1).Max(65535),
		}),
		"database": zogo.Object(zogo.Schema{
			"host":     zogo.String(),
			"port":     zogo.Number().Int(),
			"name":     zogo.String(),
			"ssl":      zogo.Boolean().Default(true),
			"poolSize": zogo.Number().Int().Min(1).Max(100).Default(10),
		}),
		"features": zogo.Record(zogo.String(), zogo.Boolean()),
		"logLevel": zogo.Enum([]interface{}{"debug", "info", "warn", "error"}),
	})

	config := map[string]interface{}{
		"server": map[string]interface{}{
			"host": "127.0.0.1",
			"port": 8080,
		},
		"database": map[string]interface{}{
			"host": "localhost",
			"port": 5432,
			"name": "myapp",
			// ssl and poolSize will use defaults
		},
		"features": map[string]interface{}{
			"authentication": true,
			"caching":        true,
			"analytics":      false,
		},
		"logLevel": "info",
	}

	result := configSchema.Parse(config)
	if result.Ok {
		fmt.Println("✓ Valid configuration (with defaults applied)")
		resultJSON, _ := json.MarshalIndent(result.Value, "  ", "  ")
		fmt.Printf("  Result: %s\n", resultJSON)
	}

	fmt.Println()
}

// Example 5: API Response Validation (Success/Error Union)
func apiResponseExample() {
	fmt.Println("--- Example 5: API Response Validation ---")

	// Success response
	successResponse := zogo.Object(zogo.Schema{
		"status": zogo.Literal("success"),
		"data": zogo.Object(zogo.Schema{
			"users": zogo.Array(zogo.Object(zogo.Schema{
				"id":    zogo.String(),
				"name":  zogo.String(),
				"email": zogo.String().Email(),
			})),
			"total": zogo.Number().Int(),
		}),
	})

	// Error response
	errorResponse := zogo.Object(zogo.Schema{
		"status":  zogo.Literal("error"),
		"message": zogo.String(),
		"code":    zogo.Number().Int(),
	})

	// Union schema
	responseSchema := zogo.Union(successResponse, errorResponse)

	// Test success response
	success := map[string]interface{}{
		"status": "success",
		"data": map[string]interface{}{
			"users": []interface{}{
				map[string]interface{}{
					"id":    "1",
					"name":  "Alice",
					"email": "alice@example.com",
				},
			},
			"total": 1,
		},
	}

	result := responseSchema.Parse(success)
	if result.Ok {
		fmt.Println("✓ Valid success response")
	}

	// Test error response
	errorResp := map[string]interface{}{
		"status":  "error",
		"message": "User not found",
		"code":    404,
	}

	result = responseSchema.Parse(errorResp)
	if result.Ok {
		fmt.Println("✓ Valid error response")
	}

	fmt.Println()
}
