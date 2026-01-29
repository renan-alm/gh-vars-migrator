# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v0.2.0] - 2026-01-29

### Added
- **Resilient API Client**: Implemented centralized retry logic with exponential backoff for GitHub API calls to handle transient errors and rate limits
- **Context Propagation**: Added cancellation and timeout support by propagating `context.Context` throughout the application
- **Rate Limit Handling**: Automatic handling of `Retry-After` headers and 429/403 responses

### Changed
- Updated internal client to use `hashicorp/go-retryablehttp`
- Refactored `Migrator` and `Client` interfaces to accept `context.Context`

## [v0.1.0] - 2026-01-23

### Added
- **Organization to Organization Migration**: Migrate organization-level GitHub Actions variables between organizations
- **Repository to Repository Migration**: Migrate repository-level variables between repositories
- **Automatic Environment Discovery**: Auto-discover all environments in source repository, create them in target if they don't exist, and migrate all environment variables
- **Dry-run Mode**: Preview changes without applying them using `--dry-run` flag
- **Force Mode**: Overwrite existing variables in target using `--force` flag
- **Skip Environments**: Option to skip environment migration with `--skip-envs` flag
- **Authentication Check**: `auth` subcommand to verify GitHub CLI authentication status
- **List Variables**: `list` subcommand to list variables in an organization
- **Verbose Output**: `--verbose` flag for detailed logging
- GitHub CLI extension support (`gh extension install`)
- Cross-platform binaries (Linux, macOS, Windows) for amd64 and arm64
- Comprehensive test suite with unit and integration tests
- CI/CD workflows for testing and releases
- Dockerfile for containerized builds

### Technical
- Built with Go 1.25+
- Uses GitHub CLI's go-gh library for API interactions
- Cobra-based CLI with flag-based command pattern
