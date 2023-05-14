package client

import (
	"fmt"
	"global_ssh/db"
	"global_ssh/termUtil"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
	"golang.org/x/term"
)

func Run() {
	db.HostMode = false
	client, _ := db.Init()
	go termUtil.SetSize(client)
	go signalHandler(client)
	go display()
	input(client)
}

func display() {
	displayWorker := make(chan string)
	go diplayWorker(displayWorker)
	for {
		data := db.AwaitData(db.Result)
		displayWorker <- data
	}
}

func diplayWorker(data chan string) {
	for {
		display := <-data
		fmt.Print(display)
	}
}

func input(client *redis.Client) {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatal(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)
	var special_command_data string
	worker := make(chan string)
	go db.SenderWorker(worker, db.Command, client)
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
		worker <- input

	}
}
