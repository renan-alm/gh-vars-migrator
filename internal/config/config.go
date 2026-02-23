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
	if cfg.OrgVisibility != "" {
		switch cfg.OrgVisibility {
		case "all", "private", "selected":
			// valid
		default:
			return fmt.Errorf("invalid org visibility %q: must be 'all', 'private', or 'selected'", cfg.OrgVisibility)
		}
	}
	return nil
}

// GetDescription returns a human-readable description of the migration
func GetDescription(cfg *types.MigrationConfig) string {
	switch cfg.Mode {
	case types.ModeRepoToRepo:
		desc := fmt.Sprintf("Repository %s/%s → %s/%s",
			cfg.SourceOwner, cfg.SourceRepo,
			cfg.TargetOwner, cfg.TargetRepo)
		if !cfg.SkipEnvs {
			desc += " (with environments)"
		}
		return desc
	case types.ModeOrgToOrg:
		return fmt.Sprintf("Organization %s → %s",
			cfg.SourceOrg, cfg.TargetOrg)
	default:
		return "Unknown migration"
	}
}
