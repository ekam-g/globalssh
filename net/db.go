package net

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"io"
	"log"

	"github.com/redis/go-redis/v9"
)

type Net struct {
	HostMode      bool
	HostName      string
	CommandStream *redis.PubSub
	ResultStream  *redis.PubSub
	Client        *redis.Client
	EncryptionKey cipher.Block
}

func (net Net) Close() {
	//Kill command
	log.SetOutput(io.Discard)
	fmt.Println()
	_ = net.ResultStream.Close()
	_ = net.CommandStream.Close()
	_ = net.Client.Close()
}

func Init(HostMode bool, name string) (Net, string) {
	key := GetKey()
	if name != "" {
		key.HostName = name
	}
	client := redis.NewClient(&redis.Options{
		Addr:     key.Addr,
		Username: key.Username,
		Password: key.Password,
		DB:       int(key.DB),
	})
	commandStream, resultStream := GetConnection(client, key.HostName)
	var EncryptionKey cipher.Block
	if key.Key != "" {
		paddedKey := padKey([]byte(key.Key))
		key, err := aes.NewCipher(paddedKey)
		if err != nil {
			log.Fatal("Failed To Create Encryption Key, Please fix it: ", err)
		}
		EncryptionKey = key
	} else {
		EncryptionKey = nil
	}
	net := Net{
		CommandStream: commandStream,
		ResultStream:  resultStream,
		Client:        client,
		HostMode:      HostMode,
		EncryptionKey: EncryptionKey,
		HostName:      key.HostName,
	}
	var err error
	if HostMode {
		err = net.Send("Server Is On", false)
	} else {
		err = net.Send("\x7f\x7f\x7f\x7f\x7f\x7f\x7f\x7f\x7f\x7f\x7f\x7fneofetch\n", true)
	}
	if err != nil {
		log.Println("Failed To Make Redis Connection, Please Review Your Config And Wifi.\nAdvanced Error Details:", err)
	} else {
		log.Println("Redis Connection Verified And Working, Starting Global SSH!")
	}
	return net, key.Shell
}
