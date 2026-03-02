package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/k8s-green/backend/internal/repository"
	"github.com/rs/zerolog/log"
)

func Audit(db *repository.Postgres) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		// Persistance asynchrone pour ne pas alourdir la réponse
		go func() {
			entry := repository.AuditEntry{
				TenantID:   TenantID(c),
				UserID:     UserID(c),
				UserEmail:  UserEmail(c),
				Action:     c.Request.Method + " " + c.FullPath(),
				Resource:   c.FullPath(),
				IPAddress:  c.ClientIP(),
				UserAgent:  c.Request.UserAgent(),
				StatusCode: c.Writer.Status(),
				DurationMs: time.Since(start).Milliseconds(),
			}
			if err := db.AuditRepo().Insert(c.Request.Context(), entry); err != nil {
				log.Debug().Err(err).Msg("audit insert failed")
			}
		}()
	}
}

// Logger retourne un middleware de logging structuré zerolog
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		log.Info().
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Int("status", c.Writer.Status()).
			Dur("duration", time.Since(start)).
			Str("tenant", TenantID(c)).
			Str("ip", c.ClientIP()).
			Msg("request")
	}
}
