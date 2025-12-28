# Changelog

All notable changes to this project will be documented in this file.

## [0.1.0] - 2025-12-28

### Added
- **Core Validators**
  - String validator with 15+ format validators
  - Number validator with comprehensive constraints
  - Boolean validator
  - Object validator with Strict/Passthrough/Strip modes
  - Array validator with typed elements
  - Date validator with temporal validation
  - Enum validator for value sets
  - Literal validator for exact value matching
  - Any and Unknown validators

- **Advanced Validators**
  - Union validator for OR logic and discriminated unions
  - Intersection validator for AND logic with transformation chaining
  - Tuple validator for fixed-length typed arrays
  - Record validator for typed dictionaries
  - Lazy validator for recursive/self-referential schemas

- **String Formats**
  - Email, URL, UUID validation
  - IP (IPv4, IPv6) validation
  - Base64 and Hex encoding validation
  - CUID, CUID2, ULID, Nanoid format validation
  - Regex, StartsWith, EndsWith, Contains patterns

- **Features**
  - Full modifier support (Required/Optional/Nullable/Default/Refine)
  - String transformations (Trim, ToUpperCase, ToLowerCase)
  - Number constraints (Min/Max/Int/Positive/Negative/MultipleOf/Safe/Finite)
  - Error path tracking with precise locations
  - Enhanced error handling with codes and helper methods
  - Zero external dependencies

- **Examples**
  - Comprehensive API validation example
  - User registration validation
  - Recursive blog post comments
  - E-commerce order with discriminated unions
  - Configuration validation with defaults
  - API response validation

### Documentation
- Complete README with API reference
- Contributing guidelines
- MIT License
- Comprehensive test coverage

[0.1.0]: https://github.com/hkurdi/zogo/releases/tag/v0.1.0