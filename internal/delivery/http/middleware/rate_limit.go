package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"scalable-reward-system/internal/utils"
)

// RateLimiter uses Redis to limit requests per IP
func RateLimiter(redisClient *redis.Client, requests int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		key := "rate_limit:" + ip

		ctx := context.Background()
		
		// Increment request count
		count, err := redisClient.Incr(ctx, key).Result()
		if err != nil {
			// If redis fails, we let request pass to not break app entirely, but log it
			c.Next()
			return
		}

		// If first request, set expiry
		if count == 1 {
			redisClient.Expire(ctx, key, window)
		}

		if count > int64(requests) {
			utils.ErrorResponse(c, http.StatusTooManyRequests, "Too many requests, please try again later.")
			return
		}

		c.Next()
	}
}
