---
description: 'Instructions for refactoring Go code following idiomatic practices and maintaining code quality'
applyTo: '**/*.go'
---

# Go Refactoring Instructions

Follow these instructions when refactoring Go code to improve readability, maintainability, and performance while preserving existing behavior.

## General Refactoring Principles

- Preserve existing functionality and behavior
- Make incremental, focused changes
- Ensure tests pass before and after refactoring
- Favor clarity and simplicity over cleverness
- Keep the happy path left-aligned (minimize indentation)
- Return early to reduce nesting
- Prefer early return over if-else chains; use `if condition { return }` pattern

## Code Simplification

### Reduce Complexity

- Extract long functions into smaller, focused functions
- Replace nested conditionals with guard clauses
- Eliminate dead code and unused variables
- Remove redundant type conversions
- Simplify boolean expressions

### Improve Readability

- Rename variables and functions to be more descriptive
- Replace magic numbers with named constants
- Group related code together
- Add blank lines to separate logical sections
- Remove unnecessary comments that state the obvious

## Structural Refactoring

### Function Extraction

- Extract repeated code into helper functions
- Keep functions focused on a single responsibility
- Limit function length (aim for under 50 lines)
- Extract complex conditionals into well-named functions

### Package Organization

- Move types closer to where they're used
- Split large files into smaller, focused files
- Use `internal/` for packages that shouldn't be exported
- Avoid circular dependencies

### Type Improvements

- Replace primitive types with domain-specific types when it adds clarity
- Convert repeated struct patterns into reusable types
- Use embedding for composition instead of copy-pasting fields
- Define interfaces close to where they're used

## Error Handling Refactoring

- Wrap errors with context using `fmt.Errorf` with `%w` verb
- Replace `panic` with proper error returns where appropriate
- Consolidate duplicate error handling logic
- Use sentinel errors or custom error types for domain errors
- Remove error swallowing (don't ignore errors with `_` without justification)

## Performance Refactoring

- Preallocate slices when size is known
- Use `strings.Builder` for string concatenation in loops
- Replace repeated map lookups with single lookups
- Use `sync.Pool` for frequently allocated objects
- Avoid unnecessary allocations in hot paths

## Concurrency Refactoring

- Replace shared memory with channel communication where appropriate
- Ensure goroutines have clear exit conditions
- Use `context.Context` for cancellation propagation
- Replace manual synchronization with `sync` primitives when clearer
- WaitGroup usage by Go version:
	- If `go >= 1.25`, prefer `WaitGroup.Go` method
	- If `go < 1.25`, use classic `Add`/`Done` pattern

## Interface Refactoring

- Extract interfaces from concrete implementations
- Keep interfaces small (1-3 methods)
- Accept interfaces, return concrete types
- Remove unused interface methods
- Define interfaces at the call site, not the implementation

## Test Refactoring

- Convert repetitive tests to table-driven tests
- Extract common setup into helper functions marked with `t.Helper()`
- Use subtests with `t.Run` for better organization
- Replace test-specific implementations with mocks or fakes
- Add missing edge case tests while refactoring

## Code Smells to Address

- Long parameter lists → use option structs or functional options
- Feature envy → move methods to the type they operate on most
- Shotgun surgery → consolidate related changes into single packages
- Primitive obsession → introduce domain types
- Large structs → split into focused types
- God functions → extract smaller, single-purpose functions

## Refactoring Checklist

Before refactoring:
- [ ] Understand the existing behavior
- [ ] Ensure tests exist and pass
- [ ] Identify the specific improvement goal

During refactoring:
- [ ] Make small, incremental changes
- [ ] Run tests frequently
- [ ] Use `go fmt` and `go vet` after changes
- [ ] Preserve package declarations (never duplicate)

After refactoring:
- [ ] Verify all tests pass
- [ ] Run `golangci-lint` for additional checks
- [ ] Review the diff for unintended changes
- [ ] Update documentation if behavior changed

## Critical Rules

- **NEVER duplicate `package` declarations** - preserve existing package lines
- **NEVER change exported API signatures** without explicit request
- **ALWAYS run tests** before considering refactoring complete
- **PRESERVE behavior** - refactoring changes structure, not functionality
