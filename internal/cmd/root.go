package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/renan-alm/gh-vars-migrator/internal/client"
	"github.com/renan-alm/gh-vars-migrator/internal/envfile"
	"github.com/renan-alm/gh-vars-migrator/internal/logger"
	"github.com/renan-alm/gh-vars-migrator/internal/migrator"
	"github.com/renan-alm/gh-vars-migrator/internal/types"
	"github.com/spf13/cobra"
)

var (
	// Version is set at build time
	Version = "dev"

	// Source flags
	sourceOrg      string
	sourceRepo     string
	sourcePAT      string
	sourceHostname string

	// Target flags
	targetOrg      string
	targetRepo     string
	targetPAT      string
	targetHostname string

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
  • Organization to organization variable migration (with automatic visibility preservation)
  • Repository to repository variable migration (with auto-discovery of environments)
  • Dry-run mode to preview changes before applying
  • Force mode to overwrite existing variables
  • Data residency compliance via custom GitHub hostnames

Mode Detection:
  - If --org-to-org flag is set → Organization migration mode
  - Otherwise → Repository-to-Repository migration mode (includes all environments)

Organization Variable Visibility:
  - Source variable visibility is automatically preserved during migration
  - Variables with 'selected' visibility have their repository selections matched
    by name in the target organisation
  - If no matching repositories are found, the variable is created with zero
    selected repositories

Authentication:
  - Primary: GITHUB_TOKEN environment variable (used for both source and target)
  - Override: --source-pat / --target-pat flags take precedence over GITHUB_TOKEN
  - Override: SOURCE_PAT / TARGET_PAT env vars (when flags are not provided)
  - Fallback: GitHub CLI authentication (gh auth login) when no tokens are set

Data Residency:
  - Use --source-hostname and --target-hostname to target specific GitHub Enterprise
    Server instances or data-residency-compliant GitHub Enterprise Cloud endpoints.
  - Variable values travel only between the specified source and target API endpoints,
    keeping data within your approved infrastructure.`,
	Example: `  # Organization to Organization migration (preserves source visibility)
  gh vars-migrator --source-org myorg --target-org targetorg --org-to-org

  # Repository to Repository migration (auto-discovers and migrates all environments)
  gh vars-migrator --source-org myorg --source-repo myrepo --target-org targetorg --target-repo targetrepo

  # Repository migration without environments
  gh vars-migrator --source-org myorg --source-repo myrepo --target-org targetorg --target-repo targetrepo --skip-envs

  # Dry-run mode (preview changes)
  gh vars-migrator --source-org myorg --target-org targetorg --org-to-org --dry-run

  # Force overwrite existing variables
  gh vars-migrator --source-org myorg --target-org targetorg --org-to-org --force

  # Using explicit PATs for different accounts
  gh vars-migrator --source-org myorg --target-org targetorg --org-to-org \
    --source-pat ghp_sourcetoken --target-pat ghp_targettoken

  # Using environment variables for tokens
  export SOURCE_PAT=ghp_sourcetoken
  export TARGET_PAT=ghp_targettoken
  gh vars-migrator --source-org myorg --target-org targetorg --org-to-org

  # Data residency: migrate between GitHub Enterprise Server instances
  gh vars-migrator --source-org myorg --target-org targetorg --org-to-org \
    --source-hostname github.source-company.com --target-hostname github.target-company.com \
    --source-pat ghp_sourcetoken --target-pat ghp_targettoken

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
	// Load .env file before registering flags so that os.Getenv picks up
	// file-defined values. Variables already set in the real environment
	// are never overwritten, and CLI flags always override env vars.
	if err := envfile.Load(".env"); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to load .env file: %v\n", err)
	}

	// Source flags
	rootCmd.Flags().StringVar(&sourceOrg, "source-org", os.Getenv("SOURCE_ORG"), "Source organization name (required) (env: SOURCE_ORG)")
	rootCmd.Flags().StringVar(&sourceRepo, "source-repo", os.Getenv("SOURCE_REPO"), "Source repository name (required for repo-to-repo) (env: SOURCE_REPO)")
	rootCmd.Flags().StringVar(&sourcePAT, "source-pat", os.Getenv("SOURCE_PAT"), "Source personal access token; overrides GITHUB_TOKEN (env: SOURCE_PAT)")
	rootCmd.Flags().StringVar(&sourceHostname, "source-hostname", os.Getenv("SOURCE_HOSTNAME"), "Source GitHub hostname for data residency (env: SOURCE_HOSTNAME)")

	// Target flags
	rootCmd.Flags().StringVar(&targetOrg, "target-org", os.Getenv("TARGET_ORG"), "Target organization name (required) (env: TARGET_ORG)")
	rootCmd.Flags().StringVar(&targetRepo, "target-repo", os.Getenv("TARGET_REPO"), "Target repository name (required for repo-to-repo) (env: TARGET_REPO)")
	rootCmd.Flags().StringVar(&targetPAT, "target-pat", os.Getenv("TARGET_PAT"), "Target personal access token; overrides GITHUB_TOKEN (env: TARGET_PAT)")
	rootCmd.Flags().StringVar(&targetHostname, "target-hostname", os.Getenv("TARGET_HOSTNAME"), "Target GitHub hostname for data residency (env: TARGET_HOSTNAME)")

	// Mode flags
	rootCmd.Flags().BoolVar(&orgToOrg, "org-to-org", envBool("ORG_TO_ORG"), "Migrate organization variables only (env: ORG_TO_ORG)")
	rootCmd.Flags().BoolVar(&skipEnvs, "skip-envs", envBool("SKIP_ENVS"), "Skip environment variable migration during repo-to-repo (env: SKIP_ENVS)")

	// Option flags
	rootCmd.Flags().BoolVar(&dryRun, "dry-run", envBool("DRY_RUN"), "Preview changes without applying them (env: DRY_RUN)")
	rootCmd.Flags().BoolVar(&force, "force", envBool("FORCE"), "Overwrite existing variables in target (env: FORCE)")

	// Global flags
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output")
}

// envBool returns true when the environment variable identified by key
// is set to a truthy value ("1", "true", "yes"). Any other value or an
// unset variable returns false.
func envBool(key string) bool {
	v := strings.ToLower(os.Getenv(key))
	return v == "1" || v == "true" || v == "yes"
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
	// Resolve tokens for source and target
	sourceToken, targetToken, err := resolveTokens()
	if err != nil {
		return err
	}

	// Create source and target clients
	sourceClient, targetClient, err := createClients(sourceToken, targetToken)
	if err != nil {
		return err
	}

	// Validate authentication
	if err := validateAuth(sourceClient, targetClient); err != nil {
		return err
	}

	// Detect migration mode
	mode := detectMigrationMode()

	// Validate PAT permissions before starting migration
	if err := validatePermissions(sourceClient, targetClient, mode); err != nil {
		return err
	}

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
		if sourceHostname != "" {
			logger.Info("Source Host: %s", sourceHostname)
		}
		logger.Info("Target: %s", targetOrg)
		if targetHostname != "" {
			logger.Info("Target Host: %s", targetHostname)
		}
		logger.Info("Org Visibility: preserve source")

	case types.ModeRepoToRepo:
		cfg.SourceOwner = sourceOrg
		cfg.SourceRepo = sourceRepo
		cfg.TargetOwner = targetOrg
		cfg.TargetRepo = targetRepo
		cfg.SkipEnvs = skipEnvs

		logger.Info("gh-vars-migrator - Repository Variable Migration")
		logger.Info("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		logger.Info("Source: %s/%s", cfg.SourceOwner, cfg.SourceRepo)
		if sourceHostname != "" {
			logger.Info("Source Host: %s", sourceHostname)
		}
		logger.Info("Target: %s/%s", cfg.TargetOwner, cfg.TargetRepo)
		if targetHostname != "" {
			logger.Info("Target Host: %s", targetHostname)
		}
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

	// Create and run migrator with both clients
	m, err := migrator.New(cfg, sourceClient, targetClient)
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

// resolveTokens determines which tokens to use for source and target.
//
// Priority per side (source / target):
//  1. --source-pat / --target-pat flag  (highest)
//  2. SOURCE_PAT / TARGET_PAT env var   (loaded as flag default)
//  3. GITHUB_TOKEN env var              (primary shared token)
//  4. GitHub CLI authentication         (lowest – empty string returned)
func resolveTokens() (sourceToken, targetToken string, err error) {
	githubToken := os.Getenv("GITHUB_TOKEN")

	// Start with GITHUB_TOKEN as the primary default for both sides.
	sourceToken = githubToken
	targetToken = githubToken

	// Override with explicit PATs when provided.
	if sourcePAT != "" {
		sourceToken = sourcePAT
	}
	if targetPAT != "" {
		targetToken = targetPAT
	}

	// Determine the label for each side's credential.
	sourceLabel := credentialLabel(sourcePAT, githubToken, "SOURCE_PAT", "GITHUB_TOKEN", "GitHub CLI")
	targetLabel := credentialLabel(targetPAT, githubToken, "TARGET_PAT", "GITHUB_TOKEN", "GitHub CLI")

	// Log which credential is used for each side.
	logger.Info("%s used for Source Org %s", sourceLabel, sourceOrg)
	logger.Info("%s used for Target Org %s", targetLabel, targetOrg)

	// Both resolved → done.
	if sourceToken != "" && targetToken != "" {
		return sourceToken, targetToken, nil
	}

	// Neither resolved → fall back to GitHub CLI authentication.
	if sourceToken == "" && targetToken == "" {
		return "", "", nil
	}

	// One side resolved, the other did not → cannot proceed.
	return "", "", fmt.Errorf("authentication required: please provide --source-pat and --target-pat flags, or set GITHUB_TOKEN environment variable")
}

// credentialLabel returns a human-readable label describing which credential
// was selected for one side of the migration (e.g. "SOURCE_PAT", "GITHUB_TOKEN",
// or "GitHub CLI").
func credentialLabel(pat, githubToken, patName, ghTokenName, cliFallback string) string {
	if pat != "" {
		return patName
	}
	if githubToken != "" {
		return ghTokenName
	}
	return cliFallback
}

// createClients creates source and target API clients
func createClients(sourceToken, targetToken string) (*client.Client, *client.Client, error) {
	var sourceClient, targetClient *client.Client
	var err error

	// Create source client
	sourceClient, err = createClientWithToken(sourceToken, sourceHostname, "source")
	if err != nil {
		return nil, nil, err
	}

	// Create target client
	targetClient, err = createClientWithToken(targetToken, targetHostname, "target")
	if err != nil {
		return nil, nil, err
	}

	return sourceClient, targetClient, nil
}

// createClientWithToken creates a client with an explicit token or default auth,
// optionally scoped to a custom GitHub hostname for data residency compliance.
func createClientWithToken(token string, hostname string, clientType string) (*client.Client, error) {
	if token != "" {
		if hostname != "" {
			c, err := client.NewWithTokenAndHost(token, hostname)
			if err != nil {
				return nil, fmt.Errorf("failed to create %s client with token and host: %w", clientType, err)
			}
			return c, nil
		}
		c, err := client.NewWithToken(token)
		if err != nil {
			return nil, fmt.Errorf("failed to create %s client with token: %w", clientType, err)
		}
		return c, nil
	}

	// Fallback to GitHub CLI authentication, with optional custom hostname
	if hostname != "" {
		c, err := client.NewWithHost(hostname)
		if err != nil {
			return nil, fmt.Errorf("failed to create %s client for host %s: %w", clientType, hostname, err)
		}
		return c, nil
	}

	c, err := client.New()
	if err != nil {
		return nil, fmt.Errorf("failed to create %s client: %w", clientType, err)
	}
	return c, nil
}

// validatePermissions validates that source and target tokens have the required
// OAuth scopes for the given migration mode. Validation is skipped when tokens
// do not expose scopes (e.g. fine-grained PATs or GITHUB_TOKEN).
func validatePermissions(sourceClient, targetClient *client.Client, mode types.MigrationMode) error {
	logger.Info("Validating token permissions...")

	switch mode {
	case types.ModeOrgToOrg:
		if err := client.ValidateOrgScopes(sourceClient, "source"); err != nil {
			return err
		}
		if err := client.ValidateOrgScopes(targetClient, "target"); err != nil {
			return err
		}
	case types.ModeRepoToRepo:
		if err := client.ValidateRepoScopes(sourceClient, "source"); err != nil {
			return err
		}
		if err := client.ValidateRepoScopes(targetClient, "target"); err != nil {
			return err
		}
	}

	logger.Success("Token permissions validated")
	return nil
}

// validateAuth validates that both source and target clients are authenticated
func validateAuth(sourceClient, targetClient *client.Client) error {
	// Validate source authentication
	sourceUser, err := sourceClient.GetUser()
	if err != nil {
		return fmt.Errorf("source authentication failed: %w\n\nPlease check your source credentials", err)
	}

	// Validate target authentication
	targetUser, err := targetClient.GetUser()
	if err != nil {
		return fmt.Errorf("target authentication failed: %w\n\nPlease check your target credentials", err)
	}

	logger.Success("Source authenticated as: %s", sourceUser)
	logger.Success("Target authenticated as: %s", targetUser)
	return nil
}

// checkAuth verifies that the user is authenticated with GitHub CLI (used by subcommands)
func checkAuth() error {
	restClient, err := api.DefaultRESTClient()
	if err != nil {
		return fmt.Errorf("failed to create GitHub API client: %w\n\nPlease authenticate using: gh auth login", err)
	}

	var user struct {
		Login string `json:"login"`
	}

	if err := restClient.Get("user", &user); err != nil {
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
