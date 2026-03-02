package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/k8s-green/backend/internal/api/middleware"
	"github.com/k8s-green/backend/internal/domain"
)

func (h *Handler) ListClusters(c *gin.Context) {
	tenantID := middleware.TenantID(c)
	clusters, err := h.db.ClusterRepo().ListByTenant(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erreur récupération clusters"})
		return
	}
	if clusters == nil {
		clusters = []domain.Cluster{}
	}
	c.JSON(http.StatusOK, gin.H{"clusters": clusters})
}

func (h *Handler) CreateCluster(c *gin.Context) {
	tenantID := middleware.TenantID(c)
	var body struct {
		Name        string `json:"name"        binding:"required"`
		Provider    string `json:"provider"    binding:"required"`
		Region      string `json:"region"      binding:"required"`
		Environment string `json:"environment"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	env := body.Environment
	if env == "" { env = "production" }

	cluster := &domain.Cluster{
		TenantID:    tenantID,
		Name:        body.Name,
		Provider:    body.Provider,
		Region:      body.Region,
		Environment: env,
		CreatedAt:   time.Now(),
	}
	if err := h.db.ClusterRepo().Create(c.Request.Context(), cluster); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "création cluster échouée"})
		return
	}
	c.JSON(http.StatusCreated, cluster)
}
