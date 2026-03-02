package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/k8s-green/backend/internal/ai"
)

// ListModels retourne les modeles Ollama disponibles + catalogue
func (h *Handler) ListModels(c *gin.Context) {
	client, ok := h.getOllamaClient()
	if !ok {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error":   "Ollama non configure",
			"hint":    "Definir OLLAMA_URL et AI_BACKEND=ollama dans les secrets",
			"models":  []interface{}{},
		})
		return
	}

	models, err := client.ListModels(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"current_model": h.cfg.AIModel,
		"ollama_url":    h.cfg.OllamaURL,
		"models":        models,
	})
}

// PullModel telecharge un modele Ollama
func (h *Handler) PullModel(c *gin.Context) {
	client, ok := h.getOllamaClient()
	if !ok {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Ollama non configure"})
		return
	}

	var body struct {
		Model string `json:"model" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Pull asynchrone ??? peut prendre plusieurs minutes
	go func() {
		if err := client.PullModel(c.Request.Context(), body.Model); err != nil {
			// Log uniquement, le client recevra la reponse 202
			_ = err
		}
	}()

	c.JSON(http.StatusAccepted, gin.H{
		"message": "Telechargement du modele " + body.Model + " lance en arriere-plan",
		"model":   body.Model,
		"status":  "pulling",
	})
}

// DeleteModel supprime un modele Ollama
func (h *Handler) DeleteModel(c *gin.Context) {
	client, ok := h.getOllamaClient()
	if !ok {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Ollama non configure"})
		return
	}

	modelName := c.Param("name")
	if modelName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "nom du modele requis"})
		return
	}

	if err := client.DeleteModel(c.Request.Context(), modelName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Modele " + modelName + " supprime",
		"model":   modelName,
	})
}

// SwitchModel change le modele actif pour les diagnostics
func (h *Handler) SwitchModel(c *gin.Context) {
	var body struct {
		Model string `json:"model" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Mettre a jour le modele actif en Redis (persistant jusqu au redemarrage)
	key := "ai:active_model"
	h.rdb.Client().Set(c.Request.Context(), key, body.Model, 0)

	// Mettre a jour le modele dans le service IA
	if h.aiSvc != nil {
		h.aiSvc.SetModel(body.Model)
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Modele actif change",
		"previous_model": h.cfg.AIModel,
		"active_model":  body.Model,
	})
}

// helper interne
func (h *Handler) getOllamaClient() (*ai.OllamaClient, bool) {
	if h.aiSvc == nil { return nil, false }
	return h.aiSvc.OllamaClient(), true
}
