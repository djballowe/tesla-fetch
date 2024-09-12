package api

import (
	"fmt"
	"io"
	"net/http"
)

type DataResponse struct {
	Body       string
	StatusCode int
}

func CallGetVehicleData() (DataResponse, error) {
	dataResponse := &DataResponse{}

	req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/data", nil)
	if err != nil {
		return DataResponse{
			StatusCode: 500,
			Body:       "Error",
		}, err
	}

	req.Header.Add("Accept", "aplication/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if resp == nil || err != nil {
		return DataResponse{
			StatusCode: 500,
			Body:       "Error",
		}, err
	}

	fmt.Println("car data")

	if resp.StatusCode == 408 {
		return DataResponse{
			StatusCode: 408,
			Body:       "Car is not awake",
		}, nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return DataResponse{
			StatusCode: 500,
			Body:       "Error",
		}, err
	}

	dataResponse.Body = string(body)
	dataResponse.StatusCode = resp.StatusCode

	return *dataResponse, nil
}
