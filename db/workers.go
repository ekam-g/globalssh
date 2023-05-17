package db

import (
	"log"

	"github.com/redis/go-redis/v9"
)

func SenderWorker(data chan string, HostMode bool, client *redis.Client) {
	for {
		send_data := BulkData(data)
		err := Send(send_data, HostMode, client)
		if err != nil {
			log.Println("Failed to send, due to: ", err)
		}
	}
}

func BulkData(ch chan string) string {
	buffer := ""
	for {
		select {
		case val, _ := <-ch:
			buffer += val
		default:
			if buffer == "" {
				return <-ch
			}
			return buffer
		}

	}
}
