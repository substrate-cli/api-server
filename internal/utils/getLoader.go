package utils

import (
	"fmt"
	"time"
)

var (
	showLoader = false
	stopChan   chan bool
)

func StartLoader(message string) {
	if showLoader {
		return // already running
	}
	showLoader = true
	stopChan = make(chan bool)

	go func() {
		chars := `|/-\`
		i := 0
		for {
			select {
			case <-stopChan:
				fmt.Print("\r\033[K") // clear line
				return
			default:
				fmt.Printf("\r\033[K%s %c", message, chars[i%len(chars)])
				i++
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()
}

// StopLoader stops the loader
func StopLoader() {
	if showLoader {
		stopChan <- true
		close(stopChan)
		showLoader = false
	}
}
