package server

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"syscall"
	"time"

	"globalssh/client"
	"globalssh/net"

	"github.com/mattn/go-isatty"

	"github.com/creack/pty"
)

func Start() {
	Net, whatShell := net.Init(true, "")
	err := os.Setenv("TERM", "xterm-256color")
	if err != nil {
		log.Println("Failed To set Env var for TERM, commands like htop may not work properly")
		return
	}
	shell := exec.Command(whatShell)
	shellPty, err := pty.Start(shell)
	if err != nil {
		log.Fatal("Failed to Start PTY due to:", err)
	}
	log.Printf("Server Name is: '%s'\n", Net.HostName)
	log.Printf("Using %s as shell\n", shell)
	tty := false
	if isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd()) {
		log.Println("Starting Share Mode(Found TTY terminal)")
		tty = true
		go net.SetLocalSize(shellPty)
		worker := make(chan string, net.LimitedWorkerLimit)
		go writerWorker(worker, shellPty)
		go client.Input(Net, worker)
	}
	go reader(shellPty, Net, tty)
	command(shellPty, Net)
}

func command(pty *os.File, Net net.Net) {
	setData := make(chan string, net.LimitedWorkerLimit)
	go writerWorker(setData, pty)
	for {
		input := Net.AwaitData(net.Command)
		if net.CheckGetSize(input, pty) {
			continue
		}
		setData <- input
	}
}

func writerWorker(setData chan string, pty *os.File) {
	for {
		input := net.BulkData(setData)
		_, err := pty.Write([]byte(input))

		if err != nil {
			log.Println(err)
		}
	}
}

func reader(pty *os.File, Net net.Net, tty bool) {
	worker := make(chan string, net.ImportantWorkerLimit)
	go Net.SenderWorker(worker, net.Result)
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
		if n < 0 {
			continue
		}
		worker <- string(buf[:n])
		if !tty {
			continue
		}
		fmt.Print(string(buf[:n]))
	}
}
