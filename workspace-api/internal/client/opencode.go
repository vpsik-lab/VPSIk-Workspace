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

type OpenCodeClient struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

type OpenCodeRepo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Language    string `json:"language"`
}

type OpenCodeChatRequest struct {
	Message  string `json:"message"`
	Context  string `json:"context,omitempty"`
	RepoPath string `json:"repo_path,omitempty"`
}

type OpenCodeChatResponse struct {
	Reply string `json:"reply"`
}

func NewOpenCodeClient(baseURL, apiKey string) *OpenCodeClient {
	if baseURL == "" {
		return nil
	}
	return &OpenCodeClient{
		baseURL:    baseURL,
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 120 * time.Second},
	}
}

func (c *OpenCodeClient) CheckHealth(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/api/health", nil)
	if err != nil {
		return err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (c *OpenCodeClient) Chat(ctx context.Context, message, contextInfo, repoPath string) (string, error) {
	body := OpenCodeChatRequest{
		Message:  message,
		Context:  contextInfo,
		RepoPath: repoPath,
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/api/chat", bytes.NewReader(payload))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("opencode api: %d - %s", resp.StatusCode, string(body))
	}

	var chatResp OpenCodeChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return "", err
	}

	return chatResp.Reply, nil
}
