package postcommand

import (
	"tesla-app/auth"
	vehiclecommand "tesla-app/command"
	"tesla-app/ui"
)

func IssueCommand(status chan ui.ProgressUpdate, command string) error {
	token, err := auth.CheckLogin(status)
	if err != nil {
		return err
	}

	err = vehiclecommand.CallIssueCommand(status, *token, command)
	if err != nil {
		return err
	}

	return nil
}
