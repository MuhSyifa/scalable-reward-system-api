package cache

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
	"scalable-reward-system/internal/config"
)

func NewRedisClient(cfg *config.Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisHost + ":" + cfg.RedisPort,
		Password: cfg.RedisPass, // no password set
		DB:       0,             // use default DB
	})

	// Test connection
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	return client
}
