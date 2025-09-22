package connections

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/substrate-cli/api-server/cmd/app/mq"
	"github.com/substrate-cli/api-server/internal/utils"
)

var ctx = context.Background()
var rdb *redis.Client

func InitRedis() {
	log.Println("Initialising redis....")
	rdb = redis.NewClient(&redis.Options{
		Addr: utils.GetRedisAddr(),
	})
	err := rdb.Set(ctx, "test_key", "hello redis", 0).Err()
	if err != nil {
		log.Println("Connection to be redis cannot be established, quitting server...")
		log.Fatal(err)
	}
	log.Println("connection to redis successfully established")
	mq.SetRedisConnection(rdb)
}
