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
	Short: "GitHub CLI extension for migrating Actions variables between organizations",
	Long: `gh-vars-migrator is a GitHub CLI extension that helps you migrate 
GitHub Actions variables from one organization to another.

It supports:
  • Organization to organization variable migration
  • Dry-run mode to preview changes before applying
  • Force mode to overwrite existing variables

Examples:
  # Dry-run migration (preview only)
  gh vars-migrator org --source my-source-org --target my-target-org --dry-run

  # Perform actual migration
  gh vars-migrator org --source my-source-org --target my-target-org

  # Force overwrite existing variables
  gh vars-migrator org --source my-source-org --target my-target-org --force`,
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
