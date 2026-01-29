package migrator

import (
	"fmt"

	"github.com/renan-alm/gh-vars-migrator/internal/client"
	"github.com/renan-alm/gh-vars-migrator/internal/config"
	"github.com/renan-alm/gh-vars-migrator/internal/logger"
	"github.com/renan-alm/gh-vars-migrator/internal/types"
)

// Migrator orchestrates the migration of GitHub Actions variables
type Migrator struct {
	sourceClient *client.Client
	targetClient *client.Client
	config       *types.MigrationConfig
}

// New creates a new Migrator instance with separate source and target clients
func New(cfg *types.MigrationConfig, sourceClient, targetClient *client.Client) (*Migrator, error) {
	// Validate configuration
	if err := config.Validate(cfg); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	if sourceClient == nil {
		return nil, fmt.Errorf("source client cannot be nil")
	}

	if targetClient == nil {
		return nil, fmt.Errorf("target client cannot be nil")
	}

	return &Migrator{
		sourceClient: sourceClient,
		targetClient: targetClient,
		config:       cfg,
	}, nil
}

// Run executes the migration based on the configuration
func (m *Migrator) Run() (*types.MigrationResult, error) {
	logger.Info("Starting migration: %s", config.GetDescription(m.config))

	if m.config.DryRun {
		logger.Warning("Running in DRY-RUN mode - no changes will be made")
	}

	var result *types.MigrationResult
	var err error

	switch m.config.Mode {
	case types.ModeRepoToRepo:
		result, err = m.migrateRepoToRepo()
	case types.ModeOrgToOrg:
		result, err = m.migrateOrgToOrg()
	default:
		return nil, fmt.Errorf("unsupported migration mode: %s", m.config.Mode)
	}

	if err != nil {
		return result, err
	}

	// Print summary
	logger.PrintSummary(result.Created, result.Updated, result.Skipped, len(result.Errors))

	// Print errors if any
	if result.HasErrors() {
		logger.Error("\nEncountered %d error(s) during migration:", len(result.Errors))
		for i, err := range result.Errors {
			logger.Error("  %d. %v", i+1, err)
		}
	}

	return result, nil
}
