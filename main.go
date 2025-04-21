package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"tesla-app/auth"
	vehiclecommand "tesla-app/command"
	"tesla-app/data"
	"tesla-app/dependencies"
	drawlogo "tesla-app/draw-status"
	"tesla-app/ui"
	"tesla-app/vehicle-state"

	"github.com/joho/godotenv"
)

type AppDependencies struct {
	Status         chan ui.ProgressUpdate
	AuthService    *auth.AuthService
	VehicleService *vehicle.VehicleService
	DrawStaus      func(vehicleData *data.VehicleData)
	IssueCommand   func(status chan ui.ProgressUpdate, token auth.Token, command string, vehicleDataService *dependencies.VehicleDataService) error
	GetData        func(token auth.Token, status chan ui.ProgressUpdate, vehicleDataService *dependencies.VehicleDataService) (*data.VehicleData, error)
}

func main() {
	args := os.Args
	if len(args) > 1 {
		log.Fatalf("\rerror: %v", errors.New("can only issue one command"))
	}

	err := loadEnv()
	if err != nil {
		log.Fatalf("\rError loading env: %s\n", err)
	}

	status := make(chan ui.ProgressUpdate)
	defer close(status)

	go func() {
		ui.LoadingSpinner(status)
	}()

	if len(args) == 1 {
		setCommand(status, args[1])
	} else {
		setGetData(status)
	}

	app := AppDependencies{
		Status:         status,
		AuthService:    &auth.AuthService{},
		VehicleService: &vehicle.VehicleService{},
		DrawStaus:      drawlogo.DrawStatus,
		IssueCommand:   vehiclecommand.IssueCommand,
		GetData:        data.GetVehicleData,
	}

	err = runApp(app, args)
	if err != nil {
		log.Fatalf("\rApp run error: %s\n", err)
	}
	return
}

func runApp(app AppDependencies, args []string) error {
	return nil
}

func loadEnv() error {
	rootDir, err := os.Executable()
	if err != nil {
		return err
	}

	path := filepath.Dir(rootDir)
	envPath := filepath.Join(path, ".env")

	err = godotenv.Load(envPath)
	if err != nil {
		return err
	}

	return nil
}

func setCommand(status chan ui.ProgressUpdate, command string) {
	status <- ui.ProgressUpdate{Message: "Issuing command"}

	vehicleService := &vehicle.VehicleService{}
	service := dependencies.VehicleDataController(vehicleService)

	token, err := auth.CheckLogin(status)
	if err != nil {
		log.Fatalf("\rError fetching auth token: %s\n", err)
	}

	err = vehiclecommand.IssueCommand(status, *token, command, service)
	if err != nil {
		log.Fatalf("\rerror: %s\n", err)
	}

	status <- ui.ProgressUpdate{Done: true}
	fmt.Printf("Command \"%s\" issued successfully\n", command)
	return
}

func setGetData(status chan ui.ProgressUpdate) {
	status <- ui.ProgressUpdate{Message: "Fetching data"}

	vehicleService := &vehicle.VehicleService{}
	service := dependencies.VehicleDataController(vehicleService)

	token, err := auth.CheckLogin(status)
	if err != nil {
		log.Fatalf("\rError fetching auth token: %s\n", err)
	}

	vehicleData, err := data.GetVehicleData(*token, status, service)
	if err != nil {
		log.Fatalf("\rCould not get vehicle data: %s\n", err)
	}

	status <- ui.ProgressUpdate{Done: true}
	drawlogo.DrawStatus(vehicleData)
	return
}
