package cmd

import (
	"fmt"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/renan-alm/gh-vars-migrator/internal/logger"
	"github.com/renan-alm/gh-vars-migrator/internal/types"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List variables in an organization",
	Long:  `List all GitHub Actions variables in the specified organization.`,
	Example: `  # List variables in an organization
  gh vars-migrator list --org renan-org`,
	RunE: runList,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if listOrg == "" {
			return fmt.Errorf("--org flag is required")
		}
		cmd.SilenceUsage = true
		return nil
	},
}

var listOrg string

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(&listOrg, "org", "o", "", "Organization name (required)")
	_ = listCmd.MarkFlagRequired("org")
}

func runList(cmd *cobra.Command, args []string) error {
	// Check authentication first
	if err := checkAuth(); err != nil {
		return err
	}

	logger.Info("Listing variables for organization: %s", listOrg)
	logger.Plain("")

	client, err := api.DefaultRESTClient()
	if err != nil {
		return fmt.Errorf("failed to create GitHub API client: %w", err)
	}

	var response struct {
		TotalCount int              `json:"total_count"`
		Variables  []types.Variable `json:"variables"`
	}

	path := fmt.Sprintf("orgs/%s/actions/variables", listOrg)
	if err := client.Get(path, &response); err != nil {
		return fmt.Errorf("failed to list variables: %w", err)
	}

	if len(response.Variables) == 0 {
		logger.Warning("No variables found in organization '%s'", listOrg)
		return nil
	}

	logger.Info("Found %d variable(s):", len(response.Variables))
	logger.Plain("")
	logger.Plain("%-30s %s", "NAME", "UPDATED AT")
	logger.Plain("%-30s %s", "----", "----------")

	for _, v := range response.Variables {
		logger.Plain("%-30s %s", v.Name, v.UpdatedAt)
	}

	logger.Plain("")
	logger.Success("Total: %d variable(s)", len(response.Variables))
	return nil
}
