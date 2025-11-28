package database

import (
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

func NewConnectionRedis() *redis.Client {
	url := os.Getenv("REDIS_URL")
	port := os.Getenv("REDIS_PORT")

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", url, port),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return rdb
}
