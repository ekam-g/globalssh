package net

import (
	speedJson "github.com/json-iterator/go"
	"log"
	"os"
	"strings"
	"time"

	"github.com/creack/pty"
	"golang.org/x/term"
)

const termCommand string = "&%#$&^!@%#$^KJH#G$@#$"

type TermSize struct {
	Width  uint16
	Length uint16
}

func (net Net) SetSize() {
	oldSize := TermSize{}
	for {
		time.Sleep(time.Second * 5)
		width, length, err := term.GetSize(int(os.Stdin.Fd()))
		if err != nil {
			log.Println("Failed to Get Size of Terminal")
			return
		}
		termSize := TermSize{
			Width:  uint16(width),
			Length: uint16(length),
		}
		if termSize == oldSize {
			time.Sleep(time.Second * 10)
		}
		oldSize = termSize
		sendData, err := speedJson.ConfigCompatibleWithStandardLibrary.Marshal(termSize)
		if err != nil {
			log.Fatal("FATAL INTERNAL ERROR\nUNABLE TO SET JSON:", err)
		}
		err = net.Send(termCommand+string(sendData), true)
		if err != nil {
			log.Println("Failed To Send Redis Data due to: ", err)
		}
	}
}

func CheckGetSize(input string, pty_term *os.File) bool {
	if !strings.Contains(input, termCommand) {
		return false
	}
	input = strings.ReplaceAll(input, termCommand, "")
	termSize := TermSize{}
	err := speedJson.ConfigCompatibleWithStandardLibrary.Unmarshal([]byte(input), &termSize)
	if err != nil {
		log.Println(err)
		return false
	}
	window := pty.Winsize{
		Rows: termSize.Length,
		Cols: termSize.Width,
	}
	err = pty.Setsize(pty_term, &window)
	if err != nil {
		log.Println("Failed to Resize Window Due to: ", err)
	}
	return true
}
