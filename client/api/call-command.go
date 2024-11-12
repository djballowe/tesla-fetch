package api

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

type CallIssueCommandResponse struct {
	Message string
	Status  int
	Error   error
}

func CallIssueCommand(command string) CallIssueCommandResponse {
	reqUrl := fmt.Sprintf("http://localhost:8080/command?command=%s", command)
	req, err := http.NewRequest(http.MethodPost, reqUrl, nil)
	if err != nil {
		return handleResponse("Error creating call issue command request", 500, err)
	}

	req.Header.Add("Accept", "aplication/json")
	client := &http.Client{}

	resp, err := client.Do(req)
	if resp == nil || err != nil {
		return handleResponse("Error calling call issue command client", 500, err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return handleResponse("Error reading body", 500, err)
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 401 {
			return handleResponse(string(body), resp.StatusCode, nil)
		}
		err = errors.New(fmt.Sprintf(string(body)))
		return handleResponse("Error reading body", resp.StatusCode, err)
	}

	return handleResponse(string(body), resp.StatusCode, err)
}

func handleResponse(message string, statusCode int, error error) CallIssueCommandResponse {
	return CallIssueCommandResponse{
		Message: message,
		Status:  statusCode,
		Error:   error,
	}
}
