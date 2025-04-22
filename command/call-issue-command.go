package vehiclecommand

import (
	"fmt"
	"tesla-app/auth"
	"tesla-app/ui"
	"tesla-app/vehicle-state"
)

func IssueCommand(status chan ui.ProgressUpdate, token auth.Token, command string, vehicleService *vehicle.VehicleService) error {
	vehicleState, err := vehicleService.VehicleState(token)
	if err != nil {
		return err
	}

	vin := vehicleState.Vin

	if vehicleState.State != "online" {
		err := vehicleService.PollWake(token, status)
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
