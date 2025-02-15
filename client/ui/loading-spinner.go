package ui

import (
	"fmt"
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
				fmt.Printf("\r%s", "                                         ")
			}
			currentMessage = state.Message

		default:
			fmt.Printf("\r%s %s", loadSpinner[idx%10], currentMessage)
			time.Sleep(50 * time.Millisecond)
			idx++
		}
	}
}
