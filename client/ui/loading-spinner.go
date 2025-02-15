package ui

import (
	"fmt"
	"time"
)

func LoadingSpinner() {
	loadSpinner := [10]string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	idx := 0

	for {
		select {
		case <-ctx.Done():
			return

		default:
			fmt.Printf("\rFetching vehicle data %s\n", loadSpinner[idx%10])
			time.Sleep(50 * time.Millisecond)
			idx++
		}
	}
}
