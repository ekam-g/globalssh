package client

import (
	"fmt"
	"global_ssh/db"
	"global_ssh/termUtil"
	"log"
	"os"
	"sync"

	"github.com/redis/go-redis/v9"
	"golang.org/x/term"
)

func Run() {
	db.HostMode = false
	client := db.Init()
	go termUtil.SetSize(client)
	go signalHandler(client)
	go display()
	input(client)
}

func display() {
	for {
		data := db.AwaitData(false)
		fmt.Print(data)
	}
}

func input(client *redis.Client) {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatal(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)
	var special_command_data string
	var waitgroup sync.WaitGroup
	var mutex sync.Mutex
	for {
		b := make([]byte, 1)
		_, err = os.Stdin.Read(b)
		if err != nil {
			log.Println(err)
		}
		input := string(b[0])
		if input == "" {
			continue
		}
		special_command_data = storeSpecialCommandData(special_command_data, input)
		if handleSpecialCommands(special_command_data) {
			continue
		}
		go func() {
			mutex.Lock()
			waitgroup.Wait()
			waitgroup.Add(1)
			mutex.Unlock()
			err = db.Send(input, true, client)
			if err != nil {
				log.Println("Failed to send, due to: ", err)
			}
			mutex.Lock()
			waitgroup.Done()
			mutex.Unlock()
		}()
	}
}
