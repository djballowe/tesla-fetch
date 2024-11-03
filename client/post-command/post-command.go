package postcommand

import (
	"errors"
	"fmt"
	"sync"
	"tesla-app/client/api"
)

func PostCommand(command string, group *sync.WaitGroup, done chan struct{}, errChan chan error) {
	defer group.Done()
	defer close(done)
	var err error = nil

	switch command {
	case "lock":
		fmt.Println("Locking car")
		err = api.CallIssueCommand("lock")
		if err != nil {
			errChan <- err
		}
		break
	case "unlock":
		fmt.Println("Unlocking car")
		err = api.CallIssueCommand("unlock")
		if err != nil {
			errChan <- err
		}
		break
	default:
		err = errors.New(fmt.Sprintf("Not a valid command"))
		errChan <- err
		break
	}

	return
}
