package data

import (
	"tesla-app/auth"
	"tesla-app/data"
	"tesla-app/ui"
)

func GetVehicleData(status chan ui.ProgressUpdate) (*data.VehicleData, error) {
	token, err := auth.CheckLogin(status)

	carDataResponse, err := data.CallGetVehicleData(*token, status)
	if err != nil {
		return nil, err
	}

	return carDataResponse, nil
}
