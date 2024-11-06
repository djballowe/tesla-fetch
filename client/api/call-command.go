package api

import (
	"fmt"
	"io"
	"net/http"
)

type CallIssueCommandResponse struct {
	StatusCode int
	Body       string
}

func CallIssueCommand(command string) (CallIssueCommandResponse, error) {
	reqUrl := fmt.Sprintf("http://localhost:8080/command?command=%s", command)
	req, err := http.NewRequest(http.MethodPost, reqUrl, nil)
	if err != nil {
		return CallIssueCommandResponse{
			StatusCode: 500,
			Body:       "Error",
		}, err
	}

	req.Header.Add("Accept", "aplication/json")
	client := &http.Client{}

	resp, err := client.Do(req)
	if resp == nil || err != nil {
		return CallIssueCommandResponse{
			StatusCode: 500,
			Body:       "Error",
		}, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return CallIssueCommandResponse{
			StatusCode: 500,
			Body:       "Error",
		}, err
	}

	if resp.StatusCode != 200 {
		return CallIssueCommandResponse{
			StatusCode: resp.StatusCode,
			Body:       string(body),
		}, nil
	}

	return CallIssueCommandResponse{
		StatusCode: 200,
		Body:       string(body),
	}, nil
}
