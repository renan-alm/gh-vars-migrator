package cmd

import (
	"os"
	"testing"
)

// TestResolveTokens_BothPATsProvided tests when both source and target PATs are explicitly provided
func TestResolveTokens_BothPATsProvided(t *testing.T) {
	// Save original values
	origSourcePAT := sourcePAT
	origTargetPAT := targetPAT
	origGitHubToken := os.Getenv("GITHUB_TOKEN")
	
	// Clean up after test
	defer func() {
		sourcePAT = origSourcePAT
		targetPAT = origTargetPAT
		if origGitHubToken != "" {
			os.Setenv("GITHUB_TOKEN", origGitHubToken)
		} else {
			os.Unsetenv("GITHUB_TOKEN")
		}
	}()

	// Set test values
	sourcePAT = "source_token_123"
	targetPAT = "target_token_456"
	os.Unsetenv("GITHUB_TOKEN")

	sourceToken, targetToken, err := resolveTokens()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if sourceToken != "source_token_123" {
		t.Errorf("Expected source token 'source_token_123', got '%s'", sourceToken)
	}

	if targetToken != "target_token_456" {
		t.Errorf("Expected target token 'target_token_456', got '%s'", targetToken)
	}
}

// TestResolveTokens_GitHubTokenFallback tests when GITHUB_TOKEN is used as fallback
func TestResolveTokens_GitHubTokenFallback(t *testing.T) {
	// Save original values
	origSourcePAT := sourcePAT
	origTargetPAT := targetPAT
	origGitHubToken := os.Getenv("GITHUB_TOKEN")
	
	// Clean up after test
	defer func() {
		sourcePAT = origSourcePAT
		targetPAT = origTargetPAT
		if origGitHubToken != "" {
			os.Setenv("GITHUB_TOKEN", origGitHubToken)
		} else {
			os.Unsetenv("GITHUB_TOKEN")
		}
	}()

	// Set test values
	sourcePAT = ""
	targetPAT = ""
	os.Setenv("GITHUB_TOKEN", "github_token_789")

	sourceToken, targetToken, err := resolveTokens()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if sourceToken != "github_token_789" {
		t.Errorf("Expected source token 'github_token_789', got '%s'", sourceToken)
	}

	if targetToken != "github_token_789" {
		t.Errorf("Expected target token 'github_token_789', got '%s'", targetToken)
	}
}

// TestResolveTokens_MixedMode tests when one PAT is provided and GITHUB_TOKEN fills the gap
func TestResolveTokens_MixedMode(t *testing.T) {
	// Save original values
	origSourcePAT := sourcePAT
	origTargetPAT := targetPAT
	origGitHubToken := os.Getenv("GITHUB_TOKEN")
	
	// Clean up after test
	defer func() {
		sourcePAT = origSourcePAT
		targetPAT = origTargetPAT
		if origGitHubToken != "" {
			os.Setenv("GITHUB_TOKEN", origGitHubToken)
		} else {
			os.Unsetenv("GITHUB_TOKEN")
		}
	}()

	// Test case 1: Only source PAT provided
	sourcePAT = "source_token_abc"
	targetPAT = ""
	os.Setenv("GITHUB_TOKEN", "github_token_xyz")

	sourceToken, targetToken, err := resolveTokens()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if sourceToken != "source_token_abc" {
		t.Errorf("Expected source token 'source_token_abc', got '%s'", sourceToken)
	}

	if targetToken != "github_token_xyz" {
		t.Errorf("Expected target token 'github_token_xyz', got '%s'", targetToken)
	}

	// Test case 2: Only target PAT provided
	sourcePAT = ""
	targetPAT = "target_token_def"
	os.Setenv("GITHUB_TOKEN", "github_token_xyz")

	sourceToken, targetToken, err = resolveTokens()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if sourceToken != "github_token_xyz" {
		t.Errorf("Expected source token 'github_token_xyz', got '%s'", sourceToken)
	}

	if targetToken != "target_token_def" {
		t.Errorf("Expected target token 'target_token_def', got '%s'", targetToken)
	}
}

// TestResolveTokens_NoTokensProvided tests error case when no tokens are available
func TestResolveTokens_NoTokensProvided(t *testing.T) {
	// Save original values
	origSourcePAT := sourcePAT
	origTargetPAT := targetPAT
	origGitHubToken := os.Getenv("GITHUB_TOKEN")
	
	// Clean up after test
	defer func() {
		sourcePAT = origSourcePAT
		targetPAT = origTargetPAT
		if origGitHubToken != "" {
			os.Setenv("GITHUB_TOKEN", origGitHubToken)
		} else {
			os.Unsetenv("GITHUB_TOKEN")
		}
	}()

	// Set test values - no tokens provided
	sourcePAT = ""
	targetPAT = ""
	os.Unsetenv("GITHUB_TOKEN")

	_, _, err := resolveTokens()
	if err == nil {
		t.Fatal("Expected error when no tokens provided, got nil")
	}

	expectedMsg := "authentication required"
	if len(err.Error()) < len(expectedMsg) || err.Error()[:len(expectedMsg)] != expectedMsg {
		t.Errorf("Expected error message to start with '%s', got '%s'", expectedMsg, err.Error())
	}
}

// TestResolveTokens_OnlySourcePATNoFallback tests error when only source PAT and no fallback
func TestResolveTokens_OnlySourcePATNoFallback(t *testing.T) {
	// Save original values
	origSourcePAT := sourcePAT
	origTargetPAT := targetPAT
	origGitHubToken := os.Getenv("GITHUB_TOKEN")
	
	// Clean up after test
	defer func() {
		sourcePAT = origSourcePAT
		targetPAT = origTargetPAT
		if origGitHubToken != "" {
			os.Setenv("GITHUB_TOKEN", origGitHubToken)
		} else {
			os.Unsetenv("GITHUB_TOKEN")
		}
	}()

	// Set test values - only source PAT, no GITHUB_TOKEN
	sourcePAT = "source_token_only"
	targetPAT = ""
	os.Unsetenv("GITHUB_TOKEN")

	_, _, err := resolveTokens()
	if err == nil {
		t.Fatal("Expected error when only source PAT provided without GITHUB_TOKEN, got nil")
	}
}
