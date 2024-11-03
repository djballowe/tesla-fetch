package api

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

func CallIssueCommand(command string) error {
	reqUrl := fmt.Sprintf("http://localhost:8080/command?command=%s", command)
	req, err := http.NewRequest(http.MethodPost, reqUrl, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Accept", "aplication/json")
	client := &http.Client{}

	resp, err := client.Do(req)
	if resp == nil || err != nil {
		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		err = errors.New(fmt.Sprintf(string(body)))
		return err
	}

	return nil
}
