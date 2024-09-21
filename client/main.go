package main

import (
	"fmt"
	"tesla-app/client/api"
	"tesla-app/client/draw-status"
)

func main() {
	carDataResponse, error := api.CallGetVehicleData()
	if error != nil {
		fmt.Println(error.Error())
		return
	}

	fmt.Println(carDataResponse.StatusCode)

	if carDataResponse.StatusCode == 401 {
		authResponse, error := api.CallAuth()
		if error != nil || authResponse.StatusCode != 200 {
			fmt.Println(error.Error())
			return
		}

		carDataResponse, error = api.CallGetVehicleData()
		if error != nil {
			fmt.Println(error.Error())
			return
		}
	}

	if carDataResponse.StatusCode != 200 {
		fmt.Printf("Something went wrong: %d\n", carDataResponse.StatusCode)
		return
	}

	drawlogo.DrawStatus(carDataResponse.Body)

	//	fmt.Println(carDataResponse.Body)

	return
}
