package main

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
)

func main() {
	cmd := exec.Command("nvim", "main.go")
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		log.Panic(err)
	}
	scanner := bufio.NewScanner(cmdReader)
	done := make(chan bool)
	go func() {
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
		done <- true
	}()
	err = cmd.Start()
	if err != nil {
		log.Panic(err)
	}
	<-done
	err = cmd.Wait()
	if err != nil {
		log.Panic(err)
	}
}
