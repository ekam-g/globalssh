package client

import (
	"fmt"
	"global_ssh/db"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
	"golang.org/x/term"
)

func Run() {
	db.HostMode = false
	client := db.Init()
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
		// if handleSpecialCommands(input) {
		// 	continue
		// }
		go func() {
			err = db.Send(input, true, client)
			if err != nil {
				log.Println("Failed to send, due to: ", err)
			}
		}()
	}
}
