package migrator

import (
	"fmt"

	"github.com/renan-alm/gh-vars-migrator/internal/logger"
	"github.com/renan-alm/gh-vars-migrator/internal/types"
)

// migrateEnvOnly handles environment-only variable migration
func (m *Migrator) migrateEnvOnly() (*types.MigrationResult, error) {
	result := &types.MigrationResult{}

	logger.Info("Fetching variables from source environment: %s (repository: %s/%s)",
		m.config.SourceEnv, m.config.SourceOwner, m.config.SourceRepo)

	// Get source environment variables
	sourceVars, err := m.client.ListEnvVariables(m.config.SourceOwner, m.config.SourceRepo, m.config.SourceEnv)
	if err != nil {
		return result, fmt.Errorf("failed to list source environment variables: %w", err)
	}

	logger.Info("Found %d variable(s) in source environment", len(sourceVars))

	// Determine target repository (same as source if not specified)
	targetOwner := m.config.TargetOwner
	targetRepo := m.config.TargetRepo
	if targetOwner == "" {
		targetOwner = m.config.SourceOwner
	}
	if targetRepo == "" {
		targetRepo = m.config.SourceRepo
	}

	logger.Info("Migrating to target environment: %s (repository: %s/%s)",
		m.config.TargetEnv, targetOwner, targetRepo)

	// Migrate each variable
	for _, variable := range sourceVars {
		if err := m.migrateEnvOnlyVariable(variable, targetOwner, targetRepo, result); err != nil {
			logger.Error("Failed to migrate variable '%s': %v", variable.Name, err)
			result.AddError(fmt.Errorf("variable '%s': %w", variable.Name, err))
		}
	}

	return result, nil
}

// migrateEnvOnlyVariable migrates a single environment variable
func (m *Migrator) migrateEnvOnlyVariable(variable types.Variable, targetOwner, targetRepo string, result *types.MigrationResult) error {
	// Check if variable exists in target environment
	existingVar, err := m.client.GetEnvVariable(targetOwner, targetRepo, m.config.TargetEnv, variable.Name)

	if err == nil && existingVar != nil {
		// Variable exists in target environment
		if !m.config.Force {
			logger.Warning("Variable '%s' already exists in target (use --force to overwrite)", variable.Name)
			result.Skipped++
			return nil
		}

		// Update existing variable
		if m.config.DryRun {
			logger.Info("[DRY-RUN] Would update variable: %s", variable.Name)
			result.Updated++
			return nil
		}

		if err := m.client.UpdateEnvVariable(targetOwner, targetRepo, m.config.TargetEnv, variable); err != nil {
			return fmt.Errorf("failed to update: %w", err)
		}

		logger.Success("Updated variable: %s", variable.Name)
		result.Updated++
		return nil
	}

	// Create new variable
	if m.config.DryRun {
		logger.Info("[DRY-RUN] Would create variable: %s", variable.Name)
		result.Created++
		return nil
	}

	if err := m.client.CreateEnvVariable(targetOwner, targetRepo, m.config.TargetEnv, variable); err != nil {
		return fmt.Errorf("failed to create: %w", err)
	}

	logger.Success("Created variable: %s", variable.Name)
	result.Created++
	return nil
}
