# Testing Documentation

This directory contains comprehensive tests for the gh-vars-migrator project.

## Test Structure

### Unit Tests

Unit tests are located alongside the source code in each package:

- **`internal/client/client_test.go`**: Tests for GitHub API client operations
  - Path construction for all API endpoints
  - Request body formatting for create/update operations
  - Tests validate the logic without requiring actual GitHub credentials
  
- **`internal/config/config_test.go`**: Tests for configuration validation
  - Validation for all migration modes (repo-to-repo, org-to-org)
  - Description generation for different configurations
  - Edge cases and error scenarios
  - Coverage: 80.0%

- **`internal/logger/logger_test.go`**: Tests for logging functionality
  - All logging functions (Info, Success, Warning, Error, Debug, Plain)
  - Summary printing with various combinations of counts
  - Output formatting validation
  - Coverage: 100.0%

- **`internal/migrator/migrator_test.go`**: Tests for migration logic
  - Configuration validation
  - Migration mode behaviors (repo-to-repo with auto-discovery, org-to-org)
  - Dry-run and force update scenarios
  - Error handling and accumulation
  - Tests validate logic paths without requiring actual API calls

- **`internal/types/types_test.go`**: Tests for core data types
  - MigrationResult methods (AddError, HasErrors, Total)
  - Migration mode constants
  - Coverage: 100.0%

### Integration Tests

Integration tests are located in `tests/integration/`:

- **`integration_test.go`**: End-to-end workflow tests
  - Configuration validation workflows
  - Migration result tracking across full workflows
  - Dry-run and force update workflows
  - Environment auto-discovery scenarios
  - Error handling across the entire migration pipeline
  - Configuration descriptions for all modes

## Running Tests

### Run All Tests
```bash
make test
```

### Run Tests with Coverage
```bash
make test-coverage
```

This generates:
- `coverage.out`: Coverage profile
- `coverage.html`: HTML coverage report (if `go tool cover` is available)

### Run Tests for Specific Package
```bash
go test -v ./internal/config/...
go test -v ./internal/migrator/...
go test -v ./tests/integration/...
```

### Run Tests with Race Detector
```bash
go test -race ./...
```

## Test Coverage Summary

| Package | Coverage | Notes |
|---------|----------|-------|
| internal/config | 80.0% | Comprehensive validation and description tests |
| internal/logger | 100.0% | All logging functions tested |
| internal/types | 100.0% | Complete coverage of data types and methods |
| internal/client | Path/Body tests | API logic tested, actual calls require credentials |
| internal/migrator | Logic tests | Migration logic tested, actual calls require credentials |

## Testing Philosophy

### What We Test

1. **Business Logic**: All configuration validation, result tracking, and workflow logic
2. **Data Transformations**: Path construction, request body formatting, response parsing
3. **Error Handling**: Error accumulation, validation failures, edge cases
4. **Integration Flows**: End-to-end workflows combining multiple components

### What We Don't Test

1. **External API Calls**: Actual GitHub API calls require authentication and would make tests brittle
2. **Network I/O**: Tests focus on logic rather than network reliability
3. **UI/CLI Interactions**: Current main.go is a placeholder and doesn't have CLI parsing yet

### Testing Approach

#### For Client Package
- Tests verify API path construction and request body formatting
- Mock-based tests are avoided due to the external API client structure
- Focus on validating the logic that prepares API calls

#### For Migrator Package  
- Tests validate configuration handling and migration logic
- Dry-run and force mode behaviors are tested without actual API calls
- Error handling and result tracking are thoroughly tested

#### For Integration Tests
- Tests combine multiple components to validate end-to-end workflows
- Configuration validation → migration execution → result tracking
- Tests work without external dependencies or credentials

## Adding New Tests

When adding new functionality:

1. **Add unit tests** in the same package as the code
2. **Add integration tests** if the feature involves multiple components
3. **Update test coverage** by running `make test-coverage`
4. **Document edge cases** that your tests cover

### Test Naming Conventions

- Unit test functions: `Test<FunctionName>_<Scenario>`
- Integration test functions: `TestEndToEnd_<Workflow>`
- Table-driven tests: Use descriptive test names in the `tests` slice

Example:
```go
func TestValidate_RepoToRepo(t *testing.T) {
    tests := []struct {
        name    string
        cfg     *types.MigrationConfig
        wantErr bool
    }{
        {
            name: "valid config",
            cfg:  &types.MigrationConfig{...},
            wantErr: false,
        },
        // ...
    }
}
```

## Continuous Integration

Tests are automatically run in CI on:
- Every push to a branch
- Every pull request
- Before merging to main

CI failures should be investigated and fixed before merging.

## Future Test Improvements

1. **CLI Integration Tests**: Once CLI parsing is implemented, add tests for command-line argument handling
2. **Mock GitHub API**: Consider adding more sophisticated mock-based tests for the client package
3. **Performance Tests**: Add benchmarks for critical paths (if needed)
4. **Contract Tests**: Add tests to verify our expectations about GitHub API responses
