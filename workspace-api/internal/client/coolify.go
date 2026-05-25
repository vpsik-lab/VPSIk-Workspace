package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type CoolifyClient struct {
	baseURL    string
	apiToken   string
	httpClient *http.Client
}

type CoolifyServer struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	IP        string `json:"ip"`
	Port      int    `json:"port"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}

type CoolifyProject struct {
	ID          string `json:"id"`
	UUID        string `json:"uuid"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CoolifyEnvironment struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}

type CoolifyApplication struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	UUID        string `json:"uuid"`
	FQDN        string `json:"fqdn"`
	Status      string `json:"status"`
	GitSource   string `json:"git_source"`
	RepoName    string `json:"repository"`
	Branch      string `json:"git_branch"`
	UpdatedAt   string `json:"updated_at"`
}

type CoolifyDeployment struct {
	ID        string `json:"id"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	Logs      string `json:"logs,omitempty"`
}

func NewCoolifyClient(baseURL, apiToken string) *CoolifyClient {
	return &CoolifyClient{
		baseURL:  baseURL,
		apiToken: apiToken,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *CoolifyClient) do(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
	var buf io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		buf = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiToken)
	req.Header.Set("Content-Type", "application/json")
	return c.httpClient.Do(req)
}

func (c *CoolifyClient) doJSON(ctx context.Context, method, path string, body, out interface{}) error {
	resp, err := c.do(ctx, method, path, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("coolify api %s %s: %d - %s", method, path, resp.StatusCode, string(respBody))
	}

	if out != nil {
		return json.NewDecoder(resp.Body).Decode(out)
	}
	return nil
}

func (c *CoolifyClient) CheckHealth(ctx context.Context) error {
	resp, err := c.do(ctx, http.MethodGet, "/api/v1/health", nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	if resp.StatusCode >= 500 {
		return fmt.Errorf("coolify returned status %d", resp.StatusCode)
	}
	return nil
}

func (c *CoolifyClient) ListServers(ctx context.Context) ([]CoolifyServer, error) {
	var servers []CoolifyServer
	err := c.doJSON(ctx, http.MethodGet, "/api/v1/servers", nil, &servers)
	return servers, err
}

func (c *CoolifyClient) ListDeployments(ctx context.Context) ([]CoolifyDeployment, error) {
	var deployments []CoolifyDeployment
	err := c.doJSON(ctx, http.MethodGet, "/api/v1/deployments", nil, &deployments)
	return deployments, err
}

func (c *CoolifyClient) ListProjects(ctx context.Context) ([]CoolifyProject, error) {
	var projects []CoolifyProject
	err := c.doJSON(ctx, http.MethodGet, "/api/v1/projects", nil, &projects)
	return projects, err
}

func (c *CoolifyClient) ListEnvironments(ctx context.Context, projectUUID string) ([]CoolifyEnvironment, error) {
	var envs []CoolifyEnvironment
	err := c.doJSON(ctx, http.MethodGet, fmt.Sprintf("/api/v1/projects/%s/environments", projectUUID), nil, &envs)
	return envs, err
}

func (c *CoolifyClient) ListApplications(ctx context.Context, projectUUID, envName string) ([]CoolifyApplication, error) {
	var apps []CoolifyApplication
	err := c.doJSON(ctx, http.MethodGet, fmt.Sprintf("/api/v1/projects/%s/%s/applications", projectUUID, envName), nil, &apps)
	return apps, err
}

func (c *CoolifyClient) Deploy(ctx context.Context, resourceUUID string) error {
	return c.doJSON(ctx, http.MethodPost, fmt.Sprintf("/api/v1/deploy?tag=%s&force=false", resourceUUID), nil, nil)
}

func (c *CoolifyClient) Restart(ctx context.Context, resourceUUID string) error {
	return c.doJSON(ctx, http.MethodGet, fmt.Sprintf("/api/v1/restart?tag=%s", resourceUUID), nil, nil)
}

func (c *CoolifyClient) GetDeploymentLogs(ctx context.Context, deploymentID string) (string, error) {
	var result struct {
		Logs string `json:"logs"`
	}
	err := c.doJSON(ctx, http.MethodGet, fmt.Sprintf("/api/v1/deployments/%s", deploymentID), nil, &result)
	if err != nil {
		return "", err
	}
	return result.Logs, nil
}

func (c *CoolifyClient) GetEnvVars(ctx context.Context, projectUUID, envName, appUUID string) (map[string]string, error) {
	var envVars map[string]string
	err := c.doJSON(ctx, http.MethodGet, fmt.Sprintf("/api/v1/projects/%s/%s/applications/%s/env", projectUUID, envName, appUUID), nil, &envVars)
	return envVars, err
}

func (c *CoolifyClient) UpdateEnvVars(ctx context.Context, projectUUID, envName, appUUID string, envVars map[string]string) error {
	return c.doJSON(ctx, http.MethodPost, fmt.Sprintf("/api/v1/projects/%s/%s/applications/%s/env", projectUUID, envName, appUUID), envVars, nil)
}
