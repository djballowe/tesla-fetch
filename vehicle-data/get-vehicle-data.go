package data

import (
	"tesla-app/auth"
	"tesla-app/data"
	"tesla-app/ui"
)

var vehicleMethods data.VehicleAPI

func GetVehicleData(status chan ui.ProgressUpdate) (*data.VehicleData, error) {
	token, err := auth.CheckLogin(status)

	carDataResponse, err := data.CallGetVehicleData(*token, status, vehicleMethods)
	if err != nil {
		return nil, err
	}

	return carDataResponse, nil
}
