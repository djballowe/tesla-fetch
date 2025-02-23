package ui

import (
	"fmt"
	"strings"
	"time"
)

type ProgressUpdate struct {
	Message string
	Done    bool
}

func LoadingSpinner(status chan ProgressUpdate) {
	loadSpinner := [10]string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	idx := 0
	currentMessage := "Fetching data"

	for {
		select {
		case state := <-status:
			if state.Done {
				erase := strings.Repeat(" ", len(currentMessage)*10)
				fmt.Printf("\r  %s\r", erase)
				fmt.Printf("%s", currentMessage)
			}
			currentMessage = state.Message

		default:
			erase := strings.Repeat(" ", len(currentMessage)*10)
			fmt.Printf("\r  %s\r", erase)
			fmt.Printf("%s %s", loadSpinner[idx%10], currentMessage)
			time.Sleep(50 * time.Millisecond)
			idx++
		}
	}
}
