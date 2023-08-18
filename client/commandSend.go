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
	data = strings.ReplaceAll(data, "\\n", "\n")
	Net, _ := net.Init(false, hostName)
	timeRead, err := strconv.Atoi(wait)
	if err != nil {
		// if wait == "n" {
		// 	err := Net.Send(data, net.Command)
		// 	checkErr(err)
		// 	result, err := Net.Read(net.Result)
		// 	checkErr(err)

		// }
		log.Fatalf("Failed to Parse int due to %s\nExiting Program\n", err)
	}
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

func checkErr(err error) {
	if err != nil {
		log.Fatalf("Exiting due to %s", err)
	}
}
