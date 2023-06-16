package client

import (
	"bufio"
	"errors"
	"fmt"
	"globalssh/net"
	"log"
	"os"
	"strings"
	"time"
)

func checkEncryptionKey(Net net.Net) {
	consoleData := make(chan string)
	go func() {
		data, err := Net.Read(net.Result)
		if err != nil {
			if err.Error() == net.DecryptError {
				log.Fatal("Failed To Decrypt Please Check Password")
			} else {
				log.Fatalf("Failed To Read Due To %s\n", err)
			}
			return
		}
		consoleData <- data
	}()
	time.Sleep(time.Millisecond * 20)
	err := Net.Send("\x7f\x7f\x7f\x7f\x7f\x7f\x7f\x7f\x7f\x7f\x7f\x7f\n", net.Command)
	if err != nil {
		log.Fatalf("Failed To Send Data To Redis Exiting\n%s", err)
	}
	data, err := readWithTimeout(consoleData, time.Second*3)
	if err != nil {
		log.Fatalf("Incorrect/No Reply Received From Server,\nThis May Happen If You Decryption Key is Bad Or Your Wifi Is Bad,\nLet's Try Changing Your Decryption Key in %s", net.KeyLocation)
	}
	fmt.Print(data)
	return
}

func readWithTimeout(ch <-chan string, timeout time.Duration) (string, error) {
	select {
	case val := <-ch:
		return val, nil
	case <-time.After(timeout):
		return "", errors.New("timeout occurred")
	}
}

func _() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Key: ")
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Error reading input %s", err)
	}
	input = strings.ReplaceAll(input, "\n", "")
	input = strings.Trim(input, " ")
	return input
}
