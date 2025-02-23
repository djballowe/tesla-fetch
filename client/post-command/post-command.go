package postcommand

import (
	"errors"
	"fmt"
	"tesla-app/client/auth"
	vehiclecommand "tesla-app/client/command"
	"tesla-app/client/ui"
)

func IssueCommand(status chan ui.ProgressUpdate, command string) error {
	token, err := auth.CheckLogin(status)
	response, err := vehiclecommand.CallIssueCommand(status, *token, command)
	if err != nil {
		return err
	}

	if response.Success == false {
		err := errors.New(fmt.Sprintf("Failed to issue command %s: %s", command, response.Message))
		return err
	}

	return nil
}
