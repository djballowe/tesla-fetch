package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"tesla-app/client/draw-status"
	postcommand "tesla-app/client/post-command"
	"tesla-app/client/ui"
	data "tesla-app/client/vehicle-data"

	"github.com/joho/godotenv"
)

func main() {
	args := os.Args

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

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
	group.Add(1)

	done := make(chan struct{})
	dataChan := make(chan data.DataResult)

	go func() {
		defer group.Done()
		data.GetVehicleData(done, dataChan)
	}()

	go func() {
		ui.LoadingSpinner(done)
	}()

	res := <-dataChan
	group.Wait()

	if res.Err != nil {
		fmt.Printf("error: %s\n", res.Err.Error())
		return
	}

	fmt.Printf("\r%s", "                                         ")
	drawlogo.DrawStatus(res.VehicleData)

	return
}
