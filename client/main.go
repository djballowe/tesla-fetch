package main

import (
	"fmt"
	"os"
	"sync"
	"tesla-app/client/api"
	"tesla-app/client/draw-status"
	"tesla-app/client/ui"
	data "tesla-app/client/vehicle-data"
)

func main() {
	args := os.Args

	switch len(args) {
	case 1:
		setGetData()
		break
	case 2:
		setCommand(args[1])
		break
	default:
		fmt.Println("error: can only issue one command")
	}

	return
}

func setCommand(command string) {
	switch command {
	case "lock":
		fmt.Println("Locking car")
		_, err := api.CallIssueCommand("lock")
		if err != nil {
			fmt.Printf("error: %s", err.Error())
		}
		break
	case "unlock":
		fmt.Println("Unlocking car")
		break
	default:
		fmt.Println("error: not a valid command")
		break
	}
	return
}

func setGetData() {
	var group sync.WaitGroup
	group.Add(2)

	done := make(chan struct{})
	dataChan := make(chan data.DataResult)

	go ui.LoadingSpinner(&group, done)
	go data.GetVehicleData(&group, done, dataChan)
	res := <-dataChan
	group.Wait()

	if res.Err != nil {
		fmt.Printf("error: %s\n", res.Err.Error())
		return
	}

	drawlogo.DrawStatus(res.VehicleData)

	return
}
