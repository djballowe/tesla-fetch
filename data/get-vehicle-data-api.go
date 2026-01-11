package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"tfetch/auth"
	"tfetch/ui"
	"tfetch/vehicle-state"
)

func GetVehicleDataApi(status ui.StatusLoggerMethods, token auth.Token, vehicleDataService vehicle.VehicleMethods) (*VehicleData, error) {
	baseUrl := os.Getenv("TESLA_BASE_URL")
	carId := os.Getenv("MY_CAR_ID")

	apiResponse := &VehicleResponse{}
	vehicleData := &VehicleData{}

	vehicleState, err := vehicleDataService.VehicleState(token)
	if err != nil {
		return nil, err
	}

	if vehicleState.State != "online" {
		err := vehicleDataService.PollWake(token, status)
		if err != nil {
			return nil, err
		}
	}

	url := fmt.Sprintf("%s/vehicles/%s/vehicle_data", baseUrl, carId)

	client := &http.Client{}
	vehicleDataRequest, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	vehicleDataRequest.Header = http.Header{
		"Content-Type":  {"application/json"},
		"Authorization": {"Bearer " + token.AccessToken},
	}

	res, err := client.Do(vehicleDataRequest)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, errors.New("Failed to get vehicle data")
	}

	err = json.NewDecoder(res.Body).Decode(&apiResponse)
	if err != nil {
		return nil, err
	}

	vehicleData = &VehicleData{
		State:                apiResponse.Response.State,
		BatteryLevel:         apiResponse.Response.ChargeState.BatteryLevel,
		BatteryRange:         apiResponse.Response.ChargeState.BatteryRange,
		ChargeRate:           apiResponse.Response.ChargeState.ChargeRate,
		ChargingState:        apiResponse.Response.ChargeState.ChargingState,
		ChargeLimitSoc:       apiResponse.Response.ChargeState.ChargeLimitSoc,
		MinutesToFullCharge:  apiResponse.Response.ChargeState.MinutesToFullCharge,
		TimeToFullCharge:     apiResponse.Response.ChargeState.TimeToFullCharge,
		InsideTemp:           int(apiResponse.Response.ClimateState.InsideTemp),
		PassengerTempSetting: apiResponse.Response.ClimateState.PassengerTempSetting,
		DriverTempSetting:    apiResponse.Response.ClimateState.DriverTempSetting,
		IsClimateOn:          apiResponse.Response.ClimateState.IsClimateOn,
		IsPreconditioning:    apiResponse.Response.ClimateState.IsPreconditioning,
		OutsideTemp:          int(apiResponse.Response.ClimateState.OutsideTemp),
		Locked:               apiResponse.Response.VehicleState.Locked,
		Odometer:             int(apiResponse.Response.VehicleState.Odometer),
		Color:                apiResponse.Response.VehicleConfig.ExteriorColor,
		VehicleName:          apiResponse.Response.VehicleState.VehicleName,
		CarType:              apiResponse.Response.VehicleConfig.CarType,
		CarSpecialType:       apiResponse.Response.VehicleConfig.CarSpecialType,
	}

	vehicleJson, err := json.Marshal(vehicleData)
	if err != nil {
		return nil, err
	}

	filePath, err := getStateFilePath()
	if err != nil {
		return nil, err
	}

	err = os.WriteFile(filePath, vehicleJson, 0600)
	if err != nil {
		return nil, err
	}

	return vehicleData, nil
}
