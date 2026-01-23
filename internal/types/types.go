package types

// Variable represents a GitHub Actions variable
type Variable struct {
	Name      string `json:"name"`
	Value     string `json:"value"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

// MigrationMode defines the type of migration to perform
type MigrationMode string

const (
	ModeRepoToRepo MigrationMode = "repo-to-repo"
	ModeOrgToOrg   MigrationMode = "org-to-org"
	ModeEnvOnly    MigrationMode = "env-only"
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
	SourceEnv string
	TargetEnv string
	SkipEnvs  bool

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
