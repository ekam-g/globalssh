package client

import (
	"bufio"
	"fmt"
	"global_ssh/redis"
	"log"
	"os"
)

func Run() {
	redis.HostMode = false
	redis.Init()
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
		in := bufio.NewReader(os.Stdin)
		input, err := in.ReadString('\n')
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
