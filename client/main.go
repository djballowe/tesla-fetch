package main

import (
	"fmt"
	"io"
	"net/http"
)

const BASE_URL = "localhost:8080"

type DataResponse struct {
	Body       string
	StatusCode int
}

type Routes struct {
	Data string
	Honk string
	Auth string
}

var routes = Routes{
	Data: "/data",
	Honk: "/honk",
	Auth: "/auth",
}

func main() {

	carDataResponse, error := callGetVehicleData()
	if error != nil {
		fmt.Println(error.Error())
		return
	}

	fmt.Println("Response: ", carDataResponse.Body, carDataResponse.StatusCode)

	if carDataResponse.StatusCode == 401 {
		callAuth()
		carDataResponse, error = callGetVehicleData()
		if error != nil {
			fmt.Println(error.Error())
			return
		}
	}

	return
}

func callGetVehicleData() (DataResponse, error) {
	response := &DataResponse{}

	req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/data", nil)
	if err != nil {
		return *response, err
	}

	req.Header.Add("Accept", "aplication/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if resp == nil || err != nil {
		return *response, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return *response, err
	}

	response.Body = string(body)
	response.StatusCode = resp.StatusCode

	return *response, nil
}

func callAuth() {
	req, err := http.NewRequest()
}
