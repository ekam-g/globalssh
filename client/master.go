package client

import (
	"bufio"
	"fmt"
	"global_ssh/db"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

func Run() {
	db.HostMode = false
	client := db.Init()
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
	for {
		in := bufio.NewReader(os.Stdin)
		input, err := in.ReadString('\n')
		if err != nil {
			log.Println(err)
			continue
		}
		if handleSpecialCommands(input) {
			continue
		}
		err = db.Send(input, true, client)
		if err != nil {
			log.Fatal(err)
		}
	}
}
