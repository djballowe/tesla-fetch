package main

import (
	"fmt"
	"tesla-app/client/api"
)

func main() {
	carDataResponse, error := api.CallGetVehicleData()
	if error != nil {
		fmt.Println(error.Error())
		return
	}

	if carDataResponse.StatusCode == 401 {
		authResponse, error := api.CallAuth()
		fmt.Println("auth response", authResponse.StatusCode)
		if error != nil || authResponse.StatusCode != 200 {
			fmt.Println(error.Error())
			return
		}

		carDataResponse, error = api.CallGetVehicleData()
		fmt.Println("second try response: ", carDataResponse.Body)
		if error != nil {
			fmt.Println(error.Error())
			return
		}
	}

	return
}
