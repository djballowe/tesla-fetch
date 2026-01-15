package services

import (
	api "tfetch/api/data"
	apitypes "tfetch/api/types"
	"tfetch/auth"
	"tfetch/ui"
)

func GetDataHandler(status ui.StatusLoggerMethods, authService auth.AuthMethods, vehicleService apitypes.VehicleMethods, wakeService apitypes.WakeMethods) (*apitypes.VehicleData, error) {
	token, err := authService.CheckLogin(status)
	if err != nil {
		return nil, err
	}

	vehicleData, err := api.GetVehicleDataApi(status, *token, vehicleService, wakeService)
	if err != nil {
		return nil, err
	}

	return vehicleData, nil
}
