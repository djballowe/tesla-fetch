package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type CallIssueCommandResponse struct {
	StatusCode int
	Body       string
}

type ResponseBody struct {
	Message string
	Success bool
}

func CallIssueCommand(command string) (CallIssueCommandResponse, error) {
	baseUrl := os.Getenv("BASE_URL")

	reqUrl := fmt.Sprintf("%s/command?command=%s", baseUrl, command)
	req, err := http.NewRequest(http.MethodPost, reqUrl, nil)
	if err != nil {
		return handleReturn("Error", 500, err)
	}

	req.Header.Add("Accept", "aplication/json")
	client := &http.Client{}

	resp, err := client.Do(req)
	if resp == nil || err != nil {
		return handleReturn("Error", 500, err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return handleReturn("Error", 500, err)
	}

	var responseBody ResponseBody
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return handleReturn("Error", 500, err)
	}

	if resp.StatusCode != 200 {
		return handleReturn(responseBody.Message, resp.StatusCode, nil)
	}

	return handleReturn(responseBody.Message, resp.StatusCode, nil)

}

func handleReturn(body string, statusCode int, err error) (CallIssueCommandResponse, error) {
	return CallIssueCommandResponse{
		StatusCode: statusCode,
		Body:       body,
	}, err
}
