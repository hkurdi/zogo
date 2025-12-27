# zog

TypeScript-inspired schema validation for Go.

```go
import "github.com/hlk/zog"

// Define a schema
schema := zog.Object(zog.Schema{
    "email": zog.String().Email().Min(5),
    "age":   zog.Number().Min(18).Max(120),
    "tags":  zog.Array(zog.String()).Min(1),
})

// Parse and validate
var data map[string]interface{}
json.Unmarshal(jsonBytes, &data)

result := schema.Parse(data)
if !result.Ok {
    for _, err := range result.Errors {
        fmt.Printf("%s: %s\n", err.Path, err.Message)
    }
}
```

## Why zog?

If you've used [Zod](https://github.com/colinhacks/zod) in TypeScript and want that same ergonomic validation in Go, zog is for you.

- **Fluent API**: Chain validators naturally
- **Type-safe**: Built with Go generics
- **Clear errors**: Get structured, actionable error messages
- **Zero dependencies**: Just the Go standard library

## Installation

```bash
go get github.com/hlk/zog
```

## Status

🚧 **Early Development** - API may change. Not production-ready yet.

## License

MIT
