package client

import (
	"testing"
)

// TestHasScope verifies that hasScope correctly identifies satisfied scopes.
func TestHasScope(t *testing.T) {
	tests := []struct {
		name     string
		scopes   []string
		required string
		want     bool
	}{
		{
			name:     "exact match",
			scopes:   []string{"repo", "admin:org"},
			required: "repo",
			want:     true,
		},
		{
			name:     "missing scope",
			scopes:   []string{"read:org"},
			required: "admin:org",
			want:     false,
		},
		{
			name:     "parent scope satisfies child",
			scopes:   []string{"admin:org"},
			required: "read:org",
			want:     true,
		},
		{
			name:     "repo parent satisfies public_repo",
			scopes:   []string{"repo"},
			required: "public_repo",
			want:     true,
		},
		{
			name:     "empty scopes",
			scopes:   []string{},
			required: "repo",
			want:     false,
		},
		{
			name:     "unrelated scopes",
			scopes:   []string{"read:user", "gist"},
			required: "admin:org",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := hasScope(tt.scopes, tt.required)
			if got != tt.want {
				t.Errorf("hasScope(%v, %q) = %v, want %v", tt.scopes, tt.required, got, tt.want)
			}
		})
	}
}

// TestIsParentScope verifies the parent-scope hierarchy relationships.
func TestIsParentScope(t *testing.T) {
	tests := []struct {
		parent   string
		required string
		want     bool
	}{
		{"admin:org", "write:org", true},
		{"admin:org", "read:org", true},
		{"admin:org", "repo", false},
		{"repo", "public_repo", true},
		{"repo", "repo:status", true},
		{"repo", "admin:org", false},
		{"read:org", "admin:org", false},
	}

	for _, tt := range tests {
		t.Run(tt.parent+"->"+tt.required, func(t *testing.T) {
			got := isParentScope(tt.parent, tt.required)
			if got != tt.want {
				t.Errorf("isParentScope(%q, %q) = %v, want %v", tt.parent, tt.required, got, tt.want)
			}
		})
	}
}

// mockClientWithScopes creates a Client whose GetTokenScopes returns the given scopes.
// It wraps the logic under test (ValidateOrgScopes / ValidateRepoScopes) directly,
// since we cannot easily mock the HTTP layer without changing production code.
// Instead, we test the scope-checking helpers independently.

// TestValidateOrgScopes_WithSufficientScopes verifies no error when admin:org is present.
func TestValidateOrgScopes_WithSufficientScopes(t *testing.T) {
	scopes := []string{"admin:org", "repo"}
	for _, required := range requiredOrgScopes {
		if !hasScope(scopes, required) {
			t.Errorf("expected scopes %v to satisfy required org scope %q", scopes, required)
		}
	}
}

// TestValidateOrgScopes_WithMissingScopes verifies that missing admin:org is detected.
func TestValidateOrgScopes_WithMissingScopes(t *testing.T) {
	scopes := []string{"repo", "read:user"}
	for _, required := range requiredOrgScopes {
		if hasScope(scopes, required) {
			t.Errorf("expected scopes %v to NOT satisfy required org scope %q", scopes, required)
		}
	}
}

// TestValidateRepoScopes_WithSufficientScopes verifies no error when repo is present.
func TestValidateRepoScopes_WithSufficientScopes(t *testing.T) {
	scopes := []string{"repo", "workflow"}
	for _, required := range requiredRepoScopes {
		if !hasScope(scopes, required) {
			t.Errorf("expected scopes %v to satisfy required repo scope %q", scopes, required)
		}
	}
}

// TestValidateRepoScopes_WithMissingScopes verifies that missing repo is detected.
func TestValidateRepoScopes_WithMissingScopes(t *testing.T) {
	scopes := []string{"read:user", "gist"}
	for _, required := range requiredRepoScopes {
		if hasScope(scopes, required) {
			t.Errorf("expected scopes %v to NOT satisfy required repo scope %q", scopes, required)
		}
	}
}

// TestValidateRepoScopes_PublicRepoNotSufficient verifies that public_repo alone
// does not satisfy the full repo scope requirement.
func TestValidateRepoScopes_PublicRepoNotSufficient(t *testing.T) {
	scopes := []string{"public_repo"}
	if hasScope(scopes, "repo") {
		t.Error("expected public_repo alone to NOT satisfy full repo scope requirement")
	}
}
