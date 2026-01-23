package main

import (
	"fmt"

	"github.com/renan-alm/gh-vars-migrator/internal/logger"
	"github.com/renan-alm/gh-vars-migrator/internal/types"
)

func main() {
	logger.Info("gh-vars-migrator - GitHub CLI extension for variables migration")

	// Example usage - in a real CLI, these would come from flags
	// For now, we'll just demonstrate the structure is wired up correctly

	// Example configuration for repo-to-repo migration
	cfg := &types.MigrationConfig{
		Mode:        types.ModeRepoToRepo,
		SourceOwner: "example-org",
		SourceRepo:  "source-repo",
		TargetOwner: "example-org",
		TargetRepo:  "target-repo",
		DryRun:      true,
		Force:       false,
		SkipEnvs:    true,
	}

	// Show usage help for now
	printUsageHelp()

	// Example of how to use the migrator (commented out to avoid API calls)
	_ = cfg
	/*
		m, err := migrator.New(cfg)
		if err != nil {
			logger.Error("Failed to initialize migrator: %v", err)
			os.Exit(1)
		}

		result, err := m.Run()
		if err != nil {
			logger.Error("Migration failed: %v", err)
			os.Exit(1)
		}

		if result.HasErrors() {
			os.Exit(1)
		}
	*/
}

func printUsageHelp() {
	logger.Plain("\nCore migrator components are now available!")
	logger.Plain("\nAvailable migration modes:")
	logger.Plain("  • repo-to-repo  - Migrate variables between repositories")
	logger.Plain("  • org-to-org    - Migrate variables between organizations")
	logger.Plain("  • env-only      - Migrate variables between environments")

	logger.Plain("\nNext steps:")
	logger.Plain("  1. Implement CLI flag parsing")
	logger.Plain("  2. Wire up CLI commands to migrator")
	logger.Plain("  3. Add proper authentication check")

	logger.Plain("\nExample programmatic usage:")
	fmt.Println(`
  cfg := &types.MigrationConfig{
    Mode:        types.ModeRepoToRepo,
    SourceOwner: "source-org",
    SourceRepo:  "source-repo",
    TargetOwner: "target-org",
    TargetRepo:  "target-repo",
    DryRun:      true,
  }
  
  m, err := migrator.New(cfg)
  if err != nil {
    log.Fatal(err)
  }
  
  result, err := m.Run()
  // Handle result...
	`)
}
