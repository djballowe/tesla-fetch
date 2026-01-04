package vehiclecommand

import (
	"fmt"
	"tfetch/auth"
	"tfetch/ui"
	"tfetch/vehicle-state"
)

func IssueCommand(status ui.StatusLoggerMethods, token auth.Token, command string, vehicleService *vehicle.VehicleService) error {
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

	status.Log("Command %s issued successfully")

	return nil
}
