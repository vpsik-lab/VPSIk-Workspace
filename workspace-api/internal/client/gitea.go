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

type GiteaWebhook struct {
	ID     int      `json:"id"`
	URL    string   `json:"url"`
	Active bool     `json:"active"`
	Events []string `json:"events"`
	Type   string   `json:"type"`
}

func (c *GiteaClient) CreateWebhook(ctx context.Context, repo, webhookURL, secret string, events []string) (*GiteaWebhook, error) {
	body := map[string]interface{}{
		"type":   "gitea",
		"config": map[string]string{"url": webhookURL, "content_type": "json", "secret": secret},
		"events": events,
		"active": true,
	}
	var hook GiteaWebhook
	err := c.doJSON(ctx, http.MethodPost, fmt.Sprintf("/api/v1/repos/%s/hooks", repo), body, &hook)
	return &hook, err
}

func (c *GiteaClient) ListWebhooks(ctx context.Context, repo string) ([]GiteaWebhook, error) {
	var hooks []GiteaWebhook
	err := c.doJSON(ctx, http.MethodGet, fmt.Sprintf("/api/v1/repos/%s/hooks", repo), nil, &hooks)
	return hooks, err
}

func (c *GiteaClient) doJSON(ctx context.Context, method, path string, body, out interface{}) error {
	var buf io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return err
		}
		buf = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, buf)
	if err != nil {
		return err
	}
	if c.token != "" {
		req.Header.Set("Authorization", "token "+c.token)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("gitea api %s: %d - %s", path, resp.StatusCode, string(respBody))
	}

	if out != nil {
		return json.NewDecoder(resp.Body).Decode(out)
	}
	return nil
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
