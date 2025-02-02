package data

import (
	// "errors"
	// "fmt"
	"tesla-app/client/api"
)

type DataResult struct {
	VehicleData api.VehicleData
	Err         error
}

func GetVehicleData(done chan struct{}, dataChan chan DataResult, statusChan chan string) {
	statusChan <- "Fetching vehicle data"
	defer close(done)
	defer close(statusChan)

	authResponse, error := api.CallAuth()
	if error != nil || authResponse.StatusCode != 200 {
		dataChan <- DataResult{Err: error}
		return
	}

	// carDataResponse, error := api.CallGetVehicleData()
	// if error != nil {
	// 	dataChan <- DataResult{Err: error}
	// 	return
	// }
	//
	// if carDataResponse.StatusCode == 401 {
	// 	statusChan <- "Fetching authentication"
	// 	authResponse, error := api.CallAuth()
	// 	if error != nil || authResponse.StatusCode != 200 {
	// 		dataChan <- DataResult{Err: error}
	// 		return
	// 	}
	//
	// 	statusChan <- "Fetching vehicle data"
	// 	carDataResponse, error = api.CallGetVehicleData()
	// 	if error != nil {
	// 		dataChan <- DataResult{Err: error}
	// 		return
	// 	}
	// }
	//
	// if carDataResponse.StatusCode != 200 {
	// 	error = errors.New(fmt.Sprintf("Error gathering vehicle data: Status Code %d", carDataResponse.StatusCode))
	// 	dataChan <- DataResult{Err: error}
	// 	return
	// }

	// dataChan <- DataResult{VehicleData: carDataResponse.Body, Err: nil}
}
