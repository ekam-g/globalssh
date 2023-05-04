package redis

import (
	"github.com/redis/go-redis/v9"
)

func init() {

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Username: "",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

}
