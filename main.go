package main

import (
	"errors"
	"log"
	"os"
	drawlogo "tesla-app/draw-status"
	postcommand "tesla-app/post-command"
	"tesla-app/ui"
	data "tesla-app/vehicle-data"

	"github.com/joho/godotenv"
)

func main() {
	args := os.Args
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	status := make(chan ui.ProgressUpdate)
	defer close(status)

	go func() {
		ui.LoadingSpinner(status)
	}()

	switch len(args) {
	case 1:
		setGetData(status)
		break
	case 2:
		setCommand(status, args[1])
		break
	default:
		log.Fatalf("error: %v", errors.New("can only issue one command"))
	}

	return
}

func setCommand(status chan ui.ProgressUpdate, command string) {
	err := postcommand.IssueCommand(status, command)
	if err != nil {
		log.Fatalf("error: %s\n", err)
	}

	return
}

func setGetData(status chan ui.ProgressUpdate) {
	vehicleData, err := data.GetVehicleData(status)
	if err != nil {
		log.Fatalf("Could not get vehicle data: %s", err)
	}

	status <- ui.ProgressUpdate{Done: true}
	drawlogo.DrawStatus(vehicleData)
	return
}
