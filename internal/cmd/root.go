package cmd

import (
	"os"

	"github.com/renan-alm/gh-vars-migrator/internal/logger"
	"github.com/spf13/cobra"
)

var (
	// Version is set at build time
	Version = "dev"
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "gh-vars-migrator",
	Short: "GitHub CLI extension for migrating Actions variables between organizations, repositories, and environments",
	Long: `gh-vars-migrator is a GitHub CLI extension that helps you migrate 
GitHub Actions variables between organizations, repositories, and environments.

It supports:
  • Organization to organization variable migration
  • Repository to repository variable migration
  • Environment to environment variable migration
  • Dry-run mode to preview changes before applying
  • Force mode to overwrite existing variables

Examples:
  # Organization to Organization migration
  gh vars-migrator migrate --source-org myorg --target-org targetorg --org-to-org

  # Repository to Repository migration
  gh vars-migrator migrate --source-org myorg --source-repo myrepo --target-org targetorg --target-repo targetrepo

  # Environment to Environment migration
  gh vars-migrator migrate --source-org myorg --source-repo myrepo --source-env staging --target-env production

  # Dry-run mode (preview changes)
  gh vars-migrator migrate --source-org myorg --target-org targetorg --org-to-org --dry-run

  # Force overwrite existing variables
  gh vars-migrator migrate --source-org myorg --target-org targetorg --org-to-org --force`,
	Version: Version,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logger.Error("%v", err)
		os.Exit(1)
	}
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output")
}
