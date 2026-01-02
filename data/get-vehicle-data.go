package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"tesla-app/auth"
	"tesla-app/ui"
	"tesla-app/vehicle-state"
)

var ErrNoStateFile = errors.New("no state file exists")

func GetVehicleData(status chan ui.ProgressUpdate, token auth.Token, vehicleDataService vehicle.VehicleMethods, flag string) (*VehicleData, error) {
	baseUrl := os.Getenv("TESLA_BASE_URL")
	carId := os.Getenv("MY_CAR_ID")

	apiResponse := &VehicleResponse{}
	vehicleData := &VehicleData{}

	vehicleState, err := vehicleDataService.VehicleState(token)
	if err != nil {
		return nil, err
	}

	if vehicleState.State != "online" {
		if flag == "-w" {
			err = getState(vehicleData)
			if err != nil {
				if errors.Is(err, ErrNoStateFile) {
					fmt.Println("State file does not exist waking car to initialize state")
				} else {
					return nil, err
				}
			} else {
				return vehicleData, nil
			}
		}
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

	fmt.Println("here")

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

func getStateFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	localDir := filepath.Join(homeDir, ".local", "share", "tfetch")

	err = os.MkdirAll(localDir, 0700)
	if err != nil {
		return "", err
	}

	statePath := filepath.Join(localDir, "vehicle-state.json")

	return statePath, nil
}

func getState(vehicleData *VehicleData) error {
	statePath, err := getStateFilePath()
	if err != nil {
		return err
	}

	data, err := os.ReadFile(statePath)
	if err != nil {
		if os.IsNotExist(err) {
			return ErrNoStateFile
		}
		return err
	}

	err = json.Unmarshal(data, vehicleData)
	if err != nil {
		return err
	}

	return nil
}
