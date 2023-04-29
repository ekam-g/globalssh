package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os/exec"
	"time"
)

func main() {
	cmd := exec.Command("nvim", "main.go")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		defer stdin.Close()
		// After 3 seconds of running, send newline to cause program to exit.
		time.Sleep(time.Second * 3)
		io.WriteString(stdin, "\n")
	}()

	cmd.Start()

	// Scan and print command's stdout
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	// Wait for program to exit.
	cmd.Wait()
}
