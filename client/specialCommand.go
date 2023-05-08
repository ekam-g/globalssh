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
	split_input := strings.Split(input, " ")
	if !strings.Contains(split_input[0], "global_ssh") {
		return false
	}
	correct_command := false
	input = strings.ReplaceAll(input, "global_ssh", "")
	input = strings.ReplaceAll(input, "\n", "")
	input = strings.Trim(input, " ")
	correct_command = exit(input)
	if !correct_command {
		fmt.Println("Incorrect Command Given")
	}
	return true

}

func exit(input string) bool {
	if input == "client exit" {
		fmt.Println("Exiting Global SSH, Goodbye!")
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
		fmt.Println("To Exit Global_SSH Please Do {global_ssh client exit}")
		return 0
	}
	error := db.Send("\x03", true, client)
	if error != nil {
		log.Println("Failed To Send Redis Data due to: ", error)
	}
	return amount_singled + 1
}
