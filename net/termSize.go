package net

import (
	"globalssh/client"
	"log"
	"os"
	"strings"
	"time"

	"github.com/creack/pty"
	speedJson "github.com/json-iterator/go"
	"github.com/mattn/go-isatty"
	"golang.org/x/term"
)

const resendTime uint8 = 30

const termCommand string = "&%#$&^!@%#$^KJH#G$@#$"

const exitCommand string = "$KU%JGH#@K$JH%$UYGT%O@&#$T%@J#H$%GOUYFRO*@$%JHLOV@#KHB$%CHKG$#F%JKL@H#$B%JKGC@$#L%IUJHG@#$&(*%@JLNB$V%GC"

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

// CheckGetSize returns true if command is run
func CheckGetSize(input string, ptyTerm *os.File) bool {
	//exit on successful command
	if onExitLocalSize(ptyTerm, input) {
		return true
	}
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

// onExitLocalSize returns true if command is run
func onExitLocalSize(ptyTerm *os.File, input string) bool {
	if strings.Contains(input, exitCommand) {
		if isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd()) {
			go func() {
				SetLocalSize(ptyTerm)
				_, err := ptyTerm.Write([]byte(client.SetDisplay))
				if err != nil {
					log.Println("Failed To Write Data Due To, ", err)
				}
			}()
		}
		log.Println("\nClient Has Exited")
		return true
	}
	return false
}

func (net Net) ExitSizeSet() {
	err := net.Send(exitCommand, Command)
	if err != nil {
		log.Println("Failed To Send Due To: ", err)
	}
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
