package server

import (
	"global_ssh/redis"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/creack/pty"
)

func Start() {
	redis.HostMode = true
	redis.Init()
	shell := exec.Command("zsh")
	shell_pty, err := pty.Start(shell)
	if err != nil {
		log.Fatal("Failed to Start PTY due to:", err)
	}
	go reader(shell_pty)
	command(shell_pty)
}

func command(pty *os.File) {
	for {
		var input string = redis.AwaitData(true)
		log.Println("Running Command: " + input)
		pty.Write([]byte(input))
	}
}

func reader(pty *os.File) {
	for {
		buf := make([]byte, 1024)
		n, err := pty.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Println("Error While Reading: ", err)
			}
			continue
		}
		err = redis.Send(string(buf[:n]), false)
		if err != nil {
			log.Println("Failed While Sending Data: ", err)
		}
	}
}
