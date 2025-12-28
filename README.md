# zogo

**A TypeScript Zod-inspired schema validation library for Go**

Zogo provides a fluent, type-safe API for validating data structures in Go. Perfect for validating JSON from APIs, configuration files, user input, and more.

```go
userSchema := zogo.Object(zogo.Schema{
    "name":  zogo.String().Min(2).Required(),
    "email": zogo.String().Email().Required(),
    "age":   zogo.Number().Min(18).Optional(),
})

result := userSchema.Parse(jsonData)
if result.Ok {
    // Data is valid!
    validatedData := result.Value
} else {
    // Handle errors
    fmt.Println(result.Errors.Error())
}
```

## Features

### ✅ **15 Built-in Validators**
- **Primitives**: String, Number, Boolean, Date
- **Collections**: Array, Object, Tuple, Record
- **Advanced**: Union, Intersection, Enum, Literal
- **Utilities**: Any, Unknown, Lazy (recursive)

### ✅ **Rich String Validation**
- **Formats**: Email, URL, UUID, IP (v4/v6)
- **Encoding**: Base64, Hex
- **IDs**: CUID, CUID2, ULID, Nanoid
- **Patterns**: Regex, StartsWith, EndsWith, Contains
- **Transforms**: Trim, ToUpperCase, ToLowerCase

### ✅ **Powerful Features**
- 🔄 **Recursive schemas** - Trees, nested comments, file systems
- 🎯 **Discriminated unions** - Type-safe polymorphic data
- 🛡️ **Error paths** - Precise error locations (`user.address[0].zip`)
- 🔧 **Transformations** - Modify data during validation
- 📦 **Zero dependencies** - Pure Go stdlib

## Installation

```bash
go get github.com/hkurdi/zogo
```

## Quick Start

### Basic Validation

```go
import "github.com/hkurdi/zogo"

// String validation
nameSchema := zogo.String().Min(2).Max(50)
result := nameSchema.Parse("John")

// Number validation
ageSchema := zogo.Number().Min(18).Max(120)
result = ageSchema.Parse(25)

// Email validation
emailSchema := zogo.String().Email()
result = emailSchema.Parse("user@example.com")
```

### Object Validation

```go
userSchema := zogo.Object(zogo.Schema{
    "username": zogo.String().Min(3).Max(20),
    "email":    zogo.String().Email(),
    "age":      zogo.Number().Min(18).Optional(),
    "metadata": zogo.Record(zogo.String(), zogo.Any()).Optional(),
})

data := map[string]interface{}{
    "username": "john_doe",
    "email":    "john@example.com",
    "age":      25,
}

result := userSchema.Parse(data)
if result.Ok {
    fmt.Println("Valid user!")
} else {
    for _, err := range result.Errors {
        fmt.Printf("%s: %s\n", err.Path, err.Message)
    }
}
```

### Recursive Schemas

```go
var commentSchema zogo.Validator
commentSchema = zogo.Lazy(func() zogo.Validator {
    return zogo.Object(zogo.Schema{
        "id":      zogo.String(),
        "text":    zogo.String(),
        "replies": zogo.Array(commentSchema).Optional(),
    })
})

// Validates comments with infinite nesting!
```

### Discriminated Unions

```go
successResponse := zogo.Object(zogo.Schema{
    "status": zogo.Literal("success"),
    "data":   zogo.Any(),
})

errorResponse := zogo.Object(zogo.Schema{
    "status":  zogo.Literal("error"),
    "message": zogo.String(),
})

responseSchema := zogo.Union(successResponse, errorResponse)
```

## API Reference

### String Validators

```go
String()
  .Min(length)
  .Max(length)
  .Length(length)
  .Email()
  .URL()
  .UUID()
  .IP() / .IPv4() / .IPv6()
  .Base64()
  .Hex()
  .CUID() / .CUID2()
  .ULID()
  .Nanoid()
  .Regex(pattern)
  .StartsWith(prefix)
  .EndsWith(suffix)
  .Contains(substring)
  .Trim()
  .ToLowerCase()
  .ToUpperCase()
  .Required() / .Optional() / .Nullable()
  .Default(value)
  .Refine(check, message)
```

### Number Validators

```go
Number()
  .Min(value)
  .Max(value)
  .Int()
  .Positive() / .Negative()
  .NonNegative() / .NonPositive()
  .Finite()
  .Safe()
  .MultipleOf(value)
  .Required() / .Optional() / .Nullable()
  .Default(value)
  .Refine(check, message)
```

### Object Validators

```go
Object(Schema{...})
  .Strict()      // Error on unknown fields
  .Passthrough() // Keep unknown fields
  .Strip()       // Remove unknown fields (default)
  .Required() / .Optional() / .Nullable()
```

### Array Validators

```go
Array(elementValidator)
  .Min(length)
  .Max(length)
  .Length(length)
  .NonEmpty()
  .Required() / .Optional() / .Nullable()
```

### Advanced Validators

```go
// Union - OR logic
Union(String(), Number())

// Intersection - AND logic  
Intersection(String().Email(), String().Min(5))

// Tuple - Fixed-length arrays
Tuple(String(), Number(), Boolean())

// Record - Typed dictionaries
Record(String(), Number())

// Enum - Value sets
Enum([]interface{}{"active", "inactive", "pending"})

// Literal - Exact values
Literal("success")

// Lazy - Recursive schemas
Lazy(func() Validator { return ... })

// Date
Date().Past() / .Future() / .Min(date) / .Max(date)
```

## Error Handling

```go
result := schema.Parse(data)

if !result.Ok {
    // Get all errors
    fmt.Println(result.Errors.Error())
    
    // Get first error
    first := result.Errors.First()
    
    // Check specific path
    if result.Errors.HasPath("email") {
        emailErrors := result.Errors.ByPath("email")
    }
    
    // Get structured issues (for JSON APIs)
    issues := result.Errors.Issues()
}
```

## Examples

See [examples/api-validation](examples/api-validation) for comprehensive examples including:
- User registration with complex validation
- Recursive blog post comments
- E-commerce orders with discriminated unions
- Configuration validation with defaults
- API response validation

## Comparison with go-playground/validator

| Feature | zogo | go-playground/validator |
|---------|------|------------------------|
| API Style | Fluent/Chainable | Struct tags |
| Dynamic Schemas | ✅ Yes | ❌ No |
| Recursive Schemas | ✅ Yes | Limited |
| Discriminated Unions | ✅ Yes | ❌ No |
| Error Paths | ✅ Precise | Tag-based |
| Type Safety | ✅ Compile-time | Runtime reflection |
| JSON Validation | ✅ First-class | Struct-dependent |
| Transformations | ✅ Yes | ❌ No |

## Roadmap

- [x] Core validators (String, Number, Boolean, Object, Array)
- [x] Advanced validators (Union, Intersection, Tuple, Record)
- [x] Recursive schemas (Lazy)
- [x] String formats (Email, URL, UUID, IP, Base64, etc.)
- [x] Transformations (Trim, case conversion)
- [x] Error handling improvements
- [ ] Object helpers (Partial, Pick, Omit, Extend, Merge)
- [ ] Async validation
- [ ] Custom error messages
- [ ] Benchmarks and performance optimization

## Contributing

Contributions are welcome! Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details.

## License

MIT License - see [LICENSE](LICENSE) for details.

## Inspiration

Zogo is inspired by [Zod](https://github.com/colinhacks/zod) from the TypeScript ecosystem, bringing its elegant API design to Go.