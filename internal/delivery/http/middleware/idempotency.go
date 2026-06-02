package middleware

import (
	"bytes"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

func Idempotency(redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		idempotencyKey := c.GetHeader("Idempotency-Key")
		if idempotencyKey == "" {
			// If no key provided, just proceed normally (or return error if strict)
			c.Next()
			return
		}

		ctx := context.Background()
		cacheKey := "idempotency:" + idempotencyKey

		// Check if response is already cached
		cachedResp, err := redisClient.Get(ctx, cacheKey).Result()
		if err == nil && cachedResp != "" {
			// Return cached response
			c.Data(http.StatusOK, "application/json", []byte(cachedResp))
			c.Abort()
			return
		}

		// Intercept the response
		w := &responseBodyWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = w

		c.Next()

		// Cache successful responses
		if c.Writer.Status() >= 200 && c.Writer.Status() < 300 {
			redisClient.Set(ctx, cacheKey, w.body.String(), 24*time.Hour)
		}
	}
}
