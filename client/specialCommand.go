package client

import (
	"fmt"
	"os"
	"strings"
)

func handleSpecialCommands(input string) bool {
	split_input := strings.Split(input, " ")
	if !strings.Contains(split_input[0], "global_ssh") {
		return false
	}
	correct_command := false
	input = strings.ReplaceAll(input, "global_ssh", "")
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
