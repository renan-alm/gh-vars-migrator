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
	sourceOrg string
	targetOrg string
	dryRun    bool
	force     bool
)

// orgCmd represents the org command for org-to-org migration
var orgCmd = &cobra.Command{
	Use:   "org",
	Short: "Migrate variables from one organization to another",
	Long: `Migrate all GitHub Actions variables from a source organization 
to a target organization.

This command will:
  1. Fetch all variables from the source organization
  2. Check if each variable exists in the target organization
  3. Create or update variables in the target organization

Use --dry-run to preview changes without applying them.
Use --force to overwrite existing variables in the target organization.`,
	Example: `  # Preview migration (dry-run)
  gh vars-migrator org --source renan-org --target demo-org-renan --dry-run

  # Perform migration
  gh vars-migrator org --source renan-org --target demo-org-renan

  # Force overwrite existing variables
  gh vars-migrator org --source renan-org --target demo-org-renan --force`,
	RunE: runOrgMigration,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// Validate required flags
		if sourceOrg == "" {
			return fmt.Errorf("--source flag is required")
		}
		if targetOrg == "" {
			return fmt.Errorf("--target flag is required")
		}
		if sourceOrg == targetOrg {
			return fmt.Errorf("source and target organizations cannot be the same")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(orgCmd)

	// Required flags
	orgCmd.Flags().StringVarP(&sourceOrg, "source", "s", "", "Source organization name (required)")
	orgCmd.Flags().StringVarP(&targetOrg, "target", "t", "", "Target organization name (required)")

	// Optional flags
	orgCmd.Flags().BoolVarP(&dryRun, "dry-run", "d", false, "Preview changes without applying them")
	orgCmd.Flags().BoolVarP(&force, "force", "f", false, "Overwrite existing variables in target")

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
	logger.Info("Source: %s", sourceOrg)
	logger.Info("Target: %s", targetOrg)
	logger.Info("Dry-run: %v", dryRun)
	logger.Info("Force: %v", force)
	logger.Info("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Create migration configuration
	cfg := &types.MigrationConfig{
		Mode:      types.ModeOrgToOrg,
		SourceOrg: sourceOrg,
		TargetOrg: targetOrg,
		DryRun:    dryRun,
		Force:     force,
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

// For testing purposes, allow checking org access before migration
func validateOrgAccess() error {
	logger.Info("Validating organization access...")

	if err := CheckOrgAccess(sourceOrg); err != nil {
		return fmt.Errorf("source organization error: %w", err)
	}
	logger.Success("✓ Source organization '%s' accessible", sourceOrg)

	if err := CheckOrgAccess(targetOrg); err != nil {
		return fmt.Errorf("target organization error: %w", err)
	}
	logger.Success("✓ Target organization '%s' accessible", targetOrg)

	return nil
}

func init() {
	// Set up pre-run validation
	orgCmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		// Validate required flags
		if sourceOrg == "" {
			return fmt.Errorf("--source flag is required")
		}
		if targetOrg == "" {
			return fmt.Errorf("--target flag is required")
		}
		if sourceOrg == targetOrg {
			return fmt.Errorf("source and target organizations cannot be the same")
		}

		// Suppress usage on runtime errors
		cmd.SilenceUsage = true
		return nil
	}
}

// Ensure proper exit code on errors
func init() {
	orgCmd.PostRunE = func(cmd *cobra.Command, args []string) error {
		return nil
	}
}

// Helper function for graceful shutdown
func exitWithError(err error) {
	logger.Error("%v", err)
	os.Exit(1)
}
