package db

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/substrate-cli/api-server/cmd/app/mq"
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

func SaveRedis(key string, value string) {
	redisClient = mq.GetRedisConnection()
	ctx := context.Background()

	log.Println("Saving cluster to redis, ", key)

	err := redisClient.Set(ctx, key, value, 0).Err()
	if err != nil {
		log.Println("Unable to save value")
	}
}

func DeleteRedisKey(key string) {
	ctx := context.Background() // you can pass a context from caller if needed
	rdb := mq.GetRedisConnection()
	result, err := rdb.Del(ctx, key).Result()
	if err != nil {
		log.Println("error in deleting key")
	}
	if result == 0 {
		log.Println("no task in processing")
	} else {
		log.Println("running tasks cleaned.")
	}
}
