package main

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"tesla-app/auth"
	vehiclecommand "tesla-app/command"
	"tesla-app/data"
	drawstatus "tesla-app/draw-status"
	"tesla-app/ui"
	"tesla-app/vehicle-state"

	"github.com/joho/godotenv"
)

// extract get data and issue command into their own services and abstract common api logic to its own package

var ErrorTooManyArgs = errors.New("usage too many args provided")
var ErrorInvalidArgs = errors.New("usage invalid arguments")

type AppDependencies struct {
	Status            chan ui.ProgressUpdate
	AuthService       auth.AuthMethods
	VehicleService    vehicle.VehicleMethods
	DrawStatusService drawstatus.DrawMethods
	IssueCommand      func(status chan ui.ProgressUpdate, token auth.Token, command string, vehicleService *vehicle.VehicleService) error
	GetData           func(status chan ui.ProgressUpdate, token auth.Token, vehicleDataService vehicle.VehicleMethods, flag string) (*data.VehicleData, error)
}

func main() {
	args := os.Args

	err := loadEnv()
	if err != nil {
		log.Fatalf("\rError loading env: %s\n", err)
	}

	status := make(chan ui.ProgressUpdate)
	defer close(status)

	go ui.LoadingSpinner(status)

	app := AppDependencies{
		Status:            status,
		AuthService:       &auth.AuthService{},
		VehicleService:    &vehicle.VehicleService{},
		DrawStatusService: &drawstatus.DrawService{},
		IssueCommand:      vehiclecommand.IssueCommand,
		GetData:           data.GetVehicleData,
	}

	err = runApp(app, args)
	if err != nil {
		log.Fatalf("\rApp run error: %s\n", err)
	}
	return
}

func runApp(app AppDependencies, args []string) error {
	flag, err := validateFlags(args)
	if err != nil {
		return err
	}

	vehicleService := app.VehicleService
	token, err := app.AuthService.CheckLogin(app.Status)
	if err != nil {
		return err
	}

	// no issue command logic right now focus on refactor
	// if len(args) == 2 {
	// 	app.Status <- ui.ProgressUpdate{Message: "Issuing command"}
	// 	command := args[1]
	//
	// 	err = app.IssueCommand(app.Status, *token, command, vehicleService)
	// 	if err != nil {
	// 		return err
	// 	}
	//
	// 	fmt.Printf("\rCommand %s issued successfully\n", command)
	// }

	app.Status <- ui.ProgressUpdate{Message: "Fetching data"}

	vehicleData, err := app.GetData(app.Status, *token, vehicleService, flag)
	if err != nil {
		return err
	}

	app.Status <- ui.ProgressUpdate{Done: true}
	if flag == "-w" {
		err = app.DrawStatusService.DrawStatusSimple(vehicleData)
		if err != nil {
			return err
		}
	} else {
		app.DrawStatusService.DrawStatus(vehicleData)
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

func validateFlags(args []string) (string, error) {
	// one flag allowed for now
	validFlags := map[string]bool{
		"-w": true,
	}
	argLimit := 2

	flag := ""

	if len(args) > argLimit {
		return "", ErrorTooManyArgs
	} else if !validFlags[args[1]] {
		return "", ErrorInvalidArgs
	}

	if len(args) > 0 && validFlags[args[1]] {
		flag = args[1]
	}

	return flag, nil
}
