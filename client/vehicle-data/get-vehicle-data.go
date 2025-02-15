package data

import (
	"tesla-app/client/api"
	"tesla-app/client/auth"
)

func GetVehicleData() (*api.VehicleData, error) {
	err := auth.CallAuth()
	if err != nil {
		return nil, err
	}

	carDataResponse, err := api.CallGetVehicleData()
	if err != nil {
		return nil, err
	}

	return carDataResponse, nil
}
