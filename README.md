# gh-vars-migrator

[![Test and Lint](https://github.com/renan-alm/gh-vars-migrator/actions/workflows/test-and-lint.yml/badge.svg)](https://github.com/renan-alm/gh-vars-migrator/actions/workflows/test-and-lint.yml)
[![Release](https://github.com/renan-alm/gh-vars-migrator/actions/workflows/release.yml/badge.svg)](https://github.com/renan-alm/gh-vars-migrator/actions/workflows/release.yml)
[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

GitHub CLI extension for migrating GitHub Actions variables between organizations, repositories, and environments.

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

The `gh-vars-migrator` extension supports two migration modes directly through command-line flags:

### Migration Modes

1. **Organization to Organization**: Migrate organization-level variables
2. **Repository to Repository**: Migrate repository-level variables with automatic environment discovery and migration

### Basic Commands

#### Organization to Organization Migration

Migrate all organization-level variables from one organization to another:

```bash
# Basic migration
gh vars-migrator --source-org myorg --target-org targetorg --org-to-org

# Dry-run mode (preview changes)
gh vars-migrator --source-org myorg --target-org targetorg --org-to-org --dry-run

# Force overwrite existing variables
gh vars-migrator --source-org myorg --target-org targetorg --org-to-org --force
```

#### Repository to Repository Migration

Migrate repository-level variables from one repository to another. The tool automatically discovers all environments in the source repository, creates them in the target if they don't exist, and migrates all environment variables:

```bash
# Basic repo migration (auto-discovers and migrates all environments)
gh vars-migrator --source-org myorg --source-repo myrepo --target-org targetorg --target-repo targetrepo

# Dry-run to preview what would be migrated
gh vars-migrator --source-org myorg --source-repo myrepo --target-org targetorg --target-repo targetrepo --dry-run

# Skip environment variable migration (repo-level variables only)
gh vars-migrator --source-org myorg --source-repo myrepo --target-org targetorg --target-repo targetrepo --skip-envs
```

### Command Options

- `--source-org` (required): Source organization name
- `--source-repo`: Source repository name (required for repo-to-repo)
- `--target-org` (required): Target organization name
- `--target-repo`: Target repository name (required for repo-to-repo)
- `--org-to-org`: Flag to enable organization-level migration mode
- `--skip-envs`: Skip environment variable migration during repo-to-repo (environments are auto-discovered by default)
- `--dry-run`: Preview changes without applying them
- `--force`: Overwrite existing variables in the target

### Global Options

These options work with all commands:

- `--verbose`, `-v`: Enable verbose output

### Mode Detection

The migration mode is automatically detected based on the flags provided:

- If `--org-to-org` flag is set → **Organization migration mode**
- Otherwise → **Repository-to-Repository migration mode** (includes automatic environment discovery and migration)

### Additional Commands

Check authentication status:
```bash
gh vars-migrator auth
```

List variables in an organization:
```bash
gh vars-migrator list --org myorg
```

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
 
