package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/k8s-green/backend/internal/auth"
	"github.com/golang-jwt/jwt/v5"
)

const ContextTenantID = "tenant_id"

func Tenant() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := ""

		// 1. Depuis le claim JWT (source principale)
		if authHeader := c.GetHeader("Authorization"); strings.HasPrefix(authHeader, "Bearer ") {
			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			if token, _, err := jwt.NewParser().ParseUnverified(tokenStr, jwt.MapClaims{}); err == nil {
				if claims, ok := token.Claims.(jwt.MapClaims); ok {
					tenantID = auth.ExtractTenantID(claims)
				}
			}
		}

		// 2. Fallback: header X-Tenant-ID (pour les appels internes / tests)
		if tenantID == "" {
			tenantID = c.GetHeader("X-Tenant-ID")
		}

		if tenantID == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "tenant_id manquant — vérifier le claim JWT ou le header X-Tenant-ID",
			})
			return
		}

		c.Set(ContextTenantID, tenantID)
		c.Next()
	}
}

func TenantID(c *gin.Context) string {
	v, _ := c.Get(ContextTenantID)
	s, _ := v.(string)
	return s
}
