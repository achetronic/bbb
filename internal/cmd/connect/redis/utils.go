package redis

import (
	"os"
	"os/signal"
	"syscall"
)

// WaitSignalAfter waits for a signal to close the process gracefully
// after executing a callback function
func WaitSignal() {

	// Capture some OS signals to close the process gracefully before closing
	SignalsToReceive := []os.Signal{
		syscall.SIGTERM, syscall.SIGINT, // Ctrl+C
		syscall.SIGTSTP, // Ctrl+Z
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, SignalsToReceive...)

	<-c
}
