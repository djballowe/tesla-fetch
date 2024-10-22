package data 

import (
	"errors"
	"fmt"
	"sync"
	"tesla-app/client/api"
)

type DataResult struct {
	VehicleData api.VehicleData
	Err         error
}

func GetVehicleData(group *sync.WaitGroup, done chan struct{}, dataChan chan DataResult) {
	defer group.Done()
	defer close(done)
	carDataResponse, error := api.CallGetVehicleData()
	if error != nil {
		dataChan <- DataResult{Err: error}
		return
	}

	if carDataResponse.StatusCode == 401 {
		authResponse, error := api.CallAuth()
		if error != nil || authResponse.StatusCode != 200 {
			dataChan <- DataResult{Err: error}
			return
		}

		carDataResponse, error = api.CallGetVehicleData()
		if error != nil {
			dataChan <- DataResult{Err: error}
			return
		}
	}

	if carDataResponse.StatusCode != 200 {
		error = errors.New(fmt.Sprintf("Error gathering vehicle data: Status Code %d", carDataResponse.StatusCode))
		dataChan <- DataResult{Err: error}
		return
	}

	dataChan <- DataResult{VehicleData: carDataResponse.Body, Err: nil}
}
