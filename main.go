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
	IssueCommand   func(status chan ui.ProgressUpdate, token auth.Token, command string, vehicleService *vehicle.VehicleService) error
	GetData        func(status chan ui.ProgressUpdate, token auth.Token, vehicleDataService vehicle.VehicleMethods) (*data.VehicleData, error)
}

func main() {
	args := os.Args
	if len(args) > 2 {
		log.Fatalf("\rerror: %v", errors.New("can only issue one command"))
	}

	err := loadEnv()
	if err != nil {
		log.Fatalf("\rError loading env: %s\n", err)
	}

	status := make(chan ui.ProgressUpdate)
	defer close(status)

	go ui.LoadingSpinner(status)

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
	vehicleService := app.VehicleService
	token, err := app.AuthService.CheckLogin(app.Status)
	if err != nil {
		return err
	}

	if len(args) == 2 {
		app.Status <- ui.ProgressUpdate{Message: "Issuing command"}
		command := args[1]

		err = app.IssueCommand(app.Status, *token, command, vehicleService)
		if err != nil {
			return err
		}

		fmt.Printf("\rCommand %s issued successfully\n", command)
	} else {
		app.Status <- ui.ProgressUpdate{Message: "Fetching data"}

		vehicleData, err := app.GetData(app.Status, *token, vehicleService)
		if err != nil {
			return err
		}

		app.Status <- ui.ProgressUpdate{Done: true}
		app.DrawStaus(vehicleData)
	}

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
