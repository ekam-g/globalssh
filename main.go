package main

import (
	"fmt"
	"global_ssh/client"
	"global_ssh/server"
	"os"
	"strings"
)

func main() {
	switch strings.Trim(os.Args[1], " ") {
	case "client":
		client.Run()
		break
	case "server":
		server.Start()
		break
	default:
		fmt.Println("Bad Arg Given, Please Put in server or client")
		os.Exit(1)
		break
	}
}
