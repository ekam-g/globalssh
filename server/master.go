package server

import (
	"global_ssh/db"
	"global_ssh/termUtil"
	"io"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"

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
	mutex := &sync.Mutex{}
	go reader(shell_pty, mutex, client)
	command(shell_pty, mutex)
}

func command(pty *os.File, mutex *sync.Mutex) {
	var waitgroup sync.WaitGroup
	for {
		var input string = db.AwaitData(true)
		if termUtil.CheckGetSize(input, pty) {
			continue
		}
		log.Println("Running Command: " + input)
		waitgroup.Wait()
		waitgroup.Add(1)
		go func() {
			log.Println("Locking in command")
			mutex.Lock()
			pty.Write([]byte(input))
			log.Println("Unlocking in command")
			mutex.Unlock()
			waitgroup.Done()
		}()
	}
}

func reader(pty *os.File, mutex *sync.Mutex, client *redis.Client) {
	var waitgroup sync.WaitGroup
	for {
		buf := make([]byte, 1024)
		tempLock(mutex)
		n, err := pty.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Println("Error While Reading: ", err)
			}
			continue
		}
		waitgroup.Wait()
		waitgroup.Add(1)
		go func() {
			err = db.Send(string(buf[:n]), false, client)
			if err != nil {
				log.Println("Failed While Sending Data: ", err)
			}
			waitgroup.Done()
		}()
	}
}

func tempLock(mutex *sync.Mutex) {
	mutex.Lock()
	time.Sleep(time.Millisecond * 5)
	mutex.Unlock()
}
