package redis

import (
	"github.com/redis/go-redis/v9"
)

func init() {
	key := GetKey()
	client := redis.NewClient(&redis.Options{
		Addr:     key.Addr,
		Username: key.Username,
		Password: key.Password,
		DB:       key.DB,
	})

}
