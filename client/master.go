package client

import (
	"fmt"
	"log"
	"os"

	"globalssh/net"

	"golang.org/x/term"
)

func Run(host string) {
	Net, serverShell := net.Init(false, host)
	if host == "" {
		log.Printf("The Server Computer is using: %s\n", serverShell)
	}
	log.Printf("Connecting to %s\n", Net.HostName)
	go Net.SetSize()
	go signalHandler(Net)
	go display(Net)
	worker := make(chan string, net.ImportantWorkerLimit)
	go Net.SenderWorker(worker, net.Command)
	Input(Net, worker)
}

func display(Net net.Net) {
	display := make(chan string, net.LimitedWorkerLimit)
	go displayWorker(display)
	for {
		data := Net.AwaitData(net.Result)
		display <- data
	}
}

func displayWorker(data chan string) {
	for {
		display := net.BulkData(data)
		fmt.Print(display)
	}
}

func Input(Net net.Net, worker chan string) {
	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		log.Fatal(err)
	}
	var specialCommandData string
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
		specialCommandData = StoreSpecialCommandData(specialCommandData, input)
		if HandleSpecialCommands(specialCommandData, fd, oldState) {
			continue
		}
		worker <- input

	}
}
