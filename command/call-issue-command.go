package vehiclecommand

import (
	"fmt"
	"tesla-app/auth"
	"tesla-app/ui"
	"tesla-app/vehicle-state"
)

func CallIssueCommand(status chan ui.ProgressUpdate, token auth.Token, command string) error {
	vehicleState, err := vehicle.VehicleState(token)
	if err != nil {
		return err
	}

	state := vehicleState.State
	vin := vehicleState.Vin

	if state != "online" {
		err := vehicle.PollWake(token, status)
		if err != nil {
			return err
		}
	}

	status <- ui.ProgressUpdate{Message: fmt.Sprintf("Issuing command: %s", command)}

	var commandReq = CommandRequest{
		AuthToken: token.AccessToken,
		Vin:       vin,
		Command:   command,
	}

	err = HandleCommand(commandReq)
	if err != nil {
		return fmt.Errorf("could not issue command: %s", err.Error())
	}

	status <- ui.ProgressUpdate{Message: fmt.Sprintf("Command %s issued successfully", command)}

	return nil
}
