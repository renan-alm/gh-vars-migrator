package types

import "testing"

func TestMigrationResult_AddError(t *testing.T) {
	result := &MigrationResult{}

	if result.HasErrors() {
		t.Error("Expected no errors initially")
	}

	result.AddError(nil)
	if !result.HasErrors() {
		t.Error("Expected to have errors after adding one")
	}
}

func TestMigrationResult_Total(t *testing.T) {
	result := &MigrationResult{
		Created: 5,
		Updated: 3,
		Skipped: 2,
	}

	expected := 10
	if result.Total() != expected {
		t.Errorf("Expected total %d, got %d", expected, result.Total())
	}
}

func TestMigrationMode_Constants(t *testing.T) {
	modes := []MigrationMode{
		ModeRepoToRepo,
		ModeOrgToOrg,
		ModeEnvOnly,
	}

	for _, mode := range modes {
		if mode == "" {
			t.Errorf("Migration mode should not be empty")
		}
	}
}
