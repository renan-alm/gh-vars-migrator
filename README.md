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

### Authentication

The tool supports multiple authentication methods:

1. **Explicit Tokens (Recommended for cross-account migrations)**: Use `--source-pat` and `--target-pat` flags or `SOURCE_PAT` and `TARGET_PAT` environment variables to specify separate tokens for source and target operations.

2. **GITHUB_TOKEN Fallback**: If `GITHUB_TOKEN` environment variable is set, it will be used for both source and target when explicit PATs are not provided.

3. **GitHub CLI Authentication**: If no tokens are provided, the tool falls back to GitHub CLI's authentication (requires `gh auth login`).

#### Authentication Examples

```bash
# Using explicit PATs for different accounts
gh vars-migrator --source-org srcorg --target-org tgtorg --org-to-org \
  --source-pat ghp_sourcetoken123 --target-pat ghp_targettoken456

# Using environment variables
export SOURCE_PAT=ghp_sourcetoken123
export TARGET_PAT=ghp_targettoken456
gh vars-migrator --source-org srcorg --target-org tgtorg --org-to-org

# Using GITHUB_TOKEN for both source and target
export GITHUB_TOKEN=ghp_yourtoken
gh vars-migrator --source-org srcorg --target-org tgtorg --org-to-org

# Using GitHub CLI authentication (default)
gh auth login
gh vars-migrator --source-org srcorg --target-org tgtorg --org-to-org
```

### Migration Modes

1. **Organization to Organization**: Migrate organization-level variables
2. **Repository to Repository**: Migrate repository-level variables with automatic environment discovery and migration

### Basic Commands

#### Organization to Organization Migration

Migrate all organization-level variables from one organization to another:

```bash
# Basic migration (preserves source visibility for each variable)
gh vars-migrator --source-org myorg --target-org targetorg --org-to-org

# Dry-run mode (preview changes)
gh vars-migrator --source-org myorg --target-org targetorg --org-to-org --dry-run

# Force overwrite existing variables
gh vars-migrator --source-org myorg --target-org targetorg --org-to-org --force
```

**Organization variable visibility**

GitHub organization variables have a visibility scope that controls which repositories can access them:

| Scope | Description |
|-------|-------------|
| `all` | Accessible by all repositories in the organization |
| `private` | Accessible only by private repositories |
| `selected` | Accessible only by explicitly selected repositories |

`gh-vars-migrator` automatically preserves the source variable's visibility when migrating. For variables with `selected` visibility, the tool fetches the selected repository names from the source organization and matches them by name in the target organization. Only repositories whose names exist in both organizations are included in the target's selection list. If no matching repositories are found, the variable is created with an empty selection (zero repositories).

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

#### Data Residency Migration

Organizations with strict data residency requirements can specify custom GitHub hostnames to control which API endpoints are used for the migration. Variable values travel only between the configured source and target endpoints, keeping data within your approved infrastructure.

Use `--source-hostname` and `--target-hostname` to target GitHub Enterprise Server (GHES) instances or data-residency-compliant GitHub Enterprise Cloud (GHEC) endpoints:

```bash
# Migrate between two GitHub Enterprise Server instances
gh vars-migrator --source-org myorg --target-org targetorg --org-to-org \
  --source-hostname github.source-company.com \
  --target-hostname github.target-company.com \
  --source-pat ghp_sourcetoken \
  --target-pat ghp_targettoken

# Migrate from a GHES instance to GitHub.com
gh vars-migrator --source-org myorg --target-org targetorg --org-to-org \
  --source-hostname github.mycompany.com \
  --source-pat ghp_sourcetoken \
  --target-pat ghp_targettoken

# Migrate from GitHub.com to a data-residency GHEC endpoint
gh vars-migrator --source-org myorg --source-repo myrepo \
  --target-org targetorg --target-repo targetrepo \
  --target-hostname api.mycompany.ghe.com \
  --source-pat ghp_sourcetoken \
  --target-pat ghp_targettoken

# Using GitHub CLI credentials stored for a specific host
# (requires: gh auth login --hostname github.mycompany.com)
gh vars-migrator --source-org myorg --target-org targetorg --org-to-org \
  --source-hostname github.mycompany.com
```

### Command Options

#### Source and Target
- `--source-org` (required): Source organization name
- `--source-repo`: Source repository name (required for repo-to-repo)
- `--target-org` (required): Target organization name
- `--target-repo`: Target repository name (required for repo-to-repo)

#### Authentication
- `--source-pat`: Source personal access token (env: `SOURCE_PAT`)
- `--target-pat`: Target personal access token (env: `TARGET_PAT`)
- If neither PAT is provided, falls back to `GITHUB_TOKEN` or GitHub CLI auth

#### Data Residency
- `--source-hostname`: Custom GitHub hostname for the source (e.g., `github.mycompany.com`). Use for GitHub Enterprise Server or data-residency GitHub Enterprise Cloud instances.
- `--target-hostname`: Custom GitHub hostname for the target (e.g., `github.mycompany.com`). Use for GitHub Enterprise Server or data-residency GitHub Enterprise Cloud instances.
- When a hostname flag is omitted, the corresponding client defaults to `github.com`.

#### Mode Options
- `--org-to-org`: Flag to enable organization-level migration mode
- `--skip-envs`: Skip environment variable migration during repo-to-repo (environments are auto-discovered by default)

#### Behavior Options
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


## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
 
