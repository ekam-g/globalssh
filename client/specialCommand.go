package client

import (
	"fmt"
	"globalssh/net"
	"io"
	"log"
	"os"
	"strings"

	"golang.org/x/term"
)

const Kill = "\x03"

const BackSpace = "\x7f"

func HandleSpecialCommands(input string, fd int, oldState *term.State, Net net.Net) bool {
	return exit(input, fd, oldState, Net)
}

func StoreSpecialCommandData(currentData string, input string) string {
	if len(currentData) > 100 {
		return ""
	}
	if input == BackSpace {
		if len(currentData) == 0 {
			return currentData
		}
		return currentData[:len(currentData)-1]
	}
	if input == " " || input == "\n" {
		return ""
	}
	return currentData + input
}

func exit(input string, fd int, oldState *term.State, Net net.Net) bool {
	if strings.Contains(input, "client-exit") {
		termClean(fd, oldState)
		Net.Close()
		log.SetOutput(io.Discard)
		fmt.Println("\nExiting Global SSH, Goodbye!")
		os.Exit(0)
	}
	return false
}

func termClean(fd int, oldState *term.State) {
	err := term.Restore(fd, oldState)
	if err != nil {
		log.Fatal("Failed To Restore Terminal due to: ", err)
	}
}
