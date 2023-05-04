package main

import (
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
	f.Write([]byte("sudo pacman -Syyu\n"))
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
			if _, err := os.Stdout.Write(buf[:n]); err != nil {
				panic(err)
			}
		}
	}()
	time.Sleep(time.Second * 2)
	f.Write([]byte("y"))
}
