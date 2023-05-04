package main

import (
	"fmt"
	"io"
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
	go func() {
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
	}()
	// f.Write([]byte("ls\n"))
	time.Sleep(time.Second * 5)
}
