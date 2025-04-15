package configs

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// InitRedis initializes the Redis client
func (config *RedisConfig) InitRedis() (*redis.Client, error) {
	// Create Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Host, strconv.Itoa(config.Port)),
		Password: config.Password,
		PoolSize: config.PoolSize,
	})

	// Set connection timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.ConnTTL)*time.Second)
	defer cancel()

	// Test Redis connection
	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("❌ Redis connection failed: %w", err)
	}

	// Assign Redis client to config
	Logger.Info().Msg("✅ Redis connected successfully!")

	return redisClient, nil
}

// ShutdownRedis closes the Redis connection safely
func ShutdownRedis(redisClient *redis.Client) {
	if redisClient == nil {
		Logger.Info().Msg("⚠️ Redis client is nil, skipping shutdown.")
		return
	}

	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := redisClient.Close(); err != nil {
		Logger.Error().Msgf("❌ Error closing Redis connection: %v", err)
	} else {
		Logger.Info().Msg("✅ Redis connection closed successfully!")
	}
}
