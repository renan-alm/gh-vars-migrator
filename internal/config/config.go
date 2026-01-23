package config

import (
	"errors"
	"fmt"

	"github.com/renan-alm/gh-vars-migrator/internal/types"
)

// Validate validates the migration configuration
func Validate(cfg *types.MigrationConfig) error {
	if cfg == nil {
		return errors.New("configuration is nil")
	}

	switch cfg.Mode {
	case types.ModeRepoToRepo:
		return validateRepoToRepo(cfg)
	case types.ModeOrgToOrg:
		return validateOrgToOrg(cfg)
	case types.ModeEnvOnly:
		return validateEnvOnly(cfg)
	default:
		return fmt.Errorf("invalid migration mode: %s", cfg.Mode)
	}
}

// validateRepoToRepo validates repository to repository migration configuration
func validateRepoToRepo(cfg *types.MigrationConfig) error {
	if cfg.SourceOwner == "" {
		return errors.New("source owner is required")
	}
	if cfg.SourceRepo == "" {
		return errors.New("source repository is required")
	}
	if cfg.TargetOwner == "" {
		return errors.New("target owner is required")
	}
	if cfg.TargetRepo == "" {
		return errors.New("target repository is required")
	}
	return nil
}

// validateOrgToOrg validates organization to organization migration configuration
func validateOrgToOrg(cfg *types.MigrationConfig) error {
	if cfg.SourceOrg == "" {
		return errors.New("source organization is required")
	}
	if cfg.TargetOrg == "" {
		return errors.New("target organization is required")
	}
	return nil
}

// validateEnvOnly validates environment-only migration configuration
func validateEnvOnly(cfg *types.MigrationConfig) error {
	if cfg.SourceOwner == "" {
		return errors.New("source owner is required")
	}
	if cfg.SourceRepo == "" {
		return errors.New("source repository is required")
	}
	if cfg.SourceEnv == "" {
		return errors.New("source environment is required")
	}
	if cfg.TargetEnv == "" {
		return errors.New("target environment is required")
	}
	// Note: TargetOwner and TargetRepo default to Source if not provided
	return nil
}

// GetDescription returns a human-readable description of the migration
func GetDescription(cfg *types.MigrationConfig) string {
	switch cfg.Mode {
	case types.ModeRepoToRepo:
		return fmt.Sprintf("Repository %s/%s → %s/%s",
			cfg.SourceOwner, cfg.SourceRepo,
			cfg.TargetOwner, cfg.TargetRepo)
	case types.ModeOrgToOrg:
		return fmt.Sprintf("Organization %s → %s",
			cfg.SourceOrg, cfg.TargetOrg)
	case types.ModeEnvOnly:
		return fmt.Sprintf("Environment %s → %s (Repository: %s/%s)",
			cfg.SourceEnv, cfg.TargetEnv,
			cfg.SourceOwner, cfg.SourceRepo)
	default:
		return "Unknown migration"
	}
}
