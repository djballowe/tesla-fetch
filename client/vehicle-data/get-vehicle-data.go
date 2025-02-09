package data

import (
	// "errors"
	// "fmt"
	"log"
	"tesla-app/client/api"
	"tesla-app/client/auth"
)

type DataResult struct {
	VehicleData api.VehicleData
	Err         error
}

func GetVehicleData() (DataResult, error) {
	var dataResult = DataResult{}
	carDataResponse, err := api.CallGetVehicleData()
	if err != nil {
		return dataResult, err
	}
	log.Println("car data response: ", carDataResponse)

	return dataResult, nil
	// if err != nil {
	// 	dataChan <- DataResult{Err: err}
	// 	return
	// }
	//
	// if carDataResponse.StatusCode == 401 {
	// 	statusChan <- "Fetching authentication"
	// 	err := api.CallAuth()
	// 	if err != nil {
	// 		dataChan <- DataResult{Err: err}
	// 		return
	// 	}
	//
	// 	statusChan <- "Fetching vehicle data"
	// 	carDataResponse, err = api.CallGetVehicleData()
	// 	if err != nil {
	// 		dataChan <- DataResult{Err: err}
	// 		return
	// 	}
	// }
	//
	// if carDataResponse.StatusCode != 200 {
	// 	err = errors.New(fmt.Sprintf("Error gathering vehicle data: Status Code %d", carDataResponse.StatusCode))
	// 	dataChan <- DataResult{Err: err}
	// 	return
	// }
	//
}
