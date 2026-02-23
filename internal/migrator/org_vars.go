package migrator

import (
	"fmt"

	"github.com/renan-alm/gh-vars-migrator/internal/logger"
	"github.com/renan-alm/gh-vars-migrator/internal/types"
)

// migrateOrgToOrg handles organization-to-organization variable migration
func (m *Migrator) migrateOrgToOrg() (*types.MigrationResult, error) {
	result := &types.MigrationResult{}

	// Check rate limit before starting the API-intensive migration
	m.sourceClient.WaitForRateLimit()

	logger.Info("Fetching variables from source organization: %s", m.config.SourceOrg)

	// Get source organization variables using source client
	sourceVars, err := m.sourceClient.ListOrgVariables(m.config.SourceOrg)
	if err != nil {
		return result, fmt.Errorf("failed to list source organization variables: %w", err)
	}

	logger.Info("Found %d variable(s) in source organization", len(sourceVars))

	// Migrate each variable, preserving source visibility
	for _, variable := range sourceVars {
		if variable.Visibility == "" {
			variable.Visibility = "all"
		}

		// For "selected" visibility, resolve the repository selection from source
		// and match by name in the target organisation.
		if variable.Visibility == "selected" {
			selectedIDs, err := m.resolveSelectedRepos(variable.Name)
			if err != nil {
				logger.Warning("Failed to resolve selected repositories for variable '%s': %v; migrating with empty repository list", variable.Name, err)
			}
			variable.SelectedRepositoryIDs = selectedIDs

			if len(selectedIDs) == 0 {
				logger.Warning("Variable '%s' has 'selected' visibility but no matching repositories were found in target organization '%s'; it will be created with zero selected repositories", variable.Name, m.config.TargetOrg)
			} else {
				logger.Info("Variable '%s': matched %d repository(ies) by name in target organization", variable.Name, len(selectedIDs))
			}
		}

		if err := m.migrateOrgVariable(variable, result); err != nil {
			logger.Error("Failed to migrate variable '%s': %v", variable.Name, err)
			result.AddError(fmt.Errorf("variable '%s': %w", variable.Name, err))
		}
	}

	return result, nil
}

// resolveSelectedRepos fetches the selected repositories for a source variable
// and looks up repositories with matching names in the target organisation.
// Returns the target repository IDs for any names that match.
func (m *Migrator) resolveSelectedRepos(varName string) ([]int64, error) {
	sourceRepos, err := m.sourceClient.ListOrgVariableSelectedRepos(m.config.SourceOrg, varName)
	if err != nil {
		return nil, fmt.Errorf("failed to list selected repos from source: %w", err)
	}

	if len(sourceRepos) == 0 {
		return nil, nil
	}

	var targetIDs []int64
	for _, srcRepo := range sourceRepos {
		targetRepo, err := m.targetClient.GetRepo(m.config.TargetOrg, srcRepo.Name)
		if err != nil {
			logger.Debug("Repository '%s' not found in target organization '%s': %v", srcRepo.Name, m.config.TargetOrg, err)
			continue
		}
		logger.Debug("Matched repository '%s' (source ID %d -> target ID %d)", srcRepo.Name, srcRepo.ID, targetRepo.ID)
		targetIDs = append(targetIDs, targetRepo.ID)
	}

	return targetIDs, nil
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
