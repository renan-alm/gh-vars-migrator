package cmd

import (
	"fmt"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/renan-alm/gh-vars-migrator/internal/logger"
	"github.com/renan-alm/gh-vars-migrator/internal/migrator"
	"github.com/renan-alm/gh-vars-migrator/internal/types"
	"github.com/spf13/cobra"
)

var (
	// Source flags
	migrateSourceOrg  string
	migrateSourceRepo string
	migrateSourceEnv  string

	// Target flags
	migrateTargetOrg  string
	migrateTargetRepo string
	migrateTargetEnv  string

	// Mode flags
	migrateOrgToOrg bool
	migrateSkipEnvs bool

	// Option flags
	migrateDryRun bool
	migrateForce  bool
)

// migrateCmd represents the unified migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate GitHub Actions variables between organizations, repositories, or environments",
	Long: `Migrate GitHub Actions variables with support for multiple migration modes:

• Organization to Organization: Migrate org-level variables (--org-to-org flag)
• Repository to Repository: Migrate repo-level variables and optionally environment variables
• Environment to Environment: Migrate environment variables within same or different repos

Mode Detection:
  - If --org-to-org flag is set → Organization migration mode
  - If --source-env and --target-env are provided → Environment-only migration mode
  - Otherwise → Repository-to-Repository migration mode

Use --dry-run to preview changes without applying them.
Use --force to overwrite existing variables in the target.`,
	Example: `  # Organization to Organization migration
  gh vars-migrator migrate --source-org myorg --target-org targetorg --org-to-org

  # Repository to Repository migration
  gh vars-migrator migrate --source-org myorg --source-repo myrepo --target-org targetorg --target-repo targetrepo

  # Repository to Repository with environment variables
  gh vars-migrator migrate --source-org myorg --source-repo myrepo --target-org targetorg --target-repo targetrepo --source-env production --target-env production

  # Skip environment migration
  gh vars-migrator migrate --source-org myorg --source-repo myrepo --target-org targetorg --target-repo targetrepo --skip-envs

  # Environment only migration (same repo, different environments)
  gh vars-migrator migrate --source-org myorg --source-repo myrepo --target-org myorg --source-env staging --target-env production

  # Dry run mode
  gh vars-migrator migrate --source-org myorg --target-org targetorg --org-to-org --dry-run

  # Force overwrite
  gh vars-migrator migrate --source-org myorg --target-org targetorg --org-to-org --force`,
	PreRunE: validateMigrateFlags,
	RunE:    runMigrate,
}

func init() {
	rootCmd.AddCommand(migrateCmd)

	// Source flags
	migrateCmd.Flags().StringVar(&migrateSourceOrg, "source-org", "", "Source organization name (required)")
	migrateCmd.Flags().StringVar(&migrateSourceRepo, "source-repo", "", "Source repository name (required for repo-to-repo and env migrations)")
	migrateCmd.Flags().StringVar(&migrateSourceEnv, "source-env", "", "Source environment name (for environment migrations)")

	// Target flags
	migrateCmd.Flags().StringVar(&migrateTargetOrg, "target-org", "", "Target organization name (required)")
	migrateCmd.Flags().StringVar(&migrateTargetRepo, "target-repo", "", "Target repository name (required for repo-to-repo, optional for org-to-org)")
	migrateCmd.Flags().StringVar(&migrateTargetEnv, "target-env", "", "Target environment name (for environment migrations)")

	// Mode flags
	migrateCmd.Flags().BoolVar(&migrateOrgToOrg, "org-to-org", false, "Migrate organization variables only")
	migrateCmd.Flags().BoolVar(&migrateSkipEnvs, "skip-envs", false, "Skip environment variable migration during repo-to-repo")

	// Option flags
	migrateCmd.Flags().BoolVar(&migrateDryRun, "dry-run", false, "Preview changes without applying them")
	migrateCmd.Flags().BoolVar(&migrateForce, "force", false, "Overwrite existing variables in target")

	// Mark required flags
	// These should never fail as the flags are defined above
	_ = migrateCmd.MarkFlagRequired("source-org")
	_ = migrateCmd.MarkFlagRequired("target-org")
}

// validateMigrateFlags validates the flags based on the detected migration mode
func validateMigrateFlags(cmd *cobra.Command, args []string) error {
	// Suppress usage on runtime errors
	cmd.SilenceUsage = true

	// Check for conflicting flags
	if migrateOrgToOrg && (migrateSourceEnv != "" || migrateTargetEnv != "") {
		return fmt.Errorf("cannot use --org-to-org with environment flags (--source-env, --target-env)")
	}

	// Detect mode and validate accordingly
	mode := detectMigrationMode()

	switch mode {
	case types.ModeOrgToOrg:
		// Org-to-org: no additional requirements
		if migrateSourceOrg == migrateTargetOrg {
			return fmt.Errorf("source and target organizations cannot be the same")
		}

	case types.ModeEnvOnly:
		// Environment-only: requires source repo, source env, and target env
		if migrateSourceRepo == "" {
			return fmt.Errorf("--source-repo is required for environment migration")
		}
		if migrateSourceEnv == "" {
			return fmt.Errorf("--source-env is required for environment migration")
		}
		if migrateTargetEnv == "" {
			return fmt.Errorf("--target-env is required for environment migration")
		}
		// Target repo defaults to source repo if not specified

	case types.ModeRepoToRepo:
		// Repo-to-repo: requires source repo and target repo
		if migrateSourceRepo == "" {
			return fmt.Errorf("--source-repo is required for repository migration")
		}
		if migrateTargetRepo == "" {
			return fmt.Errorf("--target-repo is required for repository migration")
		}
		if migrateSourceOrg == migrateTargetOrg && migrateSourceRepo == migrateTargetRepo {
			return fmt.Errorf("source and target repositories cannot be the same")
		}
	}

	return nil
}

// detectMigrationMode determines the migration mode based on the provided flags
func detectMigrationMode() types.MigrationMode {
	// If --org-to-org flag is set, it's organization migration
	if migrateOrgToOrg {
		return types.ModeOrgToOrg
	}

	// If both source-env and target-env are provided, it's environment-only migration
	if migrateSourceEnv != "" && migrateTargetEnv != "" {
		return types.ModeEnvOnly
	}

	// Default to repository-to-repository migration
	return types.ModeRepoToRepo
}

// runMigrate executes the migration based on the detected mode
func runMigrate(cmd *cobra.Command, args []string) error {
	// Check authentication first
	if err := checkMigrateAuth(); err != nil {
		return err
	}

	// Detect migration mode
	mode := detectMigrationMode()

	// Build migration configuration
	cfg := &types.MigrationConfig{
		Mode:      mode,
		SourceOrg: migrateSourceOrg,
		TargetOrg: migrateTargetOrg,
		DryRun:    migrateDryRun,
		Force:     migrateForce,
	}

	// Set mode-specific configuration
	switch mode {
	case types.ModeOrgToOrg:
		logger.Info("gh-vars-migrator - Organization Variable Migration")
		logger.Info("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		logger.Info("Source: %s", migrateSourceOrg)
		logger.Info("Target: %s", migrateTargetOrg)

	case types.ModeEnvOnly:
		cfg.SourceOwner = migrateSourceOrg
		cfg.SourceRepo = migrateSourceRepo
		cfg.SourceEnv = migrateSourceEnv
		cfg.TargetOwner = migrateTargetOrg
		if migrateTargetRepo != "" {
			cfg.TargetRepo = migrateTargetRepo
		} else {
			cfg.TargetRepo = migrateSourceRepo // Default to source repo
		}
		cfg.TargetEnv = migrateTargetEnv

		logger.Info("gh-vars-migrator - Environment Variable Migration")
		logger.Info("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		logger.Info("Source: %s/%s (env: %s)", cfg.SourceOwner, cfg.SourceRepo, cfg.SourceEnv)
		logger.Info("Target: %s/%s (env: %s)", cfg.TargetOwner, cfg.TargetRepo, cfg.TargetEnv)

	case types.ModeRepoToRepo:
		cfg.SourceOwner = migrateSourceOrg
		cfg.SourceRepo = migrateSourceRepo
		cfg.TargetOwner = migrateTargetOrg
		cfg.TargetRepo = migrateTargetRepo
		cfg.SkipEnvs = migrateSkipEnvs
		// Set environment variables if provided (for optional env migration)
		if migrateSourceEnv != "" {
			cfg.SourceEnv = migrateSourceEnv
		}
		if migrateTargetEnv != "" {
			cfg.TargetEnv = migrateTargetEnv
		}

		logger.Info("gh-vars-migrator - Repository Variable Migration")
		logger.Info("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		logger.Info("Source: %s/%s", cfg.SourceOwner, cfg.SourceRepo)
		logger.Info("Target: %s/%s", cfg.TargetOwner, cfg.TargetRepo)
		if migrateSkipEnvs {
			logger.Info("Skip Environments: true")
		}
	}

	// Common configuration display
	logger.Info("Dry-run: %v", migrateDryRun)
	logger.Info("Force: %v", migrateForce)
	logger.Info("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

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

// checkMigrateAuth verifies that the user is authenticated with GitHub CLI
func checkMigrateAuth() error {
	client, err := api.DefaultRESTClient()
	if err != nil {
		return fmt.Errorf("failed to create GitHub API client: %w\n\nPlease authenticate using: gh auth login", err)
	}

	var user struct {
		Login string `json:"login"`
	}

	if err := client.Get("user", &user); err != nil {
		return fmt.Errorf("authentication failed: %w\n\nPlease authenticate using: gh auth login", err)
	}

	logger.Success("Authenticated as: %s", user.Login)
	return nil
}
