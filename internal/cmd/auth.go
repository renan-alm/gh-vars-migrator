package cmd

import (
	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/renan-alm/gh-vars-migrator/internal/logger"
	"github.com/spf13/cobra"
)

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Check GitHub CLI authentication status",
	Long:  `Verify that you are properly authenticated with the GitHub CLI and have access to the required organizations.`,
	Example: `  # Check authentication status
  gh vars-migrator auth

  # Check access to specific organizations
  gh vars-migrator auth --check-org renan-org --check-org demo-org-renan`,
	RunE: runAuthCheck,
}

var checkOrgs []string

func init() {
	rootCmd.AddCommand(authCmd)
	authCmd.Flags().StringSliceVar(&checkOrgs, "check-org", []string{}, "Organization(s) to check access for")
}

func runAuthCheck(cmd *cobra.Command, args []string) error {
	logger.Info("Checking GitHub CLI authentication...")
	logger.Plain("")

	// Check basic authentication
	client, err := api.DefaultRESTClient()
	if err != nil {
		logger.Error("Failed to create GitHub API client: %v", err)
		logger.Plain("\nTo authenticate, run: gh auth login")
		return err
	}

	var user struct {
		Login string `json:"login"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := client.Get("user", &user); err != nil {
		logger.Error("Authentication failed: %v", err)
		logger.Plain("\nTo authenticate, run: gh auth login")
		return err
	}

	logger.Success("✓ Authenticated successfully")
	logger.Plain("  User:  %s", user.Login)
	if user.Name != "" {
		logger.Plain("  Name:  %s", user.Name)
	}
	if user.Email != "" {
		logger.Plain("  Email: %s", user.Email)
	}

	// Check organization access if specified
	if len(checkOrgs) > 0 {
		logger.Plain("")
		logger.Info("Checking organization access...")

		allOK := true
		for _, org := range checkOrgs {
			if err := CheckOrgAccess(org); err != nil {
				logger.Error("✗ Cannot access organization '%s': %v", org, err)
				allOK = false
			} else {
				logger.Success("✓ Organization '%s' accessible", org)
			}
		}

		if !allOK {
			logger.Plain("")
			logger.Warning("Some organizations are not accessible. Ensure you have the required permissions.")
			return nil
		}
	}

	logger.Plain("")
	logger.Success("Authentication check passed!")
	return nil
}
