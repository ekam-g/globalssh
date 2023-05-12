package server

import (
	"global_ssh/db"
	"global_ssh/termUtil"
	"log"
	"os"
	"os/exec"

	"github.com/creack/pty"
	"github.com/redis/go-redis/v9"
)

func Start() {
	db.HostMode = true
	client, what_shell := db.Init()
	os.Setenv("TERM", "xterm-256color")
	shell := exec.Command(what_shell)
	shell_pty, err := pty.Start(shell)
	if err != nil {
		log.Fatal("Failed to Start PTY due to:", err)
	}
	go reader(shell_pty, client)
	command(shell_pty)
}

func command(pty *os.File) {
	setData := make(chan string)
	go writerWorker(setData, pty)
	for {
		var input string = db.AwaitData(true)
		if termUtil.CheckGetSize(input, pty) {
			continue
		}
		log.Print(input)
		setData <- input
	}
}

func writerWorker(setData chan string, pty *os.File) {
	for {
		input := <-setData
		_, err := pty.Write([]byte(input))
		if err != nil {
			log.Println(err)
		}
	}
}

func reader(pty *os.File, client *redis.Client) {
	worker := make(chan string)
	go db.SenderWorker(worker, true, client)
	for {
		buf := make([]byte, 1024)
		n, err := pty.Read(buf)
		if err != nil {
			log.Println("Error While Reading: ", err)
			continue
		}
		worker <- string(buf[:n])
	}
}
