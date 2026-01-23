package cmd

import (
	"fmt"

	"github.com/renan-alm/gh-vars-migrator/internal/logger"
	"github.com/renan-alm/gh-vars-migrator/internal/migrator"
	"github.com/renan-alm/gh-vars-migrator/internal/types"
	"github.com/spf13/cobra"
)

var (
	orgSourceOrg string
	orgTargetOrg string
	orgDryRun    bool
	orgForce     bool
)

// orgCmd represents the org command for org-to-org migration
// Deprecated: Use root command flags with --org-to-org instead
var orgCmd = &cobra.Command{
	Use:        "org",
	Short:      "Migrate variables from one organization to another (deprecated: use root command flags)",
	Deprecated: "use root command flags directly instead (e.g., 'gh vars-migrator --source-org ... --target-org ... --org-to-org')",
	Long: `Migrate all GitHub Actions variables from a source organization 
to a target organization.

DEPRECATED: This command is deprecated. Please use the root command flags directly:
  gh vars-migrator --source-org SOURCE --target-org TARGET --org-to-org

This command will:
  1. Fetch all variables from the source organization
  2. Check if each variable exists in the target organization
  3. Create or update variables in the target organization

Use --dry-run to preview changes without applying them.
Use --force to overwrite existing variables in the target organization.`,
	Example: `  # OLD - Deprecated
  gh vars-migrator org --source renan-org --target demo-org-renan

  # NEW - Recommended
  gh vars-migrator --source-org renan-org --target-org demo-org-renan --org-to-org`,
	RunE: runOrgMigration,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// Validate required flags
		if orgSourceOrg == "" {
			return fmt.Errorf("--source flag is required")
		}
		if orgTargetOrg == "" {
			return fmt.Errorf("--target flag is required")
		}
		if orgSourceOrg == orgTargetOrg {
			return fmt.Errorf("source and target organizations cannot be the same")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(orgCmd)

	// Required flags
	orgCmd.Flags().StringVarP(&orgSourceOrg, "source", "s", "", "Source organization name (required)")
	orgCmd.Flags().StringVarP(&orgTargetOrg, "target", "t", "", "Target organization name (required)")

	// Optional flags
	orgCmd.Flags().BoolVarP(&orgDryRun, "dry-run", "d", false, "Preview changes without applying them")
	orgCmd.Flags().BoolVarP(&orgForce, "force", "f", false, "Overwrite existing variables in target")

	// Mark required flags
	_ = orgCmd.MarkFlagRequired("source")
	_ = orgCmd.MarkFlagRequired("target")
}

func runOrgMigration(cmd *cobra.Command, args []string) error {
	// Check authentication first
	if err := checkAuth(); err != nil {
		return err
	}

	logger.Info("gh-vars-migrator - Organization Variable Migration")
	logger.Info("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	logger.Info("Source: %s", orgSourceOrg)
	logger.Info("Target: %s", orgTargetOrg)
	logger.Info("Dry-run: %v", orgDryRun)
	logger.Info("Force: %v", orgForce)
	logger.Info("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Create migration configuration
	cfg := &types.MigrationConfig{
		Mode:      types.ModeOrgToOrg,
		SourceOrg: orgSourceOrg,
		TargetOrg: orgTargetOrg,
		DryRun:    orgDryRun,
		Force:     orgForce,
	}

	// Create and run migrator
	m, err := migrator.New(cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize migrator: %w", err)
	}

	result, err := m.Run()
	if err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	if result.HasErrors() {
		return fmt.Errorf("migration completed with %d error(s)", len(result.Errors))
	}

	logger.Success("Migration completed successfully!")
	return nil
}
