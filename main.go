package main

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
)

func main() {
	cmd := exec.Command("bash", "-c", "neofetch")

	outPipe, _ := cmd.StdoutPipe()
	stdErrorPipe, _ := cmd.StderrPipe()
	scanner := bufio.NewScanner(outPipe)
	stdErrorScanner := bufio.NewScanner(stdErrorPipe)
	go func() {
		for scanner.Scan() {
			if scanner.Text() != "" {
				fmt.Println(scanner.Text())
			} else {
				log.Println(stdErrorScanner.Text())

			}
		}
	}()
	cmd.Start()
	cmd.Wait()
	log.Println("Done")
}
