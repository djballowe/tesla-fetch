package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"tesla-app/auth"
	drawlogo "tesla-app/draw-status"
	postcommand "tesla-app/post-command"
	"tesla-app/ui"
	"tesla-app/vehicle-state"

	"github.com/joho/godotenv"
)

func main() {
	args := os.Args

	rootDir, err := os.Executable()
	if err != nil {
		log.Fatal("/rError loading .env file")
	}

	path := filepath.Dir(rootDir)
	envPath := filepath.Join(path, ".env")

	err = godotenv.Load(envPath)
	if err != nil {
		log.Fatal("\rError loading .env file")
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
		log.Fatalf("\rerror: %v", errors.New("can only issue one command"))
	}

	return
}

func setCommand(status chan ui.ProgressUpdate, command string) {
	status <- ui.ProgressUpdate{Message: "Issuing command"}

	err := postcommand.IssueCommand(status, command)
	if err != nil {
		log.Fatalf("\rerror: %s\n", err)
	}

	status <- ui.ProgressUpdate{Done: true}
	fmt.Printf("Command \"%s\" issued successfully\n", command)
	return
}

type VehicleDataController struct {
	vehicleMethods vehicle.VehicleApi
	authMethods    auth.AuthApi
}

func vehicleDataController(vehicleMethods vehicle.VehicleApi, authMethods auth.AuthApi) *VehicleDataController {
	return &VehicleDataController{
		vehicleMethods: vehicleMethods,
		authMethods:    authMethods,
	}
}

func setGetData(status chan ui.ProgressUpdate) {
	status <- ui.ProgressUpdate{Message: "Fetching data"}

	// get the interface to work
	vehicleClient := vehicle.VehicleApi
	authClient := auth.AuthApi
	service := vehicleDataController(vehicleClient, authClient)

	vehicleData, err := data.GetVehicleData(status, service)
	if err != nil {
		log.Fatalf("\rCould not get vehicle data: %s\n", err)
	}

	status <- ui.ProgressUpdate{Done: true}
	drawlogo.DrawStatus(vehicleData)
	return
}
