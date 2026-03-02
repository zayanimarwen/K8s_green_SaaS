package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type ClaudeClient struct {
	apiKey     string
	model      string
	maxTokens  int
	httpClient *http.Client
}

func NewClaudeClient(apiKey, model string, maxTokens int) *ClaudeClient {
	if maxTokens == 0 { maxTokens = 1000 }
	return &ClaudeClient{
		apiKey:     apiKey,
		model:      model,
		maxTokens:  maxTokens,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *ClaudeClient) BackendName() Backend { return BackendClaude }

func (c *ClaudeClient) Analyze(ctx context.Context, prompt string) (string, error) {
	body := map[string]interface{}{
		"model":      c.model,
		"max_tokens": c.maxTokens,
		"system":     systemPromptFR,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
	}

	data, _ := json.Marshal(body)
	req, err := http.NewRequestWithContext(ctx, "POST",
		"https://api.anthropic.com/v1/messages",
		bytes.NewReader(data))
	if err != nil { return "", err }
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil { return "", fmt.Errorf("claude request: %w", err) }
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("claude error %d: %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil { return "", err }
	for _, block := range result.Content {
		if block.Type == "text" { return block.Text, nil }
	}
	return "", fmt.Errorf("no text in claude response")
}
