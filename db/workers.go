package db

import (
	"log"
	"reflect"
	"sync"

	"github.com/redis/go-redis/v9"
)

func SenderWorker(data chan string, HostMode bool, client *redis.Client) {
	var allowSend sync.Mutex
	redis_send := make(chan string)
	go send(redis_send, client, &allowSend)
	waiting_data := ""
	for {
		send_data := <-data
		waiting_data += send_data
		if !mutexLocked(&allowSend) {
			redis_send <- waiting_data
			waiting_data = ""
		}
	}
}

func send(data chan string, client *redis.Client, allowSend *sync.Mutex) {
	for {
		send_data := <-data
		allowSend.Lock()
		err := Send(send_data, HostMode, client)
		if err != nil {
			log.Println("Failed to send, due to: ", err)
		}
		allowSend.Unlock()
	}
}

const mutex_locked = 1

func mutexLocked(m *sync.Mutex) bool {
	state := reflect.ValueOf(m).Elem().FieldByName("state")
	return state.Int()&mutex_locked == mutex_locked
}
