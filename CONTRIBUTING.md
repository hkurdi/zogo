# Contributing to zogo

Thanks for your interest in contributing to zogo!

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/YOUR_USERNAME/zogo.git`
3. Create a branch: `git checkout -b my-feature`
4. Make your changes
5. Run tests: `go test ./...`
6. Commit: `git commit -am 'Add some feature'`
7. Push: `git push origin my-feature`
8. Open a Pull Request

## Development Setup

```bash
# Clone the repo
git clone https://github.com/hkurdi/zogo.git
cd zogo

# Install dependencies
go mod download

# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...
```

## Code Style

- Follow standard Go conventions
- Run `go fmt` before committing
- Add tests for new features
- Keep functions small and focused
- Write clear commit messages

## Testing

- All new features must have tests
- Aim for >80% test coverage
- Test both success and error cases
- Use table-driven tests where appropriate

## Pull Request Guidelines

- Keep PRs focused on a single feature/fix
- Update README.md if adding new features
- Add examples for new validators
- Ensure all tests pass
- Update CHANGELOG.md (if we add one)

## Questions?

Open an issue or reach out to [@hkurdi](https://github.com/hkurdi)
