package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/k8s-green/backend/internal/domain"
)

func (h *Handler) ListTenants(c *gin.Context) {
	tenants, err := h.db.TenantRepo().List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erreur récupération tenants"})
		return
	}
	if tenants == nil { tenants = []domain.Tenant{} }
	c.JSON(http.StatusOK, gin.H{"tenants": tenants})
}

func (h *Handler) CreateTenant(c *gin.Context) {
	var body struct {
		ID   string `json:"id"   binding:"required"`
		Name string `json:"name" binding:"required"`
		Plan string `json:"plan"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if body.Plan == "" { body.Plan = "starter" }

	tenant := &domain.Tenant{ID: body.ID, Name: body.Name, Plan: body.Plan, Active: true}
	if err := h.db.TenantRepo().Create(c.Request.Context(), tenant); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "création tenant échouée"})
		return
	}
	c.JSON(http.StatusCreated, tenant)
}

func (h *Handler) GetTenant(c *gin.Context) {
	tenantID := c.Param("id")
	tenant, err := h.db.TenantRepo().GetByID(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "tenant introuvable"})
		return
	}
	c.JSON(http.StatusOK, tenant)
}
