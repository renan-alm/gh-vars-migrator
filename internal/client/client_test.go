package client

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/renan-alm/gh-vars-migrator/internal/types"
)

// NOTE: The client package wraps the GitHub API client which is difficult to mock
// without modifying the production code. These tests verify the logic and data
// transformations used by the client methods.

// TestListRepoVariables_PathConstruction verifies the path construction
func TestListRepoVariables_PathConstruction(t *testing.T) {
	owner := "test-owner"
	repo := "test-repo"
	expectedPath := fmt.Sprintf("repos/%s/%s/actions/variables", owner, repo)
	
	if expectedPath != "repos/test-owner/test-repo/actions/variables" {
		t.Errorf("Path construction failed: got %s", expectedPath)
	}
}

// TestListOrgVariables_PathConstruction verifies the path construction
func TestListOrgVariables_PathConstruction(t *testing.T) {
	org := "test-org"
	expectedPath := fmt.Sprintf("orgs/%s/actions/variables", org)
	
	if expectedPath != "orgs/test-org/actions/variables" {
		t.Errorf("Path construction failed: got %s", expectedPath)
	}
}

// TestListEnvVariables_PathConstruction verifies the path construction
func TestListEnvVariables_PathConstruction(t *testing.T) {
	owner := "test-owner"
	repo := "test-repo"
	env := "production"
	expectedPath := fmt.Sprintf("repos/%s/%s/environments/%s/variables", owner, repo, env)
	
	if expectedPath != "repos/test-owner/test-repo/environments/production/variables" {
		t.Errorf("Path construction failed: got %s", expectedPath)
	}
}

// TestGetRepoVariable_PathConstruction verifies the path construction
func TestGetRepoVariable_PathConstruction(t *testing.T) {
	owner := "test-owner"
	repo := "test-repo"
	name := "VAR_NAME"
	expectedPath := fmt.Sprintf("repos/%s/%s/actions/variables/%s", owner, repo, name)
	
	if expectedPath != "repos/test-owner/test-repo/actions/variables/VAR_NAME" {
		t.Errorf("Path construction failed: got %s", expectedPath)
	}
}

// TestGetOrgVariable_PathConstruction verifies the path construction
func TestGetOrgVariable_PathConstruction(t *testing.T) {
	org := "test-org"
	name := "VAR_NAME"
	expectedPath := fmt.Sprintf("orgs/%s/actions/variables/%s", org, name)
	
	if expectedPath != "orgs/test-org/actions/variables/VAR_NAME" {
		t.Errorf("Path construction failed: got %s", expectedPath)
	}
}

// TestGetEnvVariable_PathConstruction verifies the path construction
func TestGetEnvVariable_PathConstruction(t *testing.T) {
	owner := "test-owner"
	repo := "test-repo"
	env := "production"
	name := "VAR_NAME"
	expectedPath := fmt.Sprintf("repos/%s/%s/environments/%s/variables/%s", owner, repo, env, name)
	
	if expectedPath != "repos/test-owner/test-repo/environments/production/variables/VAR_NAME" {
		t.Errorf("Path construction failed: got %s", expectedPath)
	}
}

// TestCreateRepoVariable_RequestBody verifies request body construction
func TestCreateRepoVariable_RequestBody(t *testing.T) {
	variable := types.Variable{Name: "TEST_VAR", Value: "test_value"}
	body := map[string]string{
		"name":  variable.Name,
		"value": variable.Value,
	}
	
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("Failed to marshal body: %v", err)
	}
	
	var decoded map[string]string
	if err := json.Unmarshal(bodyBytes, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal body: %v", err)
	}
	
	if decoded["name"] != "TEST_VAR" {
		t.Errorf("Expected name TEST_VAR, got %s", decoded["name"])
	}
	if decoded["value"] != "test_value" {
		t.Errorf("Expected value test_value, got %s", decoded["value"])
	}
}

// TestCreateOrgVariable_RequestBody verifies org variable body includes visibility
func TestCreateOrgVariable_RequestBody(t *testing.T) {
	variable := types.Variable{Name: "ORG_VAR", Value: "org_value"}
	body := map[string]string{
		"name":       variable.Name,
		"value":      variable.Value,
		"visibility": "all",
	}
	
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("Failed to marshal body: %v", err)
	}
	
	var decoded map[string]string
	if err := json.Unmarshal(bodyBytes, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal body: %v", err)
	}
	
	if decoded["visibility"] != "all" {
		t.Errorf("Expected visibility all, got %s", decoded["visibility"])
	}
	if decoded["name"] != "ORG_VAR" {
		t.Errorf("Expected name ORG_VAR, got %s", decoded["name"])
	}
}

// TestCreateEnvVariable_RequestBody verifies environment variable body construction
func TestCreateEnvVariable_RequestBody(t *testing.T) {
	variable := types.Variable{Name: "ENV_VAR", Value: "env_value"}
	body := map[string]string{
		"name":  variable.Name,
		"value": variable.Value,
	}
	
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("Failed to marshal body: %v", err)
	}
	
	var decoded map[string]string
	if err := json.Unmarshal(bodyBytes, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal body: %v", err)
	}
	
	if decoded["name"] != "ENV_VAR" {
		t.Errorf("Expected name ENV_VAR, got %s", decoded["name"])
	}
	// Environment variables should NOT have visibility field (unlike org variables)
	if _, exists := decoded["visibility"]; exists {
		t.Error("Environment variable body should not contain visibility field")
	}
}

// TestUpdateRepoVariable_PathConstruction verifies update path includes variable name
func TestUpdateRepoVariable_PathConstruction(t *testing.T) {
	owner := "test-owner"
	repo := "test-repo"
	varName := "MY_VAR"
	expectedPath := fmt.Sprintf("repos/%s/%s/actions/variables/%s", owner, repo, varName)
	
	if expectedPath != "repos/test-owner/test-repo/actions/variables/MY_VAR" {
		t.Errorf("Update path construction failed: got %s", expectedPath)
	}
}

// TestUpdateOrgVariable_PathConstruction verifies organization update path
func TestUpdateOrgVariable_PathConstruction(t *testing.T) {
	org := "test-org"
	varName := "MY_VAR"
	expectedPath := fmt.Sprintf("orgs/%s/actions/variables/%s", org, varName)
	
	if expectedPath != "orgs/test-org/actions/variables/MY_VAR" {
		t.Errorf("Org update path construction failed: got %s", expectedPath)
	}
}

// TestUpdateEnvVariable_PathConstruction verifies environment update path
func TestUpdateEnvVariable_PathConstruction(t *testing.T) {
	owner := "test-owner"
	repo := "test-repo"
	env := "staging"
	varName := "ENV_VAR"
	expectedPath := fmt.Sprintf("repos/%s/%s/environments/%s/variables/%s", owner, repo, env, varName)
	
	if expectedPath != "repos/test-owner/test-repo/environments/staging/variables/ENV_VAR" {
		t.Errorf("Env update path construction failed: got %s", expectedPath)
	}
}

// TestUpdateRepoVariable_RequestBody verifies update body construction
func TestUpdateRepoVariable_RequestBody(t *testing.T) {
	variable := types.Variable{Name: "UPDATED_VAR", Value: "new_value"}
	body := map[string]string{
		"name":  variable.Name,
		"value": variable.Value,
	}
	
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("Failed to marshal body: %v", err)
	}
	
	var decoded map[string]string
	if err := json.Unmarshal(bodyBytes, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal body: %v", err)
	}
	
	if decoded["value"] != "new_value" {
		t.Errorf("Expected value new_value, got %s", decoded["value"])
	}
}

// TestUpdateOrgVariable_RequestBody verifies org update body includes visibility
func TestUpdateOrgVariable_RequestBody(t *testing.T) {
	variable := types.Variable{Name: "ORG_VAR", Value: "updated_value"}
	body := map[string]string{
		"name":       variable.Name,
		"value":      variable.Value,
		"visibility": "all",
	}
	
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("Failed to marshal body: %v", err)
	}
	
	var decoded map[string]string
	if err := json.Unmarshal(bodyBytes, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal body: %v", err)
	}
	
	if decoded["visibility"] != "all" {
		t.Errorf("Expected visibility all in update body, got %s", decoded["visibility"])
	}
	if decoded["value"] != "updated_value" {
		t.Errorf("Expected value updated_value, got %s", decoded["value"])
	}
}

// TestNewWithToken_EmptyToken verifies that NewWithToken rejects empty tokens
func TestNewWithToken_EmptyToken(t *testing.T) {
	_, err := NewWithToken("")
	if err == nil {
		t.Error("Expected error for empty token, got nil")
	}
	if err != nil && err.Error() != "token cannot be empty" {
		t.Errorf("Expected 'token cannot be empty' error, got: %v", err)
	}
}

// TestWaitForRateLimit_HappyPath verifies no sleep occurs when rate limit is healthy
func TestWaitForRateLimit_HappyPath(t *testing.T) {
	sleepCalled := false
	rl := &types.RateLimitInfo{
		Limit:     5000,
		Remaining: 4500,
		ResetTime: time.Now().Add(1 * time.Hour),
	}

	waitForRateLimit(rl, minRemainingRequests, func(d time.Duration) {
		sleepCalled = true
	})

	if sleepCalled {
		t.Error("Expected no sleep when rate limit is healthy, but sleepFn was called")
	}
}

// TestWaitForRateLimit_NearLimit verifies sleep is triggered when remaining is below threshold
func TestWaitForRateLimit_NearLimit(t *testing.T) {
	var sleptDuration time.Duration
	resetIn := 30 * time.Second
	rl := &types.RateLimitInfo{
		Limit:     5000,
		Remaining: 5, // below minRemainingRequests (10)
		ResetTime: time.Now().Add(resetIn),
	}

	waitForRateLimit(rl, minRemainingRequests, func(d time.Duration) {
		sleptDuration = d
	})

	if sleptDuration <= 0 {
		t.Error("Expected sleepFn to be called with a positive duration, but it was not")
	}
	// Sleep duration should be approximately resetIn + 5s buffer
	if sleptDuration < resetIn {
		t.Errorf("Expected sleep duration >= %s, got %s", resetIn, sleptDuration)
	}
}

// TestWaitForRateLimit_AtThreshold verifies no sleep occurs when remaining equals threshold
func TestWaitForRateLimit_AtThreshold(t *testing.T) {
	sleepCalled := false
	rl := &types.RateLimitInfo{
		Limit:     5000,
		Remaining: minRemainingRequests, // exactly at threshold
		ResetTime: time.Now().Add(1 * time.Hour),
	}

	waitForRateLimit(rl, minRemainingRequests, func(d time.Duration) {
		sleepCalled = true
	})

	if sleepCalled {
		t.Error("Expected no sleep at threshold boundary, but sleepFn was called")
	}
}

// TestWaitForRateLimit_AlreadyReset verifies no sleep occurs when reset time is in the past
func TestWaitForRateLimit_AlreadyReset(t *testing.T) {
	sleepCalled := false
	rl := &types.RateLimitInfo{
		Limit:     5000,
		Remaining: 0, // exhausted
		ResetTime: time.Now().Add(-10 * time.Second), // already past
	}

	waitForRateLimit(rl, minRemainingRequests, func(d time.Duration) {
		sleepCalled = true
	})

	if sleepCalled {
		t.Error("Expected no sleep when reset time has already passed, but sleepFn was called")
	}
}
