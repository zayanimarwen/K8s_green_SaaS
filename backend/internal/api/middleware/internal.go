package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/k8s-green/backend/internal/config"
)

// InternalKey valide les appels internes de l'agent via X-Internal-Key ou Bearer = signing key
func InternalKey(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Accepter X-Internal-Key header
		key := c.GetHeader("X-Internal-Key")
		if key == "" {
			// Fallback: Bearer token = signing key
			auth := c.GetHeader("Authorization")
			key = strings.TrimPrefix(auth, "Bearer ")
		}

		if key == "" || key != cfg.SigningKey {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "clé interne invalide"})
			return
		}

		// Extraire tenant_id du header ou body
		tenantID := c.GetHeader("X-Tenant-ID")
		if tenantID == "" {
			tenantID = "tenant-demo"
		}
		c.Set(ContextTenantID, tenantID)
		c.Next()
	}
}
