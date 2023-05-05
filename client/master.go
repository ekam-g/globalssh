package client

import (
	"fmt"
	"global_ssh/redis"
	"log"
)

func Run() {
	go display()
	input()
}

func display() {
	for {
		data := redis.AwaitData(false)
		fmt.Print(data)
	}

}

func input() {
	for {
		var input string
		_, err := fmt.Scan(&input)
		if err != nil {
			log.Println(err)
			continue
		}
		err = redis.Send(input, true)
		if err != nil {
			log.Fatal(err)
		}
	}
}
