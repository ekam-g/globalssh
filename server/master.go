package server

import (
	"errors"
	"fmt"
	"global_ssh/client"
	"global_ssh/db"
	"global_ssh/termUtil"
	"io"
	"log"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/creack/pty"
	"github.com/redis/go-redis/v9"
	"golang.org/x/term"
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
	go userReader(shell_pty)
	command(shell_pty)
}

func command(pty *os.File) {
	setData := make(chan string)
	go writerWorker(setData, pty)
	for {
		var input string = db.AwaitData(db.Command)
		if termUtil.CheckGetSize(input, pty) {
			continue
		}
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
	go db.SenderWorker(worker, db.Result, client)
	for {
		buf := make([]byte, 1024)
		n, err := pty.Read(buf)
		if err != nil {
			if err == io.EOF {
				break // Break the loop when the pty is closed
			}
			if !errors.Is(err, syscall.EAGAIN) {
				log.Println("Error while reading:", err)
			}
			time.Sleep(10 * time.Millisecond) // Wait before the next read attempt
			continue
		}
		if n > 0 {
			worker <- string(buf[:n])
			fmt.Print(string(buf[:n]))
		}
	}
}

func userReader(pty *os.File) {
	if term.IsTerminal(0) {
		var special_command_data string
		oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
		if err != nil {
			log.Fatal(err)
		}
		defer term.Restore(int(os.Stdin.Fd()), oldState)
		worker := make(chan string)
		go writerWorker(worker, pty)
		for {
			b := make([]byte, 1)
			_, err = os.Stdin.Read(b)
			if err != nil {
				log.Println(err)
			}
			input := string(b[0])
			if input == "" {
				continue
			}
			special_command_data = client.StoreSpecialCommandData(special_command_data, input)
			if client.HandleSpecialCommands(special_command_data) {
				continue
			}
			worker <- input
		}
	}
}
