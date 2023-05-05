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
		panic(err)
	}
	shell_pty.Write([]byte("neofetch\n"))
	go reader(shell_pty)
	command(shell_pty)
}

func command(f *os.File) {
	for {
		var input string = redis.AwaitData(true)
		input += "\n"
		log.Println("Running Command: " + input)
		f.Write([]byte(input))
	}
}

func reader(f *os.File) {
	for {
		buf := make([]byte, 1024)
		n, err := f.Read(buf)
		if err != nil {
			if err != io.EOF {
				panic(err)
			}
			break
		}
		err = redis.Send(string(buf[:n]), false)
		if err != nil {
			log.Println(err)
		}
		// fmt.Print(string(buf[:n]))
	}
}
