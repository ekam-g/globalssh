package redis

import (
	"log"

	"github.com/redis/go-redis/v9"
)

const (
	hostKey   string = "&&**$$&*#$"
	clientKey string = "@@$$##*()$"
)

var (
	hostMode       bool = true
	HostName       string
	client         *redis.Client
	command_stream *redis.PubSub
	result_stream  *redis.PubSub
)

func GetWriteKey() string {
	if hostMode {
		return clientKey
	} else {
		return hostKey
	}
}

func GetReadKey() string {
	if !hostMode {
		return clientKey
	} else {
		return hostKey
	}
}

func Init() {
	key := GetKey()
	client = redis.NewClient(&redis.Options{
		Addr:     key.Addr,
		Username: key.Username,
		Password: key.Password,
		DB:       key.DB,
	})
	log.Println("Redis Client Set!")
}
