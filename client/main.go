package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	drawlogo "tesla-app/client/draw-status"
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

	go func() {
		ui.LoadingSpinner()
	}()

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
	vehicleData, err := data.GetVehicleData()
	if err != nil {
		log.Fatalf("Could not get vehicle data: %s", err)
	}

	fmt.Printf("\r%s", "                                         ")
	drawlogo.DrawStatus(vehicleData)

	return
}
