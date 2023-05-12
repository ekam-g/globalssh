package termUtil

import (
	"encoding/json"
	"global_ssh/db"
	"log"
	"os"
	"strings"
	"time"

	"github.com/creack/pty"
	"github.com/redis/go-redis/v9"
	"golang.org/x/term"
)

const termCommand string = "&%#$&^!@%#$^KJH#G$@#$"

type TermSize struct {
	Width  uint16
	Length uint16
}

func SetSize(client *redis.Client) {
	old_size := TermSize{}
	for {
		time.Sleep(time.Second * 5)
		width, length, err := term.GetSize(int(os.Stdin.Fd()))
		if err != nil {
			log.Println("Failed to Get Size of Terminal")
			return
		}
		term_size := TermSize{
			Width:  uint16(width),
			Length: uint16(length),
		}
		if term_size == old_size {
			time.Sleep(time.Second * 10)
		}
		old_size = term_size
		send_data, err := json.Marshal(term_size)
		if err != nil {
			log.Fatal("FATAL INTERNAL ERROR\nUNABLE TO SET JSON:", err)
		}
		err = db.Send(termCommand+string(send_data), true, client)
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
	err := json.Unmarshal([]byte(input), &termSize)
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
