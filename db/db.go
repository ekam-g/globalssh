package db

import (
	"log"

	"github.com/redis/go-redis/v9"
)

const (
	hostKey   string = "&&**$$&*#$"
	clientKey string = "@@$$##*()$"
)

var (
	HostMode       bool = true
	HostName       string
	command_stream *redis.PubSub
	result_stream  *redis.PubSub
)

func GetWriteKey() string {
	if HostMode {
		return clientKey
	} else {
		return hostKey
	}
}

func GetReadKey() string {
	if !HostMode {
		return clientKey
	} else {
		return hostKey
	}
}

func Init() *redis.Client {
	key := GetKey()
	client := redis.NewClient(&redis.Options{
		Addr:     key.Addr,
		Username: key.Username,
		Password: key.Password,
		DB:       key.DB,
	})
	GetConnection(client)
	var err error
	if HostMode {
		err = Send("Server Is On", false, client)
	} else {
		err = Send("neofetch\n", true, client)
	}
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Redis Client Set!")
	return client
}
