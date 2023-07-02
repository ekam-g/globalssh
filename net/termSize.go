package net

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/creack/pty"
	speedJson "github.com/json-iterator/go"
	"golang.org/x/term"
)

const resendTime uint8 = 30

const termCommand string = "&%#$&^!@%#$^KJH#G$@#$"

type TermSize struct {
	Width  uint16
	Length uint16
}

func (net Net) SetSize() {
	oldSize := TermSize{}
	var amountTried uint8 = 0
	time.Sleep(time.Second * 3)
	for {
		time.Sleep(time.Second * 1)
		width, length, err := term.GetSize(int(os.Stdin.Fd()))
		if err != nil {
			log.Println("Failed to Get Size of Terminal due to: ", err)
			return
		}
		termSize := TermSize{
			Width:  uint16(width),
			Length: uint16(length),
		}
		if termSize == oldSize {
			if amountTried >= resendTime {
				amountTried = 0
			} else {
				amountTried++
				continue
			}
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

func CheckGetSize(input string, ptyTerm *os.File) bool {
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
	err = pty.Setsize(ptyTerm, &window)
	if err != nil {
		log.Println("Failed to Resize Window Due to: ", err)
	}
	return true
}

func SetLocalSize(ptyTerm *os.File) {
	width, length, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		log.Println("Failed to Get Size of Terminal due to: ", err)
		return
	}
	termSize := TermSize{
		Width:  uint16(width),
		Length: uint16(length),
	}
	window := pty.Winsize{
		Rows: termSize.Length,
		Cols: termSize.Width,
	}
	err = pty.Setsize(ptyTerm, &window)
	if err != nil {
		log.Println("Failed to Resize Window Due to: ", err)
	}

}
