package ui

import (
	"fmt"
	"time"
	// "strings"
)

type ProgressUpdate struct {
	Message string
	Done    bool
}

func LoadingSpinner(status chan ProgressUpdate) {
	loadSpinner := [10]string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	idx := 0
	curMessage := "Fetching data"

	for {
		select {
		case state := <-status:
			if state.Done {
				fmt.Printf("\r                                            ")
				return
			}

		default:
			fmt.Printf("\r %s %s", loadSpinner[idx%10], curMessage)
			time.Sleep(40 * time.Millisecond)
			idx++
		}
	}
}
