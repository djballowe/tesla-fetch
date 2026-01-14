package main

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	apitypes "tfetch/api/types"
	"tfetch/api/vehicle-status"
	"tfetch/api/wake"
	"tfetch/auth"
	vehiclecommand "tfetch/command"
	drawstatus "tfetch/draw-status"
	services "tfetch/services/get-data-handler"
	"tfetch/ui"

	"github.com/joho/godotenv"
)

var ErrorTooManyArgs = errors.New("too many args provided")
var ErrorInvalidArgs = errors.New("invalid arguments")

type AppDependencies struct {
	StatusLogger      ui.StatusLoggerMethods
	Flag              string
	AuthService       auth.AuthMethods
	VehicleService    apitypes.VehicleMethods
	WakeService       apitypes.WakeMethods
	DrawStatusService drawstatus.DrawMethods
	IssueCommand      func(status ui.StatusLoggerMethods, token auth.Token, command string, vehicleService apitypes.VehicleMethods, wakeService apitypes.WakeMethods) error
	GetDataHandler    func(status ui.StatusLoggerMethods, authService auth.AuthMethods, vehicleService apitypes.VehicleMethods, wakeService apitypes.WakeMethods) (*apitypes.VehicleData, error)
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

	var statusLogger ui.StatusLoggerMethods
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
		VehicleService:    &vehicle.VehicleService{},
		DrawStatusService: &drawstatus.DrawService{},
		WakeService:       &wake.WakeService{},
		IssueCommand:      vehiclecommand.IssueCommand,
		GetDataHandler:    services.GetDataHandler,
	}

	err = runApp(app)
	if err != nil {
		log.Fatalf("\r%s\n", err)
	}
	return
}

func runApp(app AppDependencies) error {
	vehicleService := app.VehicleService
	authService := app.AuthService
	wakeService := app.WakeService
	// token, err := app.AuthService.CheckLogin(app.StatusLogger)
	// if err != nil {
	// 	return err
	// }

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
	// one flag allowed for now
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
