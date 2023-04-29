package main

import (
	"bufio"
	"log"
	"os/exec"
)

func main() {
	cmd := exec.Command("bash", "-c", "nvim", "main.go")

	outPipe, _ := cmd.StdoutPipe()
	stdErrorPipe, _ := cmd.StderrPipe()
	scanner := bufio.NewScanner(outPipe)
	stdErrorScanner := bufio.NewScanner(stdErrorPipe)
	go func() {
		for scanner.Scan() {
			log.Println(stdErrorScanner.Text())
			log.Println(scanner.Text())
		}
	}()
	cmd.Start()
	cmd.Wait()
	log.Println("Done")
}
