package ui 

import (
	"fmt"
	"sync"
	"time"
)

func LoadingSpinner(group *sync.WaitGroup, done chan struct{}) {
	defer group.Done()
	loadSpinner := [10]string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	idx := 0

	for {
		select {
		case <-done:
			fmt.Printf("\r%s", "                         \n")
			return

		default:
			fmt.Printf("\r%s Fetching vehicle data", loadSpinner[idx%10])
			time.Sleep(50 * time.Millisecond)
			idx++
		}
	}
}
