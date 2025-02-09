package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"sync"

	//"tesla-app/client/draw-status"
	drawlogo "tesla-app/client/draw-status"
	postcommand "tesla-app/client/post-command"
	//	"tesla-app/client/ui"
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
		log.Fatalf("error: %v", errors.New("can only issue one command"))
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

	vehicleData, err := data.GetVehicleData()
	if err != nil {
		log.Fatalf("Could not get vehicle data: %s", err)
	}


	// go func() {
	// 	ui.LoadingSpinner(done)
	// }()

	// fmt.Printf("\r%s", "                                         ")
	drawlogo.DrawStatus(vehicleData)

	return
}
