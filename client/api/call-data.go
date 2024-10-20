package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type VehicleData struct {
	State                string  `json:"state"`
	BatteryLevel         int     `json:"battery_level"`
	BatteryRange         float64 `json:"battery_range"`
	ChargeRate           float64 `json:"charge_rate"`
	ChargingState        string  `json:"charging_state"`
	ChargeLimitSoc       int     `json:"charge_limit_soc"`
	MinutesToFullCharge  int     `json:"minutes_to_full_charge"`
	TimeToFullCharge     float64 `json:"time_to_full_charge"`
	InsideTemp           int     `json:"inside_temp"`
	PassengerTempSetting float64 `json:"passenger_temp_setting"`
	DriverTempSetting    float64 `json:"driver_temp_setting"`
	IsClimateOn          bool    `json:"is_climate_on"`
	IsPreconditioning    bool    `json:"is_preconditioning"`
	OutsideTemp          int     `json:"outside_temp"`
	Locked               bool    `json:"locked"`
	Odometer             int     `json:"odometer"`
	Color                string  `json:"exterior_color"`
	VehicleName          string  `json:"vehicle_name"`
	CarType              string  `json:"car_type"`
	CarSpecialType       string  `json:"car_special_type"`
}

type DataResponse struct {
	Body       VehicleData
	StatusCode int
}

func CallGetVehicleData() (DataResponse, error) {
	dataResponse := &DataResponse{}
	var vehicleData VehicleData

	// TODO handle errors better here

	req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/data", nil)
	if err != nil {
		return DataResponse{
			StatusCode: 500,
			Body:       vehicleData,
		}, err
	}

	req.Header.Add("Accept", "aplication/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if resp == nil || err != nil {
		return DataResponse{
			StatusCode: 500,
			Body:       vehicleData,
		}, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return DataResponse{
			StatusCode: 500,
			Body:       vehicleData,
		}, err
	}

	if resp.StatusCode != 200 {
		return DataResponse{
			StatusCode: resp.StatusCode,
			Body:       vehicleData,
		}, nil
	}

	err = json.Unmarshal(body, &vehicleData)
	if err != nil {
		return DataResponse{
			StatusCode: 500,
			Body:       vehicleData,
		}, errors.New(fmt.Sprintf("Error parsing json: %s", err))
	}

	dataResponse.Body = vehicleData
	dataResponse.StatusCode = resp.StatusCode

	return *dataResponse, nil
}
