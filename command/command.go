package command

import (
	"fmt"
	"tfetch/model"
)

func IssueCommand(status model.StatusLoggerMethods, token model.Token, command string, vehicleService model.VehicleMethods, wakeService model.WakeMethods) error {
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

	var commandReq = model.CommandRequest{
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
