package server

import (
	"global_ssh/db"
	"global_ssh/termUtil"
	"log"
	"os"
	"os/exec"
	"sync"

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
	var waitgroup sync.WaitGroup
	var mutex sync.Mutex
	for {
		var input string = db.AwaitData(true)
		if termUtil.CheckGetSize(input, pty) {
			continue
		}
		log.Print(input)
		go func() {
			mutex.Lock()
			waitgroup.Wait()
			waitgroup.Add(1)
			mutex.Unlock()
			pty.Write([]byte(input))
			mutex.Lock()
			waitgroup.Done()
			mutex.Unlock()
		}()
	}
}

func reader(pty *os.File, client *redis.Client) {
	var waitgroup sync.WaitGroup
	var mutex sync.Mutex
	for {
		buf := make([]byte, 1024)
		n, err := pty.Read(buf)
		if err != nil {
			log.Println("Error While Reading: ", err)
			continue
		}
		go func() {
			mutex.Lock()
			waitgroup.Wait()
			waitgroup.Add(1)
			mutex.Unlock()
			err = db.Send(string(buf[:n]), false, client)
			if err != nil {
				log.Println("Failed While Sending Data: ", err)
			}
			mutex.Lock()
			waitgroup.Done()
			mutex.Unlock()
		}()
	}
}
