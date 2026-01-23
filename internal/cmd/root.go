package cmd

import (
	"fmt"
	"os"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/renan-alm/gh-vars-migrator/internal/logger"
	"github.com/renan-alm/gh-vars-migrator/internal/migrator"
	"github.com/renan-alm/gh-vars-migrator/internal/types"
	"github.com/spf13/cobra"
)

var (
	// Version is set at build time
	Version = "dev"

	// Source flags
	sourceOrg  string
	sourceRepo string

	// Target flags
	targetOrg  string
	targetRepo string

	// Mode flags
	orgToOrg bool
	skipEnvs bool

	// Option flags
	dryRun bool
	force  bool
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "gh-vars-migrator",
	Short: "GitHub CLI extension for migrating Actions variables between organizations, repositories, and environments",
	Long: `gh-vars-migrator is a GitHub CLI extension that helps you migrate 
GitHub Actions variables between organizations, repositories, and environments.

It supports:
  • Organization to organization variable migration
  • Repository to repository variable migration (with auto-discovery of environments)
  • Dry-run mode to preview changes before applying
  • Force mode to overwrite existing variables

Mode Detection:
  - If --org-to-org flag is set → Organization migration mode
  - Otherwise → Repository-to-Repository migration mode (includes all environments)`,
	Example: `  # Organization to Organization migration
  gh vars-migrator --source-org myorg --target-org targetorg --org-to-org

  # Repository to Repository migration (auto-discovers and migrates all environments)
  gh vars-migrator --source-org myorg --source-repo myrepo --target-org targetorg --target-repo targetrepo

  # Repository migration without environments
  gh vars-migrator --source-org myorg --source-repo myrepo --target-org targetorg --target-repo targetrepo --skip-envs

  # Dry-run mode (preview changes)
  gh vars-migrator --source-org myorg --target-org targetorg --org-to-org --dry-run

  # Force overwrite existing variables
  gh vars-migrator --source-org myorg --target-org targetorg --org-to-org --force

  # Utility commands
  gh vars-migrator auth
  gh vars-migrator list --org myorg`,
	Version: Version,
	PreRunE: validateFlags,
	RunE:    runMigration,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logger.Error("%v", err)
		os.Exit(1)
	}
}

func init() {
	// Source flags
	rootCmd.Flags().StringVar(&sourceOrg, "source-org", "", "Source organization name (required)")
	rootCmd.Flags().StringVar(&sourceRepo, "source-repo", "", "Source repository name (required for repo-to-repo)")

	// Target flags
	rootCmd.Flags().StringVar(&targetOrg, "target-org", "", "Target organization name (required)")
	rootCmd.Flags().StringVar(&targetRepo, "target-repo", "", "Target repository name (required for repo-to-repo)")

	// Mode flags
	rootCmd.Flags().BoolVar(&orgToOrg, "org-to-org", false, "Migrate organization variables only")
	rootCmd.Flags().BoolVar(&skipEnvs, "skip-envs", false, "Skip environment variable migration during repo-to-repo")

	// Option flags
	rootCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Preview changes without applying them")
	rootCmd.Flags().BoolVar(&force, "force", false, "Overwrite existing variables in target")

	// Global flags
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output")
}

// validateFlags validates the flags based on the detected migration mode
func validateFlags(cmd *cobra.Command, args []string) error {
	// If a subcommand is being run, skip validation
	if cmd.Name() != "gh-vars-migrator" {
		return nil
	}

	// Check if any migration flags were provided
	if sourceOrg == "" && targetOrg == "" {
		// No flags provided, show help
		return cmd.Help()
	}

	// Suppress usage on runtime errors
	cmd.SilenceUsage = true

	// Validate required flags
	if sourceOrg == "" {
		return fmt.Errorf("--source-org flag is required")
	}
	if targetOrg == "" {
		return fmt.Errorf("--target-org flag is required")
	}

	// Detect mode and validate accordingly
	mode := detectMigrationMode()

	switch mode {
	case types.ModeOrgToOrg:
		// Org-to-org: no additional requirements
		if sourceOrg == targetOrg {
			return fmt.Errorf("source and target organizations cannot be the same")
		}

	case types.ModeRepoToRepo:
		// Repo-to-repo: requires source repo and target repo
		if sourceRepo == "" {
			return fmt.Errorf("--source-repo is required for repository migration")
		}
		if targetRepo == "" {
			return fmt.Errorf("--target-repo is required for repository migration")
		}
		if sourceOrg == targetOrg && sourceRepo == targetRepo {
			return fmt.Errorf("source and target repositories cannot be the same")
		}
	}

	return nil
}

// detectMigrationMode determines the migration mode based on the provided flags
func detectMigrationMode() types.MigrationMode {
	// If --org-to-org flag is set, it's organization migration
	if orgToOrg {
		return types.ModeOrgToOrg
	}

	// Default to repository-to-repository migration
	return types.ModeRepoToRepo
}

// runMigration executes the migration based on the detected mode
func runMigration(cmd *cobra.Command, args []string) error {
	// Check authentication first
	if err := checkAuth(); err != nil {
		return err
	}

	// Detect migration mode
	mode := detectMigrationMode()

	// Build migration configuration
	cfg := &types.MigrationConfig{
		Mode:      mode,
		SourceOrg: sourceOrg,
		TargetOrg: targetOrg,
		DryRun:    dryRun,
		Force:     force,
	}

	// Set mode-specific configuration
	switch mode {
	case types.ModeOrgToOrg:
		logger.Info("gh-vars-migrator - Organization Variable Migration")
		logger.Info("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		logger.Info("Source: %s", sourceOrg)
		logger.Info("Target: %s", targetOrg)

	case types.ModeRepoToRepo:
		cfg.SourceOwner = sourceOrg
		cfg.SourceRepo = sourceRepo
		cfg.TargetOwner = targetOrg
		cfg.TargetRepo = targetRepo
		cfg.SkipEnvs = skipEnvs

		logger.Info("gh-vars-migrator - Repository Variable Migration")
		logger.Info("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		logger.Info("Source: %s/%s", cfg.SourceOwner, cfg.SourceRepo)
		logger.Info("Target: %s/%s", cfg.TargetOwner, cfg.TargetRepo)
		if skipEnvs {
			logger.Info("Skip Environments: true")
		} else {
			logger.Info("Environments: auto-discover and migrate")
		}
	}

	// Common configuration display
	logger.Info("Dry-run: %v", dryRun)
	logger.Info("Force: %v", force)
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

// checkAuth verifies that the user is authenticated with GitHub CLI
func checkAuth() error {
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

// CheckOrgAccess verifies the user has access to the specified organization
func CheckOrgAccess(orgName string) error {
	client, err := api.DefaultRESTClient()
	if err != nil {
		return err
	}

	var org struct {
		Login string `json:"login"`
	}

	path := fmt.Sprintf("orgs/%s", orgName)
	if err := client.Get(path, &org); err != nil {
		return fmt.Errorf("cannot access organization '%s': %w", orgName, err)
	}

	return nil
}
