package vehiclecommand

import (
	"errors"
	"fmt"
	"tesla-app/auth"
	"tesla-app/ui"
	"tesla-app/vehicle-state"
)

func CallIssueCommand(status chan ui.ProgressUpdate, token auth.Token, command string) (*CommandResponse, error) {
	vehicleState, err := vehicle.VehicleState(token)
	if err != nil {
		return nil, err
	}

	state := vehicleState.State
	vin := vehicleState.Vin

	if state != "online" {
		err := vehicle.PollWake(token, status)
		if err != nil {
			return nil, err
		}
	}

	status <- ui.ProgressUpdate{Message: fmt.Sprintf("Issuing command: %s", command)}

	var commandReq = CommandRequest{
		AuthToken: token.AccessToken,
		Vin:       vin,
		Command:   command,
	}

	commandResp := HandleCommand(commandReq)

	if !commandResp.Success {
		return nil, errors.New("Could not issue command")
	}

	status <- ui.ProgressUpdate{Message: fmt.Sprintf("Command %s issued successfully", command)}

	return &commandResp, nil
}
