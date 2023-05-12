package client

import (
	"fmt"
	"global_ssh/db"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/redis/go-redis/v9"
)

func handleSpecialCommands(input string) bool {
	return exit(input)
}

func storeSpecialCommandData(currentData string, input string) string {
	if input == " " || input == "\n" {
		return ""
	}
	return currentData + input
}

func exit(input string) bool {
	//take f8 singal to end code
	if strings.Contains(input, "client-exit") {
		fmt.Println("\nExiting Global SSH, Goodbye!")
		os.Exit(0)
	}
	return false
}

func signalHandler(client *redis.Client) {
	amount_singled := 0
	for {
		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		amount_singled = sigtermHandler(amount_singled, client)
	}

}

func sigtermHandler(amount_singled int, client *redis.Client) int {
	if amount_singled > 10 {
		fmt.Println("To Exit Global_SSH Please Do {client-exit}")
		return 0
	}
	error := db.Send("\x03", true, client)
	if error != nil {
		log.Println("Failed To Send Redis Data due to: ", error)
	}
	return amount_singled + 1
}
