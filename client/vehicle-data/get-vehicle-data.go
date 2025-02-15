package data

import (
	"tesla-app/client/api"
	"tesla-app/client/auth"
	"tesla-app/client/ui"
)

func GetVehicleData(status chan ui.ProgressUpdate) (*api.VehicleData, error) {
	err := auth.CallAuth()
	if err != nil {
		return nil, err
	}

	carDataResponse, err := api.CallGetVehicleData(status)
	if err != nil {
		return nil, err
	}

	return carDataResponse, nil
}
