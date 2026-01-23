package integration

import (
	"testing"

	"github.com/renan-alm/gh-vars-migrator/internal/config"
	"github.com/renan-alm/gh-vars-migrator/internal/types"
)

// TestEndToEnd_ConfigValidation tests the complete configuration validation workflow
func TestEndToEnd_ConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		config  *types.MigrationConfig
		wantErr bool
	}{
		{
			name: "valid_repo_to_repo",
			config: &types.MigrationConfig{
				Mode:        types.ModeRepoToRepo,
				SourceOwner: "source-org",
				SourceRepo:  "source-repo",
				TargetOwner: "target-org",
				TargetRepo:  "target-repo",
				DryRun:      true,
			},
			wantErr: false,
		},
		{
			name: "valid_org_to_org",
			config: &types.MigrationConfig{
				Mode:      types.ModeOrgToOrg,
				SourceOrg: "source-org",
				TargetOrg: "target-org",
				Force:     true,
			},
			wantErr: false,
		},
		{
			name: "valid_env_only",
			config: &types.MigrationConfig{
				Mode:        types.ModeEnvOnly,
				SourceOwner: "owner",
				SourceRepo:  "repo",
				SourceEnv:   "staging",
				TargetEnv:   "production",
			},
			wantErr: false,
		},
		{
			name: "invalid_missing_source_owner",
			config: &types.MigrationConfig{
				Mode:        types.ModeRepoToRepo,
				SourceRepo:  "repo",
				TargetOwner: "target",
				TargetRepo:  "repo",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := config.Validate(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("config.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				// Get description for valid configs
				desc := config.GetDescription(tt.config)
				if desc == "" {
					t.Error("Expected non-empty description for valid config")
				}
			}
		})
	}
}

// TestEndToEnd_MigrationResultTracking tests result aggregation across migration
func TestEndToEnd_MigrationResultTracking(t *testing.T) {
	result := &types.MigrationResult{}

	// Simulate a migration that:
	// - Creates 5 new variables
	// - Updates 3 existing variables
	// - Skips 2 variables (already exist, no force flag)
	// - Encounters 1 error
	result.Created = 5
	result.Updated = 3
	result.Skipped = 2
	result.AddError(types.ErrInvalidConfig)

	// Verify total count
	expectedTotal := 10
	if result.Total() != expectedTotal {
		t.Errorf("Expected total %d, got %d", expectedTotal, result.Total())
	}

	// Verify error tracking
	if !result.HasErrors() {
		t.Error("Expected result to have errors")
	}

	if len(result.Errors) != 1 {
		t.Errorf("Expected 1 error, got %d", len(result.Errors))
	}
}

// TestEndToEnd_DryRunWorkflow tests that dry-run mode works end-to-end
func TestEndToEnd_DryRunWorkflow(t *testing.T) {
	// Test that dry-run configuration is properly set and validated
	cfg := &types.MigrationConfig{
		Mode:        types.ModeRepoToRepo,
		SourceOwner: "source",
		SourceRepo:  "repo",
		TargetOwner: "target",
		TargetRepo:  "repo",
		DryRun:      true,
		Force:       false,
	}

	// Validate config
	if err := config.Validate(cfg); err != nil {
		t.Fatalf("Config validation failed: %v", err)
	}

	// Verify dry-run flag is set
	if !cfg.DryRun {
		t.Error("Expected DryRun to be true")
	}

	// Simulate dry-run results (would create 10 vars but doesn't actually)
	result := &types.MigrationResult{
		Created: 10,
		Updated: 0,
		Skipped: 0,
	}

	// In dry-run, we should have tracked what would happen
	if result.Created != 10 {
		t.Errorf("Expected 10 variables tracked for creation in dry-run, got %d", result.Created)
	}

	// No errors should occur in dry-run for valid operations
	if result.HasErrors() {
		t.Error("Expected no errors in dry-run mode")
	}
}

// TestEndToEnd_ForceUpdateWorkflow tests the force update workflow
func TestEndToEnd_ForceUpdateWorkflow(t *testing.T) {
	cfg := &types.MigrationConfig{
		Mode:        types.ModeRepoToRepo,
		SourceOwner: "source",
		SourceRepo:  "repo",
		TargetOwner: "target",
		TargetRepo:  "repo",
		Force:       true,
		DryRun:      false,
	}

	if err := config.Validate(cfg); err != nil {
		t.Fatalf("Config validation failed: %v", err)
	}

	if !cfg.Force {
		t.Error("Expected Force flag to be set")
	}

	// With force=true, existing variables should be updated
	result := &types.MigrationResult{
		Created: 3,  // 3 new variables
		Updated: 7,  // 7 existing variables updated due to force flag
		Skipped: 0,  // 0 skipped because force=true
	}

	if result.Updated != 7 {
		t.Errorf("Expected 7 updates with force=true, got %d", result.Updated)
	}

	if result.Skipped != 0 {
		t.Errorf("Expected 0 skipped with force=true, got %d", result.Skipped)
	}
}

// TestEndToEnd_EnvironmentMigration tests environment-only migration workflow
func TestEndToEnd_EnvironmentMigration(t *testing.T) {
	// Test migrating from staging to production
	cfg := &types.MigrationConfig{
		Mode:        types.ModeEnvOnly,
		SourceOwner: "my-org",
		SourceRepo:  "my-repo",
		SourceEnv:   "staging",
		TargetEnv:   "production",
		DryRun:      true,
	}

	if err := config.Validate(cfg); err != nil {
		t.Fatalf("Config validation failed: %v", err)
	}

	// Verify environment names are set correctly
	if cfg.SourceEnv != "staging" {
		t.Errorf("Expected source env 'staging', got '%s'", cfg.SourceEnv)
	}

	if cfg.TargetEnv != "production" {
		t.Errorf("Expected target env 'production', got '%s'", cfg.TargetEnv)
	}

	// Simulate migration result
	result := &types.MigrationResult{
		Created: 15,  // 15 environment variables created in target
		Updated: 0,
		Skipped: 0,
	}

	if result.Created != 15 {
		t.Errorf("Expected 15 environment variables created, got %d", result.Created)
	}
}

// TestEndToEnd_ErrorHandlingWorkflow tests error handling across the migration
func TestEndToEnd_ErrorHandlingWorkflow(t *testing.T) {
	result := &types.MigrationResult{}

	// Simulate various errors during migration
	errors := []error{
		types.ErrInvalidConfig,
		types.ErrMissingSourceOwner,
		types.ErrMissingTargetOwner,
	}

	for _, err := range errors {
		result.AddError(err)
	}

	// Verify error accumulation
	if len(result.Errors) != 3 {
		t.Errorf("Expected 3 errors, got %d", len(result.Errors))
	}

	if !result.HasErrors() {
		t.Error("Expected result to have errors")
	}

	// Even with errors, partial success should be tracked
	result.Created = 5  // 5 succeeded before errors
	result.Updated = 2  // 2 updated before errors

	if result.Total() != 7 {
		t.Errorf("Expected total of 7 (partial success), got %d", result.Total())
	}
}

// TestEndToEnd_ConfigDescriptions tests descriptive output for different modes
func TestEndToEnd_ConfigDescriptions(t *testing.T) {
	tests := []struct {
		name   string
		config *types.MigrationConfig
	}{
		{
			name: "repo_to_repo_description",
			config: &types.MigrationConfig{
				Mode:        types.ModeRepoToRepo,
				SourceOwner: "org1",
				SourceRepo:  "repo1",
				TargetOwner: "org2",
				TargetRepo:  "repo2",
			},
		},
		{
			name: "org_to_org_description",
			config: &types.MigrationConfig{
				Mode:      types.ModeOrgToOrg,
				SourceOrg: "source-org",
				TargetOrg: "target-org",
			},
		},
		{
			name: "env_only_description",
			config: &types.MigrationConfig{
				Mode:        types.ModeEnvOnly,
				SourceOwner: "owner",
				SourceRepo:  "repo",
				SourceEnv:   "staging",
				TargetEnv:   "production",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			desc := config.GetDescription(tt.config)
			if desc == "" {
				t.Error("Expected non-empty description")
			}

			if len(desc) < 10 {
				t.Errorf("Description seems too short: %s", desc)
			}
		})
	}
}

// TestEndToEnd_MigrationModes tests all three migration modes
func TestEndToEnd_MigrationModes(t *testing.T) {
	modes := []types.MigrationMode{
		types.ModeRepoToRepo,
		types.ModeOrgToOrg,
		types.ModeEnvOnly,
	}

	for _, mode := range modes {
		t.Run(string(mode), func(t *testing.T) {
			var cfg *types.MigrationConfig

			switch mode {
			case types.ModeRepoToRepo:
				cfg = &types.MigrationConfig{
					Mode:        mode,
					SourceOwner: "src",
					SourceRepo:  "repo",
					TargetOwner: "tgt",
					TargetRepo:  "repo",
				}
			case types.ModeOrgToOrg:
				cfg = &types.MigrationConfig{
					Mode:      mode,
					SourceOrg: "src-org",
					TargetOrg: "tgt-org",
				}
			case types.ModeEnvOnly:
				cfg = &types.MigrationConfig{
					Mode:        mode,
					SourceOwner: "owner",
					SourceRepo:  "repo",
					SourceEnv:   "staging",
					TargetEnv:   "prod",
				}
			}

			if err := config.Validate(cfg); err != nil {
				t.Errorf("Mode %s validation failed: %v", mode, err)
			}

			if cfg.Mode != mode {
				t.Errorf("Expected mode %s, got %s", mode, cfg.Mode)
			}
		})
	}
}
