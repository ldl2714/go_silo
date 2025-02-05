package redis

import (
	"context"
	"log"
	"sync"

	"github.com/go-redis/redis/v8"
)

var (
	ctx         = context.Background()
	redisClient *redis.Client
	once        sync.Once
)

func GetRedisClient() *redis.Client {
	once.Do(func() {
		redisClient = redis.NewClient(&redis.Options{
			Addr:     "localhost:6379", // Redis server address
			Password: "",               // No password set
			DB:       0,                // Use default DB
		})

		// Test the connection
		_, err := redisClient.Ping(ctx).Result()
		if err != nil {
			log.Fatalf("Could not connect to Redis: %v", err)
		} else {
			log.Println("Redis连接成功")
		}
	})

	return redisClient
}
