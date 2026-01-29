package migrator

import (
	"fmt"

	"github.com/renan-alm/gh-vars-migrator/internal/logger"
	"github.com/renan-alm/gh-vars-migrator/internal/types"
)

// migrateOrgToOrg handles organization-to-organization variable migration
func (m *Migrator) migrateOrgToOrg() (*types.MigrationResult, error) {
	result := &types.MigrationResult{}

	logger.Info("Fetching variables from source organization: %s", m.config.SourceOrg)

	// Get source organization variables using source client
	sourceVars, err := m.sourceClient.ListOrgVariables(m.config.SourceOrg)
	if err != nil {
		return result, fmt.Errorf("failed to list source organization variables: %w", err)
	}

	logger.Info("Found %d variable(s) in source organization", len(sourceVars))

	// Migrate each variable
	for _, variable := range sourceVars {
		if err := m.migrateOrgVariable(variable, result); err != nil {
			logger.Error("Failed to migrate variable '%s': %v", variable.Name, err)
			result.AddError(fmt.Errorf("variable '%s': %w", variable.Name, err))
		}
	}

	return result, nil
}

// migrateOrgVariable migrates a single organization variable
func (m *Migrator) migrateOrgVariable(variable types.Variable, result *types.MigrationResult) error {
	// Check if variable exists in target using target client
	existingVar, err := m.targetClient.GetOrgVariable(m.config.TargetOrg, variable.Name)

	if err == nil && existingVar != nil {
		// Variable exists in target
		if !m.config.Force {
			logger.Warning("Variable '%s' already exists in target (use --force to overwrite)", variable.Name)
			result.Skipped++
			return nil
		}

		// Update existing variable using target client
		if m.config.DryRun {
			logger.Info("[DRY-RUN] Would update variable: %s", variable.Name)
			result.Updated++
			return nil
		}

		if err := m.targetClient.UpdateOrgVariable(m.config.TargetOrg, variable); err != nil {
			return fmt.Errorf("failed to update: %w", err)
		}

		logger.Success("Updated variable: %s", variable.Name)
		result.Updated++
		return nil
	}

	// Create new variable using target client
	if m.config.DryRun {
		logger.Info("[DRY-RUN] Would create variable: %s", variable.Name)
		result.Created++
		return nil
	}

	if err := m.targetClient.CreateOrgVariable(m.config.TargetOrg, variable); err != nil {
		return fmt.Errorf("failed to create: %w", err)
	}

	logger.Success("Created variable: %s", variable.Name)
	result.Created++
	return nil
}
