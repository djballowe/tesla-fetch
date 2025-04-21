package vehiclecommand

import (
	"fmt"
	"tesla-app/auth"
	"tesla-app/dependencies"
	"tesla-app/ui"
)

func IssueCommand(status chan ui.ProgressUpdate, token auth.Token, command string, vehicleDataService *dependencies.VehicleDataService) error {
	vehicleState, err := vehicleDataService.VehicleMethods.VehicleState(token)
	if err != nil {
		return err
	}

	vin := vehicleState.Vin

	if vehicleState.State != "online" {
		err := vehicleDataService.VehicleMethods.PollWake(token, status)
		if err != nil {
			return err
		}
	}

	var commandReq = CommandRequest{
		AuthToken: token.AccessToken,
		Vin:       vin,
		Command:   command,
	}

	err = HandleCommand(commandReq)
	if err != nil {
		return fmt.Errorf("\rcould not issue command: %s", err.Error())
	}

	status <- ui.ProgressUpdate{Message: fmt.Sprintf("Command %s issued successfully", command)}

	return nil
}
