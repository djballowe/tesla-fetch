package main

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"tfetch/auth"
	"tfetch/command"
	"tfetch/model"
	"tfetch/tesla-data"
	"tfetch/ui"

	"github.com/joho/godotenv"
)

var ErrorTooManyArgs = errors.New("too many args provided")
var ErrorInvalidArgs = errors.New("invalid arguments")

type AppDependencies struct {
	StatusLogger      model.StatusLoggerMethods
	Flag              string
	AuthService       auth.AuthMethods
	VehicleService    model.VehicleMethods
	WakeService       model.WakeMethods
	DrawStatusService model.DrawMethods
	IssueCommand      func(status model.StatusLoggerMethods, token model.Token, command string, vehicleService model.VehicleMethods, wakeService model.WakeMethods) error
	GetDataHandler    func(status model.StatusLoggerMethods, authService auth.AuthMethods, vehicleService model.VehicleMethods, wakeService model.WakeMethods) (*model.VehicleData, error)
}

func main() {
	args := os.Args

	err := loadEnv()
	if err != nil {
		log.Fatalf("\rError loading env: %s\n", err)
	}

	flag, err := validateFlags(args)
	if err != nil {
		log.Fatalf("\rUsage: %s", err)
	}

	var statusLogger model.StatusLoggerMethods
	if flag == "-w" {
		statusLogger = ui.NewNoopLogger()
	} else {
		status := make(chan ui.ProgressUpdate)
		statusLogger = ui.NewStatusLogger(status)
		defer close(status)
		go ui.LoadingSpinner(status)
	}

	app := AppDependencies{
		StatusLogger:      statusLogger,
		Flag:              flag,
		AuthService:       &auth.AuthService{},
		VehicleService:    &tesla.VehicleService{},
		DrawStatusService: &ui.DrawService{},
		WakeService:       &tesla.WakeService{},
		IssueCommand:      command.IssueCommand,
		GetDataHandler:    tesla.GetDataHandler,
	}

	err = runApp(app)
	if err != nil {
		log.Fatalf("\r%s\n", err)
	}
}

func runApp(app AppDependencies) error {
	vehicleService := app.VehicleService
	authService := app.AuthService
	wakeService := app.WakeService

	vehicleData, err := app.GetDataHandler(app.StatusLogger, authService, vehicleService, wakeService)
	if err != nil {
		return err
	}

	app.StatusLogger.Done()
	if app.Flag == "-w" {
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
	validFlags := map[string]bool{
		"-w": true,
	}
	argLimit := 2

	flag := ""

	if len(args) > argLimit {
		return "", ErrorTooManyArgs
	} else if len(args) > 1 && !validFlags[args[1]] {
		return "", ErrorInvalidArgs
	}

	if len(args) > 1 && validFlags[args[1]] {
		flag = args[1]
	}

	return flag, nil
}
