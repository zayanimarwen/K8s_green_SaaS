package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/k8s-green/backend/internal/auth"
)

func RequireRole(required string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roles := Roles(c)
		if !auth.HasRole(roles, required) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "accès refusé — rôle insuffisant",
				"required": required,
			})
			return
		}
		c.Next()
	}
}
