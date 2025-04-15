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
	state := <-status

	for {
		select {
		case state := <-status:
			if state.Done {
				fmt.Printf("\r                                            ")
				fmt.Printf("\r")
				return
			}

		default:
			fmt.Printf("\r %s %s  ", loadSpinner[idx%10], state.Message)
			time.Sleep(40 * time.Millisecond)
			idx++
		}
	}
}
