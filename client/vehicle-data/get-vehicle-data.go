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

func GetVehicleData() (DataResult, error) {
	var dataResult = DataResult{}
	carDataResponse, err := api.CallGetVehicleData()
	if err != nil {
		return dataResult, nil
	}

	err = api.CallAuth()
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
	return dataResult, nil
}
