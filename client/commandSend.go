package client

import (
	"globalssh/net"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func CommandSend(hostName string, wait string, data string) {
	timeRead, err := strconv.Atoi(wait)
	data = strings.ReplaceAll(data, "\\n", "\n")
	if err != nil {
		log.Fatalf("Failed to Parse int due to %s\nExiting Program\n", err)
	}
	Net, _ := net.Init(false, hostName)
	if timeRead == 0 {
		Net.Send(data, net.Command)
		os.Exit(0)
	}
	go display(Net)
	time.Sleep(time.Millisecond * 10)
	go Net.Send(data, net.Command)
	time.Sleep(time.Second * time.Duration(timeRead))
	os.Exit(0)
}
