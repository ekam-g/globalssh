package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func GetConnection() {
	ctx := context.Background()
	command_stream = client.Subscribe(ctx, HostName+"command")
	result_stream = client.Subscribe(ctx, HostName+"result")
}

func Send(val string, command_send bool) error {
	ctx := context.Background()
	return client.Publish(ctx, HostName+Extention(command_send), val).Err()
}

func Extention(command_send bool) string {
	if command_send {
		return "command"
	} else {
		return "result"
	}
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
	return data.Payload, nil
}

func AwaitData(command_version bool) string {
	for {
		data, err := Read(command_version)
		if err == nil {
			return data
		}
		fmt.Println(err)
		// time.Sleep(time.Millisecond * 10)
	}
}
