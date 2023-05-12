package server

import (
	"global_ssh/db"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/creack/pty"
	"github.com/redis/go-redis/v9"
)

func Start() {
	db.HostMode = true
	client := db.Init()
	os.Setenv("TERM", "xterm-256color")
	shell := exec.Command("zsh")
	shell_pty, err := pty.Start(shell)
	if err != nil {
		log.Fatal("Failed to Start PTY due to:", err)
	}
	go reader(shell_pty, client)
	command(shell_pty)
}

func command(pty *os.File) {
	for {
		var input string = db.AwaitData(true)
		log.Println("Running Command: " + input)
		go pty.Write([]byte(input))
	}
}

func reader(pty *os.File, client *redis.Client) {
	for {
		buf := make([]byte, 1024)
		n, err := pty.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Println("Error While Reading: ", err)
			}
			continue
		}
		go func() {
			err = db.Send(string(buf[:n]), false, client)
			if err != nil {
				log.Println("Failed While Sending Data: ", err)
			}
		}()
	}
}
