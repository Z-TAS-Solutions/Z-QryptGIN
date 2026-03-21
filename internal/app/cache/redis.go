package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(address, password string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         address,
		Password:     password,
		DB:           0,  // Use default DB
		PoolSize:     100, // Maximum number of socket connections
		MinIdleConns: 10,
	})

	// Ping the Redis server to ensure connection is alive
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return client, nil
}