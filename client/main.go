package main

import (
	"fmt"
	"os"
	"sync"
	"tesla-app/client/draw-status"
	postcommand "tesla-app/client/post-command"
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
	err := postcommand.PostCommand(command)
	if err != nil {
		fmt.Printf("error: %s\n", err)
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
