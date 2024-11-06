package postcommand

import (
	"errors"
	"fmt"
	"tesla-app/client/api"
)

func PostCommand(command string) error {
	response, err := api.CallIssueCommand(command)
	if err != nil {
		return err
	}

	if response.StatusCode == 401 {
		authResp, err := api.CallAuth()
		if err != nil || authResp.StatusCode != 200 {
			return err
		}

		response, err = api.CallIssueCommand(command)
		if err != nil {
			return err
		}
	}

	if response.StatusCode != 200 {
		return errors.New(fmt.Sprintf("Error issuing command: %s\n", response.Body))
	}

	return nil
}
