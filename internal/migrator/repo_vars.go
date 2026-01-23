package migrator

import (
	"fmt"

	"github.com/renan-alm/gh-vars-migrator/internal/logger"
	"github.com/renan-alm/gh-vars-migrator/internal/types"
)

// migrateRepoToRepo handles repository-to-repository variable migration
func (m *Migrator) migrateRepoToRepo() (*types.MigrationResult, error) {
	result := &types.MigrationResult{}

	logger.Info("Fetching variables from source repository: %s/%s", m.config.SourceOwner, m.config.SourceRepo)

	// Get source repository variables
	sourceVars, err := m.client.ListRepoVariables(m.config.SourceOwner, m.config.SourceRepo)
	if err != nil {
		return result, fmt.Errorf("failed to list source repository variables: %w", err)
	}

	logger.Info("Found %d variable(s) in source repository", len(sourceVars))

	// Migrate repository-level variables
	if err := m.migrateRepoVariables(sourceVars, result); err != nil {
		return result, err
	}

	// Migrate environment variables if not skipped
	if !m.config.SkipEnvs && m.config.SourceEnv != "" && m.config.TargetEnv != "" {
		logger.Info("Migrating environment variables from '%s' to '%s'", m.config.SourceEnv, m.config.TargetEnv)
		if err := m.migrateEnvironmentVariables(result); err != nil {
			logger.Warning("Failed to migrate environment variables: %v", err)
			result.AddError(fmt.Errorf("environment migration failed: %w", err))
		}
	} else if !m.config.SkipEnvs {
		logger.Debug("No environment variables to migrate (source or target environment not specified)")
	}

	return result, nil
}

// migrateRepoVariables migrates repository-level variables
func (m *Migrator) migrateRepoVariables(sourceVars []types.Variable, result *types.MigrationResult) error {
	for _, variable := range sourceVars {
		if err := m.migrateRepoVariable(variable, result); err != nil {
			logger.Error("Failed to migrate variable '%s': %v", variable.Name, err)
			result.AddError(fmt.Errorf("variable '%s': %w", variable.Name, err))
		}
	}
	return nil
}

// migrateRepoVariable migrates a single repository variable
func (m *Migrator) migrateRepoVariable(variable types.Variable, result *types.MigrationResult) error {
	// Check if variable exists in target
	existingVar, err := m.client.GetRepoVariable(m.config.TargetOwner, m.config.TargetRepo, variable.Name)

	if err == nil && existingVar != nil {
		// Variable exists in target
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

		if err := m.client.UpdateRepoVariable(m.config.TargetOwner, m.config.TargetRepo, variable); err != nil {
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

	if err := m.client.CreateRepoVariable(m.config.TargetOwner, m.config.TargetRepo, variable); err != nil {
		return fmt.Errorf("failed to create: %w", err)
	}

	logger.Success("Created variable: %s", variable.Name)
	result.Created++
	return nil
}

// migrateEnvironmentVariables migrates environment-level variables
func (m *Migrator) migrateEnvironmentVariables(result *types.MigrationResult) error {
	// Get source environment variables
	sourceEnvVars, err := m.client.ListEnvVariables(m.config.SourceOwner, m.config.SourceRepo, m.config.SourceEnv)
	if err != nil {
		return fmt.Errorf("failed to list source environment variables: %w", err)
	}

	logger.Info("Found %d environment variable(s) in '%s'", len(sourceEnvVars), m.config.SourceEnv)

	for _, variable := range sourceEnvVars {
		if err := m.migrateEnvVariable(variable, result); err != nil {
			logger.Error("Failed to migrate environment variable '%s': %v", variable.Name, err)
			result.AddError(fmt.Errorf("env variable '%s': %w", variable.Name, err))
		}
	}

	return nil
}

// migrateEnvVariable migrates a single environment variable
func (m *Migrator) migrateEnvVariable(variable types.Variable, result *types.MigrationResult) error {
	// Check if variable exists in target environment
	existingVar, err := m.client.GetEnvVariable(m.config.TargetOwner, m.config.TargetRepo, m.config.TargetEnv, variable.Name)

	if err == nil && existingVar != nil {
		// Variable exists in target environment
		if !m.config.Force {
			logger.Warning("Environment variable '%s' already exists in target (use --force to overwrite)", variable.Name)
			result.Skipped++
			return nil
		}

		// Update existing variable
		if m.config.DryRun {
			logger.Info("[DRY-RUN] Would update environment variable: %s", variable.Name)
			result.Updated++
			return nil
		}

		if err := m.client.UpdateEnvVariable(m.config.TargetOwner, m.config.TargetRepo, m.config.TargetEnv, variable); err != nil {
			return fmt.Errorf("failed to update: %w", err)
		}

		logger.Success("Updated environment variable: %s", variable.Name)
		result.Updated++
		return nil
	}

	// Create new environment variable
	if m.config.DryRun {
		logger.Info("[DRY-RUN] Would create environment variable: %s", variable.Name)
		result.Created++
		return nil
	}

	if err := m.client.CreateEnvVariable(m.config.TargetOwner, m.config.TargetRepo, m.config.TargetEnv, variable); err != nil {
		return fmt.Errorf("failed to create: %w", err)
	}

	logger.Success("Created environment variable: %s", variable.Name)
	result.Created++
	return nil
}
