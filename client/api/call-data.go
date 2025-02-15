package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"tesla-app/client/helpers"
	"tesla-app/client/ui"
	"tesla-app/client/vehicle-state"
)

func CallGetVehicleData(status chan ui.ProgressUpdate) (*VehicleData, error) {
	tokenStore, state := helpers.GetTokenStore()
	if state == "" || tokenStore[state].AccessToken == "" {
		return nil, errors.New("Invalid or no access token")
	}
	baseUrl := os.Getenv("TESLA_BASE_URL")
	carId := os.Getenv("MY_CAR_ID")

	var apiResponse = &VehicleResponse{}

	vehicleState, err := vehicle.VehicleState()
	status <- ui.ProgressUpdate{Message: fmt.Sprintf("Vehicle State: %s", vehicleState.State)}
	if err != nil {
		return nil, err
	}

	if vehicleState.State != "online" {
		err := vehicle.PollWake(status)
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
		"Authorization": {"Bearer " + tokenStore[state].AccessToken},
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

	vehicleData := &VehicleData{
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

	return vehicleData, nil
}
