package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/k8s-green/backend/internal/api/middleware"
	"github.com/k8s-green/backend/internal/domain"
)

func (h *Handler) GetScore(c *gin.Context) {
	clusterID := c.Param("id")
	tenantID  := middleware.TenantID(c)
	cacheKey  := fmt.Sprintf("tenant:%s:score:%s", tenantID, clusterID)
	ctx       := c.Request.Context()

	// 1. Cache Redis (TTL 30s)
	if cached, err := h.rdb.Client().Get(ctx, cacheKey).Bytes(); err == nil {
		var score domain.GreenScore
		if json.Unmarshal(cached, &score) == nil {
			c.Header("X-Cache", "HIT")
			c.JSON(http.StatusOK, score)
			return
		}
	}

	// 2. Fallback TimescaleDB
	score, err := h.db.TimescaleRepo().GetLatestScore(ctx, clusterID)
	if err != nil {
		// Aucun score encore : cluster trop récent
		c.JSON(http.StatusOK, gin.H{
			"cluster_id": clusterID,
			"score": 0, "grade": "N/A",
			"message": "En attente du premier scan agent",
		})
		return
	}

	// 3. Mettre en cache
	if b, err := json.Marshal(score); err == nil {
		h.rdb.Client().Set(context.Background(), cacheKey, b, 30*time.Second)
	}

	c.Header("X-Cache", "MISS")
	c.JSON(http.StatusOK, score)
}
