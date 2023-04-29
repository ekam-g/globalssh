package main

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
)

func main() {
	cmd := exec.Command("echo", "hello")
	cmdReader, _ := cmd.StdoutPipe()
	scanner := bufio.NewScanner(cmdReader)
	done := make(chan bool)
	go func() {
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
		done <- true
	}()
	err := cmd.Start()
	if err != nil {
		log.Fatalln(err)
	}
	<-done
	err = cmd.Wait()
	if err != nil {
		log.Fatalln(err)
	}
}
