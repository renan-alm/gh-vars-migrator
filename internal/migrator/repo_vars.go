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
	if !m.config.SkipEnvs {
		if err := m.migrateAllEnvironments(result); err != nil {
			logger.Warning("Failed to migrate environments: %v", err)
			result.AddError(fmt.Errorf("environment migration failed: %w", err))
		}
	} else {
		logger.Info("Skipping environment variable migration (--skip-envs)")
	}

	return result, nil
}

// migrateAllEnvironments discovers all environments from source repo and migrates them
func (m *Migrator) migrateAllEnvironments(result *types.MigrationResult) error {
	logger.Info("Discovering environments from source repository: %s/%s", m.config.SourceOwner, m.config.SourceRepo)

	// List all environments from source repository
	environments, err := m.client.ListEnvironments(m.config.SourceOwner, m.config.SourceRepo)
	if err != nil {
		return fmt.Errorf("failed to list environments: %w", err)
	}

	if len(environments) == 0 {
		logger.Info("No environments found in source repository")
		return nil
	}

	logger.Info("Found %d environment(s): %v", len(environments), getEnvNames(environments))

	// Migrate each environment
	for _, env := range environments {
		if err := m.migrateEnvironment(env.Name, result); err != nil {
			logger.Error("Failed to migrate environment '%s': %v", env.Name, err)
			result.AddError(fmt.Errorf("environment '%s': %w", env.Name, err))
		}
	}

	return nil
}

// getEnvNames extracts environment names for logging
func getEnvNames(envs []types.Environment) []string {
	names := make([]string, len(envs))
	for i, env := range envs {
		names[i] = env.Name
	}
	return names
}

// migrateEnvironment migrates a single environment and its variables
func (m *Migrator) migrateEnvironment(envName string, result *types.MigrationResult) error {
	logger.Info("Migrating environment: %s", envName)

	// Check if environment exists in target, create if not
	if err := m.ensureEnvironmentExists(envName); err != nil {
		return fmt.Errorf("failed to ensure environment exists: %w", err)
	}

	// Get variables from source environment
	sourceEnvVars, err := m.client.ListEnvVariables(m.config.SourceOwner, m.config.SourceRepo, envName)
	if err != nil {
		return fmt.Errorf("failed to list environment variables: %w", err)
	}

	logger.Info("Found %d variable(s) in environment '%s'", len(sourceEnvVars), envName)

	// Migrate each variable in this environment
	for _, variable := range sourceEnvVars {
		if err := m.migrateEnvVariable(envName, variable, result); err != nil {
			logger.Error("Failed to migrate environment variable '%s': %v", variable.Name, err)
			result.AddError(fmt.Errorf("env '%s' variable '%s': %w", envName, variable.Name, err))
		}
	}

	return nil
}

// ensureEnvironmentExists creates the environment in the target repo if it doesn't exist
func (m *Migrator) ensureEnvironmentExists(envName string) error {
	// Check if environment already exists in target
	_, err := m.client.GetEnvironment(m.config.TargetOwner, m.config.TargetRepo, envName)
	if err == nil {
		logger.Debug("Environment '%s' already exists in target repository", envName)
		return nil
	}

	// Environment doesn't exist, create it
	if m.config.DryRun {
		logger.Info("[DRY-RUN] Would create environment: %s", envName)
		return nil
	}

	logger.Info("Creating environment '%s' in target repository", envName)
	if err := m.client.CreateEnvironment(m.config.TargetOwner, m.config.TargetRepo, envName); err != nil {
		return fmt.Errorf("failed to create environment: %w", err)
	}

	logger.Success("Created environment: %s", envName)
	return nil
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

// migrateEnvVariable migrates a single environment variable
func (m *Migrator) migrateEnvVariable(envName string, variable types.Variable, result *types.MigrationResult) error {
	// Check if variable exists in target environment
	existingVar, err := m.client.GetEnvVariable(m.config.TargetOwner, m.config.TargetRepo, envName, variable.Name)

	if err == nil && existingVar != nil {
		// Variable exists in target environment
		if !m.config.Force {
			logger.Warning("Environment variable '%s' already exists in target (use --force to overwrite)", variable.Name)
			result.Skipped++
			return nil
		}

		// Update existing variable
		if m.config.DryRun {
			logger.Info("[DRY-RUN] Would update environment variable: %s (env: %s)", variable.Name, envName)
			result.Updated++
			return nil
		}

		if err := m.client.UpdateEnvVariable(m.config.TargetOwner, m.config.TargetRepo, envName, variable); err != nil {
			return fmt.Errorf("failed to update: %w", err)
		}

		logger.Success("Updated environment variable: %s (env: %s)", variable.Name, envName)
		result.Updated++
		return nil
	}

	// Create new environment variable
	if m.config.DryRun {
		logger.Info("[DRY-RUN] Would create environment variable: %s (env: %s)", variable.Name, envName)
		result.Created++
		return nil
	}

	if err := m.client.CreateEnvVariable(m.config.TargetOwner, m.config.TargetRepo, envName, variable); err != nil {
		return fmt.Errorf("failed to create: %w", err)
	}

	logger.Success("Created environment variable: %s (env: %s)", variable.Name, envName)
	result.Created++
	return nil
}
