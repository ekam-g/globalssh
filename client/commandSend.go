package client

import (
	"globalssh/net"
	"log"
	"os"
	"strconv"
	"time"
)

func CommandSend(hostName string, wait string, data string) {
	timeRead, err := strconv.Atoi(wait)
	if err != nil {
		log.Fatalf("Failed to Parse int due to %s\nExiting Program\n", err)
	}
	Net, _ := net.Init(false, hostName)
	go display(Net)
	time.Sleep(time.Millisecond * 20)
	go Net.Send(data, net.Command)
	time.Sleep(time.Second * time.Duration(timeRead))
	os.Exit(0)
}
