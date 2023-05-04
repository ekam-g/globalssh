package redis

import (
	"context"
	"errors"
	"strings"
)

func Send(val string) error {
	ctx := context.Background()

	return client.Set(ctx, HostName, GetWriteKey()+val, 0).Err()
}

func Read() (string, error) {
	ctx := context.Background()
	value, err := client.Get(ctx, HostName).Result()
	if err != nil {
		return "", err
	}
	if !strings.Contains(GetReadKey(), hostKey) {
		return value, errors.New("Not Allowed to read")
	}
	return value, nil

}
