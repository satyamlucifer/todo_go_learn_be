package config

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var RedisCtx = context.Background()

func InitRedis() *redis.Client {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"), // e.g. "localhost:6379"
		Password: "",                     // no password by default
		DB:       0,                      // use default DB
	})

	_, err := RedisClient.Ping(RedisCtx).Result()
	if err != nil {
		log.Fatalf("❌ Failed to connect to Redis: %v", err)
	}

	log.Println("✅ Connected to Redis")
	return RedisClient
}
