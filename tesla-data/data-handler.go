package tesla

import (
	"tfetch/auth"
	"tfetch/model"
)

func GetDataHandler(status model.StatusLoggerMethods, authService auth.AuthMethods, vehicleService model.VehicleMethods, wakeService model.WakeMethods) (*model.VehicleData, error) {
	token, err := authService.CheckLogin(status)
	if err != nil {
		return nil, err
	}

	vehicleData, err := GetVehicleData(status, *token, vehicleService, wakeService)
	if err != nil {
		return nil, err
	}

	return vehicleData, nil
}
