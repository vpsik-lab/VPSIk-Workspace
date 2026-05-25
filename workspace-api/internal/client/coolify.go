package client

import (
	"context"
	"encoding/json"
	"fmt"
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

type CoolifyDeployment struct {
	ID        int    `json:"id"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}

func NewCoolifyClient(baseURL, apiToken string) *CoolifyClient {
	return &CoolifyClient{
		baseURL:  baseURL,
		apiToken: apiToken,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *CoolifyClient) CheckHealth(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/api/v1/health", nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 500 {
		return fmt.Errorf("coolify returned status %d", resp.StatusCode)
	}
	return nil
}

func (c *CoolifyClient) ListServers(ctx context.Context) ([]CoolifyServer, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/api/v1/servers", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("coolify api: %d", resp.StatusCode)
	}

	var servers []CoolifyServer
	if err := json.NewDecoder(resp.Body).Decode(&servers); err != nil {
		return nil, err
	}
	return servers, nil
}

func (c *CoolifyClient) ListDeployments(ctx context.Context) ([]CoolifyDeployment, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/api/v1/deployments", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("coolify api: %d", resp.StatusCode)
	}

	var deployments []CoolifyDeployment
	if err := json.NewDecoder(resp.Body).Decode(&deployments); err != nil {
		return nil, err
	}
	return deployments, nil
}
