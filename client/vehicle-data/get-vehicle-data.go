package data

import (
	"tesla-app/client/auth"
	"tesla-app/client/data"
	"tesla-app/client/ui"
)

func GetVehicleData(status chan ui.ProgressUpdate) (*data.VehicleData, error) {
	token, err := auth.CheckLogin(status)

	carDataResponse, err := data.CallGetVehicleData(*token, status)
	if err != nil {
		return nil, err
	}

	return carDataResponse, nil
}
