package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type GiteaClient struct {
	baseURL   string
	token     string
	httpClient *http.Client
}

type GiteaRepo struct {
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	Private     bool   `json:"private"`
	HTMLURL     string `json:"html_url"`
	Language    string `json:"language"`
	Stars       int    `json:"stars_count"`
	Forks       int    `json:"forks_count"`
	UpdatedAt   string `json:"updated_at"`
}

type GiteaIssue struct {
	Number    int    `json:"number"`
	Title     string `json:"title"`
	State     string `json:"state"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func NewGiteaClient(baseURL, token string) *GiteaClient {
	return &GiteaClient{
		baseURL: baseURL,
		token:   token,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *GiteaClient) CheckHealth(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/api/healthz", nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 500 {
		return fmt.Errorf("gitea returned status %d", resp.StatusCode)
	}
	return nil
}

func (c *GiteaClient) ListRepos(ctx context.Context) ([]GiteaRepo, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/api/v1/user/repos", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "token "+c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("gitea api: %d", resp.StatusCode)
	}

	var repos []GiteaRepo
	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		return nil, err
	}
	return repos, nil
}

func (c *GiteaClient) ListIssues(ctx context.Context) ([]GiteaIssue, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/api/v1/user/issues", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "token "+c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("gitea api: %d", resp.StatusCode)
	}

	var issues []GiteaIssue
	if err := json.NewDecoder(resp.Body).Decode(&issues); err != nil {
		return nil, err
	}
	return issues, nil
}
