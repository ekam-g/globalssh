package client

import (
	"fmt"
	"os"
	"strings"
)

func handleSpecialCommands(input string) bool {
	if strings.Contains(input, "global_ssh") {
		correct_command := false
		input = strings.ReplaceAll(input, "global_ssh", "")
		input = strings.Trim(input, " ")
		correct_command = exit(input)
		if !correct_command {
			fmt.Println("Incorrect Command Given")
		}
		return true
	}
	return false
}

func exit(input string) bool {
	if input == "client exit" {
		fmt.Println("Exiting Global SSH, Goodbye!")
		os.Exit(0)
	}
	return false
}
