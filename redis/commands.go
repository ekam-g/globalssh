package redis

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	old_key_error = "Old Data, Skipping"
)

var (
	old_write_key = 0
	old_read_key  = 0
)

func GetConnection() {
	ctx := context.Background()
	command_stream = client.Subscribe(ctx, HostName+"command")
	result_stream = client.Subscribe(ctx, HostName+"result")
}

func Send(val string, command_send bool) error {
	ctx := context.Background()
	// client.Set(ctx, "hello", "works",).Err()
	return client.Publish(ctx, HostName+Extention(command_send), NewWriteKey()+val).Err()
}

func Extention(command_send bool) string {
	if command_send {
		return "command"
	} else {
		return "result"
	}
}

func CheckReadKey(data string) (string, bool) {
	keys := strings.Split(data, "|")
	if len(keys) != 2 {
		log.Println("Bad Data Recived, '|' missing")
		return data, false
	}
	key, err := strconv.Atoi(keys[0])
	if err != nil {
		log.Println("Failed to Parse Key Int: ", err, key)
		return keys[1], false
	}
	if key == old_read_key {
		return keys[1], false
	}
	old_read_key = key
	return keys[1], true
}

func NewWriteKey() string {
	random := rand.Int()
	if random == old_write_key {
		NewWriteKey()
	}
	old_write_key = random
	return fmt.Sprint(random) + "|"
}

func Read(commmand_version bool) (string, error) {
	ctx := context.Background()
	var stream *redis.PubSub
	if commmand_version {
		stream = command_stream
	} else {
		stream = result_stream
	}
	data, err := stream.ReceiveMessage(ctx)
	if err != nil {
		return "", err
	}
	clean_data, allow := CheckReadKey(data.Payload)
	if !allow {
		return clean_data, errors.New("Not Allowed to read")
	}
	return clean_data, nil
}

func AwaitData(command_version bool) string {
	for {
		data, err := Read(command_version)
		if err == nil {
			return data
		}
		fmt.Println(err)
		time.Sleep(time.Millisecond * 10)
	}
}
