package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/creack/pty"
)

func main() {
	c := exec.Command("zsh")
	f, err := pty.Start(c)
	if err != nil {
		panic(err)
	}
	f.Write([]byte("neofetch\n"))
	go reader(f)
	command(f)
	// f.Write([]byte("ls\n"))
	time.Sleep(time.Second * 5)
}

func command(f *os.File) {
	for {
		var input string
		fmt.Scan(&input)
		input += "\n"
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
		fmt.Print(string(buf[:n]))
	}
}
