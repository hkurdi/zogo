# Setup Instructions

## Initial Setup

1. **Create the directory structure:**
```bash
mkdir -p zogo/examples/basic
cd zogo
```

2. **Initialize the Go module:**
```bash
go mod init github.com/hkurdi/zogo
```

3. **Copy all files:**
   - LICENSE
   - README.md
   - .gitignore
   - CONTRIBUTING.md

4. **Initialize git repository:**
```bash
git init
git add .
git commit -m "Initial commit: Project setup"
```

5. **Create GitHub repository:**
   - Go to https://github.com/new
   - Repository name: `zogo`
   - Description: "Schema validation for Go, inspired by Zod"
   - Public repository
   - **Don't** initialize with README (we have our own)

6. **Push to GitHub:**
```bash
git remote add origin https://github.com/hkurdi/zogo.git
git branch -M main
git push -u origin main
```

## Next Steps

After setup, we'll create files in this order:

1. **Core files:**
   - `zogo.go` - Package documentation + core types
   - `validator.go` - Validator interface
   - `errors.go` - Error types
   - `result.go` - ParseResult type

2. **Validators (in order):**
   - `string.go` + `string_test.go`
   - `number.go` + `number_test.go`
   - `boolean.go` + `boolean_test.go`
   - `array.go` + `array_test.go`
   - `object.go` + `object_test.go`

3. **Advanced validators:**
   - `date.go`
   - `enum.go`
   - `union.go`
   - `literal.go`
   - etc.

4. **Examples:**
   - `examples/basic/main.go`
   - `examples/api_validation/main.go`
   - `examples/nested/main.go`

## Verify Setup

```bash
# Should show go.mod
ls -la

# Should be on main branch
git branch

# Should show github remote
git remote -v
```

You're ready to start coding! 🚀
