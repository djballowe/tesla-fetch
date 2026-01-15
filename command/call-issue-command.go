package vehiclecommand

import (
	"fmt"
	apitypes "tfetch/api/types"
	"tfetch/auth"
	"tfetch/ui"
)

func IssueCommand(status ui.StatusLoggerMethods, token auth.Token, command string, vehicleService apitypes.VehicleMethods, wakeService apitypes.WakeMethods) error {
	vehicleState, err := vehicleService.VehicleState(token)
	if err != nil {
		return err
	}

	vin := vehicleState.Vin

	if vehicleState.State != "online" {
		err := wakeService.PollWake(token, status)
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
