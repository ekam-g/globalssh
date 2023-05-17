package db

import (
	"log"

	"github.com/redis/go-redis/v9"
)

var (
	HostMode       bool = true
	HostName       string
	command_stream *redis.PubSub
	result_stream  *redis.PubSub
)

func Init() (*redis.Client, string) {
	key := GetKey()
	client := redis.NewClient(&redis.Options{
		Addr:     key.Addr,
		Username: key.Username,
		Password: key.Password,
		DB:       key.DB,
	})
	GetConnection(client)
	EncryptionKey = []byte(key.Key)
	var err error
	if HostMode {
		err = Send("Server Is On", false, client)
	} else {
		err = Send("neofetch\n", true, client)
	}
	if err != nil {
		log.Println("Failed To Make Redis Connection, Please Review Your Config And Wifi.\nAdvanced Error Details:", err)
	} else {
		log.Println("Redis Connection Verfied And Working, Starting Global SSH!")
	}
	return client, key.Shell
}
