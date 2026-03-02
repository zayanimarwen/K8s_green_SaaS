package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

// ModelInfo decrit un modele Ollama disponible
type ModelInfo struct {
	Name        string `json:"name"`
	Size        int64  `json:"size"`
	Description string `json:"description"`
	Pulled      bool   `json:"pulled"`
}

// Catalogue des modeles recommandes avec descriptions
var RecommendedModels = []ModelInfo{
	{Name: "llama3.2",      Description: "Diagnostic rapide - recommande pour la plupart des cas"},
	{Name: "llama3.1:8b",   Description: "Analyse plus precise - bon equilibre vitesse/qualite"},
	{Name: "mistral",       Description: "Excellent pour l analyse de logs et evenements K8s"},
	{Name: "codellama",     Description: "Specialise YAML/code - ideal pour les suggestions kubectl"},
	{Name: "phi3:mini",     Description: "Tres leger (2 Go) - pour machines avec peu de RAM"},
	{Name: "gemma2:2b",     Description: "Google Gemma 2 - rapide et efficace"},
}

type OllamaClient struct {
	baseURL    string
	model      string
	maxTokens  int
	httpClient *http.Client
}

func NewOllamaClient(baseURL, model string, maxTokens int) *OllamaClient {
	if maxTokens == 0 { maxTokens = 2000 }
	if baseURL == ""  { baseURL = "http://localhost:11434" }
	if model == ""    { model = "llama3.2" }
	return &OllamaClient{
		baseURL:    baseURL,
		model:      model,
		maxTokens:  maxTokens,
		httpClient: &http.Client{Timeout: 120 * time.Second},
	}
}

func (c *OllamaClient) BackendName() Backend { return BackendOllama }
func (c *OllamaClient) ModelName() string    { return c.model }

func (c *OllamaClient) Analyze(ctx context.Context, prompt string) (string, error) {
	body := map[string]interface{}{
		"model":  c.model,
		"prompt": systemPromptFR + "\n\n" + prompt,
		"stream": false,
		"options": map[string]interface{}{
			"num_predict": c.maxTokens,
			"temperature": 0.3,
			"top_p":       0.9,
		},
	}
	data, _ := json.Marshal(body)

	req, err := http.NewRequestWithContext(ctx, "POST",
		c.baseURL+"/api/generate", bytes.NewReader(data))
	if err != nil { return "", err }
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("ollama non joignable sur %s: %w", c.baseURL, err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("ollama erreur %d: %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Response string `json:"response"`
		Done     bool   `json:"done"`
		DoneReason string `json:"done_reason"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return string(respBody), nil
	}
	return result.Response, nil
}

// ListModels retourne les modeles installes + le catalogue recommande
func (c *OllamaClient) ListModels(ctx context.Context) ([]ModelInfo, error) {
	req, _ := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/api/tags", nil)
	resp, err := c.httpClient.Do(req)
	if err != nil { return nil, fmt.Errorf("ollama non joignable: %w", err) }
	defer resp.Body.Close()

	var tags struct {
		Models []struct {
			Name  string `json:"name"`
			Size  int64  `json:"size"`
		} `json:"models"`
	}
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &tags)

	// Construire la liste avec statut pulled
	pulledMap := map[string]int64{}
	for _, m := range tags.Models {
		pulledMap[m.Name] = m.Size
	}

	var result []ModelInfo
	for _, rec := range RecommendedModels {
		info := rec
		if size, ok := pulledMap[rec.Name]; ok {
			info.Pulled = true
			info.Size = size
		} else if size, ok := pulledMap[rec.Name+":latest"]; ok {
			info.Pulled = true
			info.Size = size
		}
		result = append(result, info)
	}

	// Ajouter les modeles installes non presents dans le catalogue
	for name, size := range pulledMap {
		inCatalogue := false
		for _, rec := range RecommendedModels {
			if rec.Name == name || rec.Name+":latest" == name { inCatalogue = true; break }
		}
		if !inCatalogue {
			result = append(result, ModelInfo{
				Name: name, Size: size, Pulled: true,
				Description: "Modele installe",
			})
		}
	}
	return result, nil
}

// PullModel telecharge un modele depuis le catalogue Ollama
func (c *OllamaClient) PullModel(ctx context.Context, modelName string) error {
	log.Info().Str("model", modelName).Msg("Pull modele Ollama...")
	body := map[string]interface{}{"name": modelName, "stream": false}
	data, _ := json.Marshal(body)

	client := &http.Client{Timeout: 30 * time.Minute}
	req, err := http.NewRequestWithContext(ctx, "POST",
		c.baseURL+"/api/pull", bytes.NewReader(data))
	if err != nil { return err }
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil { return fmt.Errorf("pull model: %w", err) }
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("pull erreur %d: %s", resp.StatusCode, string(b))
	}
	log.Info().Str("model", modelName).Msg("Modele pret")
	return nil
}

// DeleteModel supprime un modele pour liberer de l espace
func (c *OllamaClient) DeleteModel(ctx context.Context, modelName string) error {
	body := map[string]string{"name": modelName}
	data, _ := json.Marshal(body)
	req, _ := http.NewRequestWithContext(ctx, "DELETE",
		c.baseURL+"/api/delete", bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil { return err }
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("delete erreur %d: %s", resp.StatusCode, string(b))
	}
	return nil
}

func (c *OllamaClient) IsAvailable(ctx context.Context) bool {
	req, _ := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/api/tags", nil)
	resp, err := c.httpClient.Do(req)
	if err != nil { return false }
	io.ReadAll(resp.Body)
	resp.Body.Close()
	return resp.StatusCode == 200
}
