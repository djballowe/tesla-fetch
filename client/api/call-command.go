package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
	reqUrl := fmt.Sprintf("http://localhost:8080/command?command=%s", command)
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
