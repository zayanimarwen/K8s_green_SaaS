package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/k8s-green/backend/internal/domain"
	"github.com/rs/zerolog/log"
)

// LiveWebSocket pousse le score en temps réel via SSE (Server-Sent Events).
func (h *Handler) LiveWebSocket(c *gin.Context) {
	clusterID := c.Param("id")
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	ctx := c.Request.Context()
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	h.pushScoreSSE(c, clusterID)

	for {
		select {
		case <-ctx.Done():
			log.Debug().Str("cluster", clusterID).Msg("SSE déconnecté")
			return
		case <-ticker.C:
			h.pushScoreSSE(c, clusterID)
		}
	}
}

func (h *Handler) pushScoreSSE(c *gin.Context, clusterID string) {
	score, err := h.db.TimescaleRepo().GetLatestScore(c.Request.Context(), clusterID)
	if err != nil {
		c.SSEvent("error", `{"message":"score indisponible"}`)
		c.Writer.Flush()
		return
	}
	data, _ := json.Marshal(score)
	c.SSEvent("score", string(data))
	c.Writer.Flush()
}

// GetScoreHistory retourne l'historique sur N jours
func (h *Handler) GetScoreHistory(c *gin.Context) {
	clusterID := c.Param("id")
	days := 30
	if d := c.Query("days"); d != "" {
		if n, err := strconv.Atoi(d); err == nil && n > 0 && n <= 365 {
			days = n
		}
	}

	scores, err := h.db.TimescaleRepo().GetScoreHistory(c.Request.Context(), clusterID, days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erreur historique"})
		return
	}
	var result interface{} = scores
	if scores == nil {
		result = []domain.GreenScore{}
	}
	c.JSON(http.StatusOK, gin.H{
		"cluster_id": clusterID,
		"days":       days,
		"history":    result,
	})
}
