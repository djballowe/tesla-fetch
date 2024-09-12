package main

import (
	"fmt"
	"tesla-app/client/api"
	drawlogo "tesla-app/client/draw-status"
)

func main() {
	carDataResponse, error := api.CallGetVehicleData()
	if error != nil {
		fmt.Println(error.Error())
		return
	}

	fmt.Println(carDataResponse.StatusCode)

	if carDataResponse.StatusCode == 401 {
		fmt.Println("This shouldnt fire")
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

	drawlogo.DrawStatus()

	if carDataResponse.StatusCode != 200 {
		fmt.Printf("Something went wrong: %s\n", carDataResponse.Body)
		return
	}

	fmt.Println(carDataResponse.Body)

	return
}
