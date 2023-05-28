package net

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

const (
	Command = true
	Result  = false
)

func GetConnection(client *redis.Client, HostName string) (*redis.PubSub, *redis.PubSub) {
	ctx := context.Background()
	commandStream := client.Subscribe(ctx, HostName+"command")
	resultStream := client.Subscribe(ctx, HostName+"result")
	return commandStream, resultStream
}

func (net Net) Send(val string, commandSend bool) error {
	ctx := context.Background()
	data, err := net.encrypt(val)
	if err != nil {
		return err
	}
	return net.Client.Publish(ctx, net.HostName+Extension(commandSend), data).Err()
}

func Extension(commandSend bool) string {
	if commandSend {
		return "command"
	} else {
		return "result"
	}
}

func (net Net) Read(commandVersion bool) (string, error) {
	ctx := context.Background()
	var stream *redis.PubSub
	if commandVersion {
		stream = net.CommandStream
	} else {
		stream = net.ResultStream
	}
	data, err := stream.ReceiveMessage(ctx)
	if err != nil {
		return "", err
	}
	decrypted, err := net.decrypt(data.Payload)
	if err != nil {
		return "", err
	}
	return decrypted, nil
}

func (net Net) AwaitData(commandVersion bool) string {
	for {
		data, err := net.Read(commandVersion)
		if err == nil {
			return data
		}
		fmt.Println(err)
	}
}
