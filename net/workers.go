package net

import (
	"log"
	"strings"
)

const ImportantWorkerLimit = 4000

const LimitedWorkerLimit = 100

func (net Net) SenderWorker(data chan string, HostMode bool) {
	for {
		sendData := BulkData(data)
		err := net.Send(sendData, HostMode)
		if err != nil {
			log.Println("Failed to send, due to: ", err)
		}
	}
}

func BulkData(ch chan string) string {
	var builder strings.Builder
	for {
		select {
		case val := <-ch:
			_, err := builder.WriteString(val)
			if err != nil {
				return builder.String()
			}
		default:
			if builder.Len() == 0 {
				return <-ch
			}
			return builder.String()
		}
	}
}
