package config

import (
	"testing"

	"github.com/renan-alm/gh-vars-migrator/internal/types"
)

func TestValidate_NilConfig(t *testing.T) {
	err := Validate(nil)
	if err == nil {
		t.Error("Expected error for nil config")
	}
}

func TestValidate_RepoToRepo(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *types.MigrationConfig
		wantErr bool
	}{
		{
			name: "valid config",
			cfg: &types.MigrationConfig{
				Mode:        types.ModeRepoToRepo,
				SourceOwner: "source-owner",
				SourceRepo:  "source-repo",
				TargetOwner: "target-owner",
				TargetRepo:  "target-repo",
			},
			wantErr: false,
		},
		{
			name: "missing source owner",
			cfg: &types.MigrationConfig{
				Mode:        types.ModeRepoToRepo,
				SourceRepo:  "source-repo",
				TargetOwner: "target-owner",
				TargetRepo:  "target-repo",
			},
			wantErr: true,
		},
		{
			name: "missing source repo",
			cfg: &types.MigrationConfig{
				Mode:        types.ModeRepoToRepo,
				SourceOwner: "source-owner",
				TargetOwner: "target-owner",
				TargetRepo:  "target-repo",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidate_OrgToOrg(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *types.MigrationConfig
		wantErr bool
	}{
		{
			name: "valid config",
			cfg: &types.MigrationConfig{
				Mode:      types.ModeOrgToOrg,
				SourceOrg: "source-org",
				TargetOrg: "target-org",
			},
			wantErr: false,
		},
		{
			name: "missing source org",
			cfg: &types.MigrationConfig{
				Mode:      types.ModeOrgToOrg,
				TargetOrg: "target-org",
			},
			wantErr: true,
		},
		{
			name: "missing target org",
			cfg: &types.MigrationConfig{
				Mode:      types.ModeOrgToOrg,
				SourceOrg: "source-org",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetDescription(t *testing.T) {
	tests := []struct {
		name string
		cfg  *types.MigrationConfig
		want string
	}{
		{
			name: "repo to repo with envs",
			cfg: &types.MigrationConfig{
				Mode:        types.ModeRepoToRepo,
				SourceOwner: "org1",
				SourceRepo:  "repo1",
				TargetOwner: "org2",
				TargetRepo:  "repo2",
				SkipEnvs:    false,
			},
			want: "Repository org1/repo1 → org2/repo2 (with environments)",
		},
		{
			name: "repo to repo skip envs",
			cfg: &types.MigrationConfig{
				Mode:        types.ModeRepoToRepo,
				SourceOwner: "org1",
				SourceRepo:  "repo1",
				TargetOwner: "org2",
				TargetRepo:  "repo2",
				SkipEnvs:    true,
			},
			want: "Repository org1/repo1 → org2/repo2",
		},
		{
			name: "org to org",
			cfg: &types.MigrationConfig{
				Mode:      types.ModeOrgToOrg,
				SourceOrg: "org1",
				TargetOrg: "org2",
			},
			want: "Organization org1 → org2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetDescription(tt.cfg)
			if got != tt.want {
				t.Errorf("GetDescription() = %v, want %v", got, tt.want)
			}
		})
	}
}
