package postcommand

import (
	"errors"
	"fmt"
	"tesla-app/client/api"
	"tesla-app/client/auth"
	"tesla-app/client/ui"
)

func PostCommand(status chan ui.ProgressUpdate, command string) error {
	response, err := api.CallIssueCommand(command)
	if err != nil {
		return err
	}

	if response.StatusCode == 401 {
		// change to use new auth
		_, err := auth.CallAuth()
		if err != nil {
			return err
		}

		response, err = api.CallIssueCommand(command)
		if err != nil {
			return err
		}
	}

	if response.StatusCode != 200 {
		err := errors.New(fmt.Sprintf("Failed to issue command: %s Status Code: %d", response.Body, response.StatusCode))
		return err
	}

	fmt.Println(response.Body)

	return nil
}
