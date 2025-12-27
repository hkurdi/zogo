package zogo

// Validator is the core interface that all validators implement
type Validator interface {
	// Parse validates the input value and returns a ParseResult
	Parse(value any) ParseResult
}
