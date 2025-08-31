package db

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/sshfz/api-server-substrate/cmd/app/mq"
)

var redisClient *redis.Client

func ReadValueFromKey(key string) (string, error) {
	redisClient := mq.GetRedisConnection()
	ctx := context.Background()

	val, err := redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", err
	} else if err != nil {
		log.Println("Error reading value from redis:", err)
		return "", err
	}
	return val, nil
}
