package main

import (
	"fmt"
	"global_ssh/client"
	"global_ssh/server"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Please give an arg like client or server")
		os.Exit(1)
	}
	switch strings.Trim(os.Args[1], " ") {
	case "client":
		client.Run()
	case "server":
		server.Start()
	default:
		fmt.Println("Bad Arg Given, Please Put in server or client")
		os.Exit(1)

	}
}
