package client

import (
	"fmt"
	"strings"
)

// requiredOrgScopes lists the OAuth scopes needed for organization variable migration.
var requiredOrgScopes = []string{"admin:org"}

// requiredRepoScopes lists the OAuth scopes needed for repository and environment variable migration.
var requiredRepoScopes = []string{"repo"}

// hasScope reports whether a required scope is satisfied by any scope in the provided list.
// It handles parent–child relationships where a broader scope (e.g. "repo") implies
// narrower ones (e.g. "public_repo").
func hasScope(scopes []string, required string) bool {
	for _, s := range scopes {
		if s == required {
			return true
		}
		if isParentScope(s, required) {
			return true
		}
	}
	return false
}

// isParentScope reports whether parent is a superset that satisfies required.
func isParentScope(parent, required string) bool {
	hierarchy := map[string][]string{
		"admin:org": {"write:org", "read:org"},
		"repo":      {"public_repo", "repo:status", "repo:deployment", "repo:invite"},
	}
	for _, child := range hierarchy[parent] {
		if child == required {
			return true
		}
	}
	return false
}

// ValidateOrgScopes checks that the client token has the required scopes for
// organization variable migration. If the X-OAuth-Scopes header is absent
// (fine-grained PAT or GITHUB_TOKEN), validation is skipped.
func ValidateOrgScopes(c *Client, role string) error {
	scopes, err := c.GetTokenScopes()
	if err != nil {
		return fmt.Errorf("failed to retrieve %s token scopes: %w", role, err)
	}
	if scopes == nil {
		return nil
	}
	for _, required := range requiredOrgScopes {
		if !hasScope(scopes, required) {
			return fmt.Errorf(
				"%s token is missing required scope %q for organization variable migration\n"+
					"  Current scopes: %s\n"+
					"  Please create a PAT with the 'admin:org' scope at https://github.com/settings/tokens",
				role, required, strings.Join(scopes, ", "),
			)
		}
	}
	return nil
}

// ValidateRepoScopes checks that the client token has the required scopes for
// repository and environment variable migration. If the X-OAuth-Scopes header
// is absent (fine-grained PAT or GITHUB_TOKEN), validation is skipped.
func ValidateRepoScopes(c *Client, role string) error {
	scopes, err := c.GetTokenScopes()
	if err != nil {
		return fmt.Errorf("failed to retrieve %s token scopes: %w", role, err)
	}
	if scopes == nil {
		return nil
	}
	for _, required := range requiredRepoScopes {
		if !hasScope(scopes, required) {
			return fmt.Errorf(
				"%s token is missing required scope %q for repository variable migration\n"+
					"  Current scopes: %s\n"+
					"  Please create a PAT with the 'repo' scope at https://github.com/settings/tokens",
				role, required, strings.Join(scopes, ", "),
			)
		}
	}
	return nil
}
