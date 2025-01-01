package ui

import (
	"fmt"
	"time"
)

func LoadingSpinner(done chan struct{}) {
	loadSpinner := [10]string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	idx := 0

	for {
		select {
		case <-done:
			return

		default:
			fmt.Printf("\r%s", loadSpinner[idx%10])
			time.Sleep(50 * time.Millisecond)
			idx++
		}
	}
}
