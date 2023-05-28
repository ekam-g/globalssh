package server

import (
	"testing"
	"time"
)

func TestApp(t *testing.T) {
	// This is used for profiling the app
	go Start()
	time.Sleep(time.Second * 30)
}
