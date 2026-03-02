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

type OpenAIClient struct {
	apiKey    string
	model     string
	maxTokens int
	httpClient *http.Client
}

func NewOpenAIClient(apiKey, model string, maxTokens int) *OpenAIClient {
	if maxTokens == 0 { maxTokens = 1000 }
	return &OpenAIClient{
		apiKey:     apiKey,
		model:      model,
		maxTokens:  maxTokens,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *OpenAIClient) BackendName() Backend { return BackendOpenAI }

func (c *OpenAIClient) Analyze(ctx context.Context, prompt string) (string, error) {
	body := map[string]interface{}{
		"model": c.model,
		"max_tokens": c.maxTokens,
		"messages": []map[string]string{
			{"role": "system", "content": systemPromptFR},
			{"role": "user",   "content": prompt},
		},
	}

	data, _ := json.Marshal(body)
	req, err := http.NewRequestWithContext(ctx, "POST",
		"https://api.openai.com/v1/chat/completions",
		bytes.NewReader(data))
	if err != nil { return "", err }
	req.Header.Set("Authorization", "Bearer " + c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil { return "", fmt.Errorf("openai request: %w", err) }
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("openai error %d: %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Choices []struct {
			Message struct { Content string `json:"content"` } `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil { return "", err }
	if len(result.Choices) == 0 { return "", fmt.Errorf("no choices in response") }
	return result.Choices[0].Message.Content, nil
}
