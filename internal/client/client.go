package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/renan-alm/gh-vars-migrator/internal/types"
)

// Client is a wrapper around the GitHub API client
type Client struct {
	restClient *api.RESTClient
}

// New creates a new GitHub API client using default authentication
func New() (*Client, error) {
	restClient, err := api.DefaultRESTClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create GitHub API client: %w", err)
	}

	return &Client{
		restClient: restClient,
	}, nil
}

// NewWithToken creates a new GitHub API client with an explicit token
func NewWithToken(token string) (*Client, error) {
	if token == "" {
		return nil, fmt.Errorf("token cannot be empty")
	}

	opts := api.ClientOptions{
		AuthToken: token,
	}

	restClient, err := api.NewRESTClient(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create GitHub API client with token: %w", err)
	}

	return &Client{
		restClient: restClient,
	}, nil
}

// ListRepoVariables lists all variables for a repository
func (c *Client) ListRepoVariables(owner, repo string) ([]types.Variable, error) {
	var response struct {
		Variables []types.Variable `json:"variables"`
	}

	path := fmt.Sprintf("repos/%s/%s/actions/variables", owner, repo)
	err := c.restClient.Get(path, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to list repository variables: %w", err)
	}

	return response.Variables, nil
}

// ListOrgVariables lists all variables for an organization
func (c *Client) ListOrgVariables(org string) ([]types.Variable, error) {
	var response struct {
		Variables []types.Variable `json:"variables"`
	}

	path := fmt.Sprintf("orgs/%s/actions/variables", org)
	err := c.restClient.Get(path, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to list organization variables: %w", err)
	}

	return response.Variables, nil
}

// ListEnvVariables lists all variables for a repository environment
func (c *Client) ListEnvVariables(owner, repo, env string) ([]types.Variable, error) {
	var response struct {
		Variables []types.Variable `json:"variables"`
	}

	path := fmt.Sprintf("repos/%s/%s/environments/%s/variables", owner, repo, env)
	err := c.restClient.Get(path, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to list environment variables: %w", err)
	}

	return response.Variables, nil
}

// GetRepoVariable gets a specific variable from a repository
func (c *Client) GetRepoVariable(owner, repo, name string) (*types.Variable, error) {
	var variable types.Variable

	path := fmt.Sprintf("repos/%s/%s/actions/variables/%s", owner, repo, name)
	err := c.restClient.Get(path, &variable)
	if err != nil {
		return nil, err
	}

	return &variable, nil
}

// GetOrgVariable gets a specific variable from an organization
func (c *Client) GetOrgVariable(org, name string) (*types.Variable, error) {
	var variable types.Variable

	path := fmt.Sprintf("orgs/%s/actions/variables/%s", org, name)
	err := c.restClient.Get(path, &variable)
	if err != nil {
		return nil, err
	}

	return &variable, nil
}

// GetEnvVariable gets a specific variable from an environment
func (c *Client) GetEnvVariable(owner, repo, env, name string) (*types.Variable, error) {
	var variable types.Variable

	path := fmt.Sprintf("repos/%s/%s/environments/%s/variables/%s", owner, repo, env, name)
	err := c.restClient.Get(path, &variable)
	if err != nil {
		return nil, err
	}

	return &variable, nil
}

// CreateRepoVariable creates a new variable in a repository
func (c *Client) CreateRepoVariable(owner, repo string, variable types.Variable) error {
	path := fmt.Sprintf("repos/%s/%s/actions/variables", owner, repo)
	body := map[string]string{
		"name":  variable.Name,
		"value": variable.Value,
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	err = c.restClient.Post(path, bytes.NewReader(bodyBytes), nil)
	if err != nil {
		return fmt.Errorf("failed to create repository variable: %w", err)
	}

	return nil
}

// CreateOrgVariable creates a new variable in an organization
func (c *Client) CreateOrgVariable(org string, variable types.Variable) error {
	path := fmt.Sprintf("orgs/%s/actions/variables", org)
	body := map[string]string{
		"name":       variable.Name,
		"value":      variable.Value,
		"visibility": "all",
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	err = c.restClient.Post(path, bytes.NewReader(bodyBytes), nil)
	if err != nil {
		return fmt.Errorf("failed to create organization variable: %w", err)
	}

	return nil
}

// CreateEnvVariable creates a new variable in an environment
func (c *Client) CreateEnvVariable(owner, repo, env string, variable types.Variable) error {
	path := fmt.Sprintf("repos/%s/%s/environments/%s/variables", owner, repo, env)
	body := map[string]string{
		"name":  variable.Name,
		"value": variable.Value,
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	err = c.restClient.Post(path, bytes.NewReader(bodyBytes), nil)
	if err != nil {
		return fmt.Errorf("failed to create environment variable: %w", err)
	}

	return nil
}

// UpdateRepoVariable updates an existing variable in a repository
func (c *Client) UpdateRepoVariable(owner, repo string, variable types.Variable) error {
	path := fmt.Sprintf("repos/%s/%s/actions/variables/%s", owner, repo, variable.Name)
	body := map[string]string{
		"name":  variable.Name,
		"value": variable.Value,
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	err = c.restClient.Patch(path, bytes.NewReader(bodyBytes), nil)
	if err != nil {
		return fmt.Errorf("failed to update repository variable: %w", err)
	}

	return nil
}

// UpdateOrgVariable updates an existing variable in an organization
func (c *Client) UpdateOrgVariable(org string, variable types.Variable) error {
	path := fmt.Sprintf("orgs/%s/actions/variables/%s", org, variable.Name)
	body := map[string]string{
		"name":       variable.Name,
		"value":      variable.Value,
		"visibility": "all",
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	err = c.restClient.Patch(path, bytes.NewReader(bodyBytes), nil)
	if err != nil {
		return fmt.Errorf("failed to update organization variable: %w", err)
	}

	return nil
}

// UpdateEnvVariable updates an existing variable in an environment
func (c *Client) UpdateEnvVariable(owner, repo, env string, variable types.Variable) error {
	path := fmt.Sprintf("repos/%s/%s/environments/%s/variables/%s", owner, repo, env, variable.Name)
	body := map[string]string{
		"name":  variable.Name,
		"value": variable.Value,
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	err = c.restClient.Patch(path, bytes.NewReader(bodyBytes), nil)
	if err != nil {
		return fmt.Errorf("failed to update environment variable: %w", err)
	}

	return nil
}

// ListEnvironments lists all environments for a repository
func (c *Client) ListEnvironments(owner, repo string) ([]types.Environment, error) {
	var response struct {
		TotalCount   int                 `json:"total_count"`
		Environments []types.Environment `json:"environments"`
	}

	path := fmt.Sprintf("repos/%s/%s/environments", owner, repo)
	err := c.restClient.Get(path, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to list environments: %w", err)
	}

	return response.Environments, nil
}

// GetEnvironment gets a specific environment from a repository
func (c *Client) GetEnvironment(owner, repo, envName string) (*types.Environment, error) {
	var env types.Environment

	path := fmt.Sprintf("repos/%s/%s/environments/%s", owner, repo, envName)
	err := c.restClient.Get(path, &env)
	if err != nil {
		return nil, err
	}

	return &env, nil
}

// CreateEnvironment creates a new environment in a repository
func (c *Client) CreateEnvironment(owner, repo, envName string) error {
	path := fmt.Sprintf("repos/%s/%s/environments/%s", owner, repo, envName)

	// GitHub API requires PUT with empty body to create an environment
	err := c.restClient.Put(path, bytes.NewReader([]byte("{}")), nil)
	if err != nil {
		return fmt.Errorf("failed to create environment: %w", err)
	}

	return nil
}

// GetTokenScopes returns the OAuth scopes associated with the token by inspecting
// the X-OAuth-Scopes response header. Returns nil if the header is absent (e.g.
// fine-grained PATs or GITHUB_TOKEN from Actions), indicating scope validation
// should be skipped.
func (c *Client) GetTokenScopes() ([]string, error) {
	resp, err := c.restClient.Request("GET", "user", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve token scopes: %w", err)
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)

	scopesHeader := resp.Header.Get("X-OAuth-Scopes")
	if scopesHeader == "" {
		return nil, nil
	}

	parts := strings.Split(scopesHeader, ",")
	scopes := make([]string, 0, len(parts))
	for _, s := range parts {
		if trimmed := strings.TrimSpace(s); trimmed != "" {
			scopes = append(scopes, trimmed)
		}
	}
	return scopes, nil
}

// GetUser retrieves the authenticated user information
func (c *Client) GetUser() (string, error) {
	var user struct {
		Login string `json:"login"`
	}

	if err := c.restClient.Get("user", &user); err != nil {
		return "", fmt.Errorf("failed to get user: %w", err)
	}

	return user.Login, nil
}
