# zogo

Schema validation for Go, inspired by [Zod](https://github.com/colinhacks/zod).

```go
import "github.com/hkurdi/zogo"

// Validate incoming API data
userSchema := zogo.Object(zogo.Schema{
    "email": zogo.String().Email().Required(),
    "age":   zogo.Number().Min(18).Max(120).Required(),
    "tags":  zogo.Array(zogo.String()).Optional(),
})

result := userSchema.Parse(data)
if !result.Ok {
    // Handle validation errors
    for _, err := range result.Errors {
        fmt.Printf("%s: %s\n", err.Path, err.Message)
    }
    return
}

// Use validated data
fmt.Println("Valid user data:", result.Value)
```

## Why zogo?

- **Familiar API** - If you know Zod, you know zogo
- **Type-safe** - Built with Go generics, fluent chainable methods
- **Validate anything** - JSON APIs, user input, configuration files
- **Clear errors** - Structured errors with field paths and messages

## Installation

```bash
go get github.com/hkurdi/zogo
```

## Features

### Core Types

**Primitives**
- `String()` - min/max length, email, URL, UUID, regex, startsWith, endsWith, contains
- `Number()` - min/max, int, positive, negative, finite, safe, multipleOf
- `Boolean()`
- `Date()` - min/max dates
- `BigInt()` - min/max
- `Literal()` - exact value matching
- `Enum()` - string/number enums
- `Any()` - accepts anything
- `Unknown()` - like any but safer
- `Void()` - undefined/null only
- `Undefined()` - undefined only
- `Null()` - null only
- `Never()` - no valid input

**Complex Types**
- `Object()` - nested schemas, strict/passthrough/strip modes
- `Array()` - typed arrays, min/max length, nonempty
- `Tuple()` - fixed-length arrays with different types
- `Record()` - map[string]T validation
- `Map()` - map validation
- `Set()` - set validation

**Advanced Types**
- `Union()` - OR types (A | B | C)
- `Intersection()` - AND types (A & B)
- `Discriminated()` - discriminated unions
- `Lazy()` - recursive/circular schemas

### Modifiers & Transformations

- `Required()` / `Optional()` - field presence
- `Nullable()` - allow null
- `Default()` - default values
- `Catch()` - fallback on error
- `Transform()` - map validated values
- `Refine()` - custom validation logic
- `Preprocess()` - transform before validation

### String Methods

```go
schema := zogo.String().
    Min(5).                    // Minimum length
    Max(100).                  // Maximum length
    Length(10).                // Exact length
    Email().                   // Valid email
    URL().                     // Valid URL
    UUID().                    // Valid UUID
    Regex("^[a-z]+$").        // Match pattern
    StartsWith("https://").    // Must start with
    EndsWith(".com").          // Must end with
    Contains("@").             // Must contain
    Trim().                    // Remove whitespace
    ToLowerCase().             // Convert to lowercase
    ToUpperCase().             // Convert to uppercase
    Required()                 // Cannot be null/undefined
```

### Number Methods

```go
schema := zogo.Number().
    Min(0).                    // Minimum value
    Max(100).                  // Maximum value
    Int().                     // Must be integer
    Positive().                // Must be > 0
    Negative().                // Must be < 0
    NonNegative().             // Must be >= 0
    NonPositive().             // Must be <= 0
    Finite().                  // No Infinity/NaN
    Safe().                    // Safe integer range
    MultipleOf(5).             // Divisible by N
    Required()
```

### Object Methods

```go
baseSchema := zogo.Object(zogo.Schema{
    "name": zogo.String(),
    "age":  zogo.Number(),
})

// Extend with new fields
extended := baseSchema.Extend(zogo.Schema{
    "email": zogo.String().Email(),
})

// Merge two schemas
merged := schema1.Merge(schema2)

// Pick specific fields
picked := schema.Pick("name", "email")

// Omit specific fields
omitted := schema.Omit("password")

// Make all fields optional
partial := schema.Partial()

// Make all fields required
required := schema.Required()

// Recursively make optional
deepPartial := schema.DeepPartial()

// Unknown field handling
strict := schema.Strict()        // Error on unknown fields
passthrough := schema.Passthrough() // Keep unknown fields
strip := schema.Strip()          // Remove unknown fields
```

## Quick Start

### Basic Validation

```go
import "github.com/hkurdi/zogo"

// String validation
emailSchema := zogo.String().Email()
result := emailSchema.Parse("[email protected]")

// Number validation
ageSchema := zogo.Number().Min(18).Max(120)
result := ageSchema.Parse(25)

// Boolean validation
boolSchema := zogo.Boolean()
result := boolSchema.Parse(true)
```

### Object Validation

```go
userSchema := zogo.Object(zogo.Schema{
    "name":  zogo.String().Min(2).Required(),
    "email": zogo.String().Email().Required(),
    "age":   zogo.Number().Min(18).Optional(),
    "tags":  zogo.Array(zogo.String()).Optional(),
})

var data map[string]interface{}
json.Unmarshal(jsonBytes, &data)

result := userSchema.Parse(data)
if result.Ok {
    fmt.Println("Valid!", result.Value)
}
```

### Nested Objects

```go
postSchema := zogo.Object(zogo.Schema{
    "title": zogo.String().Min(5).Max(100).Required(),
    "author": zogo.Object(zogo.Schema{
        "name":  zogo.String().Required(),
        "email": zogo.String().Email().Required(),
    }).Required(),
    "tags": zogo.Array(zogo.String()).Min(1).Required(),
    "metadata": zogo.Object(zogo.Schema{
        "views":     zogo.Number().Optional(),
        "likes":     zogo.Number().Optional(),
        "createdAt": zogo.Date().Required(),
    }).Optional(),
})
```

### Array Validation

```go
// Array of strings
tagsSchema := zogo.Array(zogo.String()).Min(1).Max(5)

// Array of objects
usersSchema := zogo.Array(
    zogo.Object(zogo.Schema{
        "name":  zogo.String().Required(),
        "email": zogo.String().Email().Required(),
    }),
).Nonempty()

// Tuple (fixed length, different types)
tupleSchema := zogo.Tuple(
    zogo.String(),
    zogo.Number(),
    zogo.Boolean(),
)
result := tupleSchema.Parse([]interface{}{"hello", 42, true})
```

### Transformations

```go
// Trim and lowercase email
emailSchema := zogo.String().Email().Trim().ToLowerCase()
result := emailSchema.Parse("  [email protected]  ")
// result.Value = "[email protected]"

// Round number to integer
priceSchema := zogo.Number().Int()
result := priceSchema.Parse(19.99)
// result.Value = 20

// Default value
nameSchema := zogo.String().Default("Anonymous")
result := nameSchema.Parse(nil)
// result.Value = "Anonymous"
```

### Custom Validation

```go
// Using Refine
passwordSchema := zogo.String().
    Min(8).
    Refine(func(s string) bool {
        return strings.ContainsAny(s, "ABCDEFGHIJKLMNOPQRSTUVWXYZ")
    }, "Password must contain uppercase letter").
    Refine(func(s string) bool {
        return strings.ContainsAny(s, "0123456789")
    }, "Password must contain number")

// Using Transform
slugSchema := zogo.String().Transform(func(s string) (interface{}, error) {
    return strings.ToLower(strings.ReplaceAll(s, " ", "-")), nil
})
```

### Union Types

```go
// String OR Number
schema := zogo.Union(
    zogo.String(),
    zogo.Number(),
)
result := schema.Parse("hello") // OK
result = schema.Parse(42)       // OK
result = schema.Parse(true)     // Error

// Discriminated Union
shapeSchema := zogo.Discriminated("type", zogo.Schema{
    "circle": zogo.Object(zogo.Schema{
        "type":   zogo.Literal("circle"),
        "radius": zogo.Number(),
    }),
    "rectangle": zogo.Object(zogo.Schema{
        "type":   zogo.Literal("rectangle"),
        "width":  zogo.Number(),
        "height": zogo.Number(),
    }),
})
```

### Enums

```go
// String enum
roleSchema := zogo.Enum("admin", "user", "guest")

// Number enum
statusSchema := zogo.Enum(200, 404, 500)
```

### Error Handling

```go
result := userSchema.Parse(data)
if !result.Ok {
    // Get all errors
    for _, err := range result.Errors {
        fmt.Printf("Field '%s': %s\n", err.Path, err.Message)
    }
    
    // Or get first error
    if firstErr := result.Errors.First(); firstErr != nil {
        fmt.Println(firstErr.Error())
    }
}

// Output example:
// Field 'email': Invalid email format
// Field 'age': Must be at least 18
// Field 'tags.0': Expected string, received number
```

### Safe Parsing

```go
// Parse returns ParseResult instead of panicking
result := schema.Parse(data)
if result.Ok {
    // Use result.Value
} else {
    // Handle result.Errors
}

// Versus Parse which panics on error (coming soon)
value := schema.MustParse(data) // Panics if invalid
```

## Advanced Examples

### Recursive Schemas

```go
type Category struct {
    Name     string
    Children []Category
}

categorySchema := zogo.Lazy(func() zogo.Validator {
    return zogo.Object(zogo.Schema{
        "name":     zogo.String().Required(),
        "children": zogo.Array(categorySchema).Optional(),
    })
})
```

### Preprocessing

```go
// Parse string to number before validation
schema := zogo.Preprocess(
    func(val interface{}) (interface{}, error) {
        if s, ok := val.(string); ok {
            return strconv.Atoi(s)
        }
        return val, nil
    },
    zogo.Number().Min(0),
)

result := schema.Parse("42") // Converts "42" to 42, then validates
```

### Error Messages

```go
// Custom error messages
schema := zogo.String().
    Min(8, "Password must be at least 8 characters").
    Email("Please provide a valid email address")

// Or use error map
schema := zogo.Object(zogo.Schema{
    "email": zogo.String().Email(),
    "age":   zogo.Number().Min(18),
}).WithErrorMap(map[string]string{
    "email": "Invalid email address",
    "age":   "Must be 18 or older",
})
```

## Comparison with go-playground/validator

| Feature | zogo | go-playground/validator |
|---------|------|------------------------|
| API Style | Fluent/Chainable | Struct tags |
| Type Safety | ✅ Full | ⚠️ String tags |
| Dynamic Schemas | ✅ Yes | ❌ No |
| Validate Maps/JSON | ✅ Easy | ⚠️ Harder |
| Validate Structs | ✅ Yes | ✅ Yes |
| Transformations | ✅ Built-in | ❌ Manual |
| Composability | ✅ High | ⚠️ Limited |

**Use zogo when:**
- Validating external JSON/API data
- Building schemas dynamically
- Coming from TypeScript/Zod
- Want fluent, type-safe API

**Use go-playground/validator when:**
- Only validating your own structs
- Prefer declarative struct tags
- Need maximum performance

## Contributing

Contributions welcome! Please open an issue or PR.

## License

MIT © 2025 Hamza L Kurdi

## Roadmap

- [ ] Core types (String, Number, Boolean)
- [ ] Complex types (Object, Array)
- [ ] Advanced types (Union, Intersection, Tuple)
- [ ] Date validation
- [ ] Custom error messages
- [ ] Async validation
- [ ] JSON Schema export
- [ ] Performance benchmarks
- [ ] Full test coverage

## Inspiration

This project is heavily inspired by:
- [Zod](https://github.com/colinhacks/zod) - TypeScript-first schema validation
- [go-playground/validator](https://github.com/go-playground/validator) - Go struct validation
