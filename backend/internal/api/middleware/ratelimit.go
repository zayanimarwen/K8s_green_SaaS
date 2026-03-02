package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/k8s-green/backend/internal/repository"
)

const rateLimit = 1000 // req/min/tenant

func RateLimit(rdb *repository.Redis) gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := TenantID(c)
		if tenantID == "" {
			c.Next()
			return
		}

		minute := time.Now().Format("200601021504") // YYYYMMDDHHmm
		key := fmt.Sprintf("rl:tenant:%s:%s", tenantID, minute)
		ctx := context.Background()

		// Increment + TTL atomique
		pipe := rdb.Client().Pipeline()
		incr := pipe.Incr(ctx, key)
		pipe.Expire(ctx, key, 2*time.Minute) // TTL 2min pour couvrir la fenêtre
		pipe.Exec(ctx)

		count := incr.Val()

		c.Header("X-RateLimit-Limit", strconv.Itoa(rateLimit))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(max(0, rateLimit-int(count))))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(time.Minute).Unix(), 10))

		if count > int64(rateLimit) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":       "rate limit dépassé",
				"limit":       rateLimit,
				"retry_after": 60,
			})
			return
		}
		c.Next()
	}
}

func max(a, b int) int {
	if a > b { return a }
	return b
}
