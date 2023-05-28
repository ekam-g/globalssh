package client

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"globalssh/net"

	"golang.org/x/term"
)

func HandleSpecialCommands(input string, fd int, oldState *term.State) bool {
	return exit(input, fd, oldState)
}

func StoreSpecialCommandData(currentData string, input string) string {
	if input == " " || input == "\n" {
		return ""
	}
	return currentData + input
}

func exit(input string, fd int, oldState *term.State) bool {
	if strings.Contains(input, "client-exit") {
		fmt.Println("\nExiting Global SSH, Goodbye!")
		os.Exit(termClean(fd, oldState))
	}
	return false
}

func signalHandler(Net net.Net) {
	amountSingled := 0
	for {
		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		amountSingled = sigtermHandler(amountSingled, Net)
	}

}

func sigtermHandler(amountSingled int, Net net.Net) int {
	if amountSingled > 10 {
		fmt.Println("To Exit Global_SSH Please Do {client-exit}")
		return 0
	}
	err := Net.Send("\x03", true)
	if err != nil {
		log.Println("Failed To Send Redis Data due to: ", err)
	}
	return amountSingled + 1
}

func termClean(fd int, oldState *term.State) int {
	err := term.Restore(fd, oldState)
	if err != nil {
		log.Fatal("Failed To Restore Terminal due to: ", err)
	}
	return 0
}
