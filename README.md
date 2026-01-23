# gh-vars-migrator

GitHub CLI extension for variables migration between GitHub Organizations.

## Installation

### Prerequisites

- [GitHub CLI](https://cli.github.com/) (gh) installed and authenticated
- Go 1.25.0 or later (for building from source)

### Installing as a GitHub CLI Extension

Install directly from GitHub:

```bash
gh extension install renan-alm/gh-vars-migrator
```

### Manual Installation

1. Clone the repository:
```bash
git clone https://github.com/renan-alm/gh-vars-migrator.git
cd gh-vars-migrator
```

2. Build and install:
```bash
make install
```

Or install using Go:
```bash
go install github.com/renan-alm/gh-vars-migrator@latest
```

## Usage

Once installed, run the extension using:

```bash
gh vars-migrator
```

The extension will authenticate using your GitHub CLI credentials and display your authenticated user information.

## Development

### Building from Source

Build the binary:

```bash
make build
```

The compiled binary will be in the `bin/` directory.

### Testing

Run tests:

```bash
make test
```

Run tests with coverage:

```bash
make test-coverage
```

### Linting

Run the linter (requires [golangci-lint](https://golangci-lint.run/usage/install/)):

```bash
make lint
```

### Available Make Targets

- `make build` - Build the binary
- `make test` - Run tests
- `make test-coverage` - Run tests with coverage report
- `make lint` - Run linting
- `make install` - Build and install the binary
- `make clean` - Remove build artifacts
- `make help` - Display help message

## Release Process

This project uses GitHub Actions to automatically build and release binaries for multiple platforms when a new tag is pushed.

To create a new release:

1. Tag a new version:
```bash
git tag v1.0.0
git push origin v1.0.0
```

2. GitHub Actions will automatically:
   - Build binaries for multiple platforms (Linux, macOS, Windows)
   - Create a GitHub release
   - Upload the binaries as release assets

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
 
