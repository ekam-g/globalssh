package db

import (
	"log"
	"reflect"
	"sync"
	"sync/atomic"

	"github.com/redis/go-redis/v9"
)

func SenderWorker(data chan string, HostMode bool, client *redis.Client) {
	var waiting_data_mx sync.Mutex
	waiting_data := ""
	var allowSend sync.Mutex
	redis_send := make(chan string)
	go send(redis_send, client, &allowSend, HostMode)
	var thread_on atomic.Bool
	thread_on.Store(false)
	for {
		send_data := <-data
		if mutexLocked(&allowSend) {
			waiting_data_mx.Lock()
			waiting_data += send_data
			waiting_data_mx.Unlock()
			if !thread_on.Load() {
				go waitingWorker(&waiting_data, &waiting_data_mx, &allowSend, &redis_send, &thread_on)
			}
		} else {
			redis_send <- send_data
		}
	}
}

func waitingWorker(waiting_data *string, waiting_data_mx *sync.Mutex, allowSend *sync.Mutex, redis_send *chan string, thread_on *atomic.Bool) {
	thread_on.Store(true)
	allowSend.Lock()
	allowSend.Unlock()
	waiting_data_mx.Lock()
	*redis_send <- *waiting_data
	waiting_data_mx.Unlock()
	thread_on.Store(false)
}

func send(data chan string, client *redis.Client, allowSend *sync.Mutex, HostMode bool) {
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
