package db

import (
	"log"
	"reflect"
	"sync"

	"github.com/redis/go-redis/v9"
)

func SenderWorker(data chan string, HostMode bool, client *redis.Client) {
	var waiting_data_mx sync.Mutex
	run := make(chan struct{})
	waiting_data := ""
	go waitingWorker(run, &waiting_data, &waiting_data_mx, HostMode, client)
	for {
		send_data := <-data
		waiting_data_mx.Lock()
		waiting_data += send_data
		waiting_data_mx.Unlock()
		run <- struct{}{}
	}

}

func waitingWorker(wait chan struct{}, waiting_data *string, waiting_data_mx *sync.Mutex, HostMode bool, client *redis.Client) {
	var allowSend sync.Mutex
	redis_send := make(chan string)
	go send(redis_send, client, &allowSend, HostMode)
	for {
		<-wait
		waiting_data_mx.Lock()
		if *waiting_data == "" {
			continue
		}
		if !mutexLocked(&allowSend) {
			redis_send <- *waiting_data
			*waiting_data = ""
		}
		waiting_data_mx.Unlock()
	}
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
