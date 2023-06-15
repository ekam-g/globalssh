package net

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/creack/pty"
	"golang.org/x/term"
)

func Update() {
	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := term.Restore(fd, oldState)
		if err != nil {
			log.Fatal("Failed To Restore Terminal due to: ", err)
		}
	}()
	shellCode := [1]string{"bash <( curl -s https://raw.githubusercontent.com/carghai/globalssh/main/install.sh)\n"}
	shellPty, err := pty.Start(exec.Command("bash"))
	if err != nil {
		log.Fatalf("Failed to update due to a pty failure: %s\n", err)
	}
	go func() {
		for _, x := range shellCode {
			_, err = shellPty.Write([]byte(x))
			if err != nil {
				log.Fatalf("Failed to update due to unable to write to pty: %s\n", err)
			}
		}
	}()
	go func() {
		_, err := io.Copy(shellPty, os.Stdin)
		if err != nil {
			log.Fatalf("Failed to copy Stdin: %s\n", err)
		}
	}()
	fmt.Println("----Start Stdout----")
	var builder strings.Builder
	for {
		buf := make([]byte, 1024)
		n, err := shellPty.Read(buf)
		if err != nil {
			fmt.Printf("Error when reading Stdout: %s", err)
		}
		if n < 0 {
			continue
		}
		stdout := string(buf[:n])
		builder.WriteString(stdout)
		fmt.Print(stdout)
		if strings.Contains(builder.String(), "Welcome") {
			break
		}
	}
	fmt.Println("----End Stdout----")
}
