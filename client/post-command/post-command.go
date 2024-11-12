package postcommand

import (
	"errors"
	"fmt"
	"tesla-app/client/api"
)

func PostCommand(command string) error {
	response := api.CallIssueCommand(command)
	if response.Error != nil {
		err := response.Error
		return err
	}

	if response.StatusCode == 401 {
		authResp, err := api.CallAuth()
		if err != nil || authResp.StatusCode != 200 {
			return err
		}

		response = api.CallIssueCommand(command)
		if response.Error != nil {
			err = response.Error
			return response.Error
		}
	}

	if response.Status != 200 {
		err := errors.New(fmt.Sprintf("Failed to issue command: %s Status Code: %d", response.Message, response.Status))
		return err
	}


	return nil

}
