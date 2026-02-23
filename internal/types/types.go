package types

import (
	"errors"
	"time"
)

// Error definitions for migration operations
var (
	ErrInvalidConfig      = errors.New("invalid configuration")
	ErrMissingSourceOwner = errors.New("missing source owner")
	ErrMissingTargetOwner = errors.New("missing target owner")
	ErrMissingSourceRepo  = errors.New("missing source repository")
	ErrMissingTargetRepo  = errors.New("missing target repository")
	ErrMissingSourceOrg   = errors.New("missing source organization")
	ErrMissingTargetOrg   = errors.New("missing target organization")
)

// RateLimitInfo holds rate limit information from the GitHub API
type RateLimitInfo struct {
	Limit     int
	Remaining int
	ResetTime time.Time
}

// Variable represents a GitHub Actions variable
type Variable struct {
	Name      string `json:"name"`
	Value     string `json:"value"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

// Environment represents a GitHub repository environment
type Environment struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

// MigrationMode defines the type of migration to perform
type MigrationMode string

const (
	ModeRepoToRepo MigrationMode = "repo-to-repo"
	ModeOrgToOrg   MigrationMode = "org-to-org"
)

// MigrationConfig holds the configuration for a migration
type MigrationConfig struct {
	Mode MigrationMode

	// Source
	SourceOwner string
	SourceRepo  string
	SourceOrg   string

	// Target
	TargetOwner string
	TargetRepo  string
	TargetOrg   string

	// Environment variables settings
	SkipEnvs bool

	// Options
	DryRun bool
	Force  bool
}

// MigrationResult holds the result of a migration
type MigrationResult struct {
	Created int
	Updated int
	Skipped int
	Errors  []error
}

// AddError adds an error to the result
func (r *MigrationResult) AddError(err error) {
	r.Errors = append(r.Errors, err)
}

// HasErrors returns true if there are any errors
func (r *MigrationResult) HasErrors() bool {
	return len(r.Errors) > 0
}

// Total returns the total number of variables processed
func (r *MigrationResult) Total() int {
	return r.Created + r.Updated + r.Skipped
}
