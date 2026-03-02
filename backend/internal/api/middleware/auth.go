package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/k8s-green/backend/internal/auth"
	"github.com/k8s-green/backend/internal/config"
)

const (
	ContextUserID    = "user_id"
	ContextUserEmail = "user_email"
	ContextRoles     = "roles"
)

func Auth(cfg *config.Config) gin.HandlerFunc {
	var validator *auth.JWTValidator
	if cfg.Env != "development" && cfg.JWTIssuer != "" {
		validator = auth.NewJWTValidator(cfg.JWTIssuer, cfg.KeycloakClientID)
	}

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token manquant"})
			return
		}
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		var claims jwt.MapClaims

		if cfg.Env == "development" || validator == nil {
			// Dev: ParseUnverified — facilite les tests sans Keycloak
			token, _, err := jwt.NewParser().ParseUnverified(tokenStr, jwt.MapClaims{})
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token invalide"})
				return
			}
			var ok bool
			claims, ok = token.Claims.(jwt.MapClaims)
			if !ok {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "claims invalides"})
				return
			}
		} else {
			// Production: validation RS256 + expiration via JWKS Keycloak
			var err error
			claims, err = validator.Validate(tokenStr)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token invalide ou expiré"})
				return
			}
		}

		// Extraire sub, email, roles depuis les claims Keycloak
		sub, _ := claims["sub"].(string)
		email, _ := claims["email"].(string)
		roles := auth.ExtractRoles(claims["realm_access"])

		c.Set(ContextUserID, sub)
		c.Set(ContextUserEmail, email)
		c.Set(ContextRoles, roles)
		c.Next()
	}
}

// Helpers pour récupérer les valeurs du context depuis les handlers

func UserID(c *gin.Context) string {
	v, _ := c.Get(ContextUserID)
	s, _ := v.(string)
	return s
}

func UserEmail(c *gin.Context) string {
	v, _ := c.Get(ContextUserEmail)
	s, _ := v.(string)
	return s
}

func Roles(c *gin.Context) []string {
	v, _ := c.Get(ContextRoles)
	r, _ := v.([]string)
	return r
}
