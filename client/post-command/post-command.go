package postcommand

import (
	"tesla-app/client/api"
)

func PostCommand(command string) error {
	err := api.CallIssueCommand(command)
	if err != nil {
		return err
	}
	return nil
}
