package data

import (
	"fmt"
	"tfetch/auth"
	"tfetch/ui"
	"tfetch/vehicle-state"
)

func GetDataHandler(status ui.StatusLoggerMethods, authService auth.AuthMethods, vehicleDataService vehicle.VehicleMethods, flag string) (*VehicleData, error) {
	vehicleData := &VehicleData{}

	if flag == "-w" {
		err := getDataCache(vehicleDataService, vehicleData)
		if err != nil {
			// need error hanlding for no json file
			return nil, err
		}
		return vehicleData, nil
	}

	fmt.Println("shouldnt see this")

	status.Log("Fetching Data")

	token, err := authService.CheckLogin(status)
	if err != nil {
		return nil, err
	}

	vehicleData, err = GetVehicleDataApi(status, *token, vehicleDataService)
	if err != nil {
		return nil, err
	}

	return vehicleData, nil
}
