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

type AuthResponse struct {
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

	if carDataResponse.StatusCode == 401 {
		authResponse, error := callAuth()
		fmt.Println("auth response", authResponse.StatusCode)
		if error != nil || authResponse.StatusCode != 200 {
			fmt.Println(error.Error())
			return
		}

		carDataResponse, error = callGetVehicleData()
		fmt.Println("second try response: ", carDataResponse.Body)
		if error != nil {
			fmt.Println(error.Error())
			return
		}
	}

	return
}

func callGetVehicleData() (DataResponse, error) {
	dataResponse := &DataResponse{}

	req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/data", nil)
	if err != nil {
		return *dataResponse, err
	}

	req.Header.Add("Accept", "aplication/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if resp == nil || err != nil {
		return *dataResponse, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return *dataResponse, err
	}

	dataResponse.Body = string(body)
	dataResponse.StatusCode = resp.StatusCode

	return *dataResponse, nil
}

func callAuth() (AuthResponse, error) {
	// this needs to redirect the user to login screen
	authResponse := &AuthResponse{}

	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/auth", nil)
	if err != nil {
		return *authResponse, err
	}

	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return *authResponse, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return *authResponse, err
	}

	authResponse.Body = string(body)
	authResponse.StatusCode = resp.StatusCode

	return *authResponse, nil
}
