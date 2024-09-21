package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type VehicleData struct {
	State               string `json:"state"`
	BatteryLevel        int    `json:"battery_level"`
	ChargeRate          int    `json:"charge_rate"`
	CharginState        string `json:"chargin_state"`
	MinutesToFullCharge int    `json:"minutes_to_full_charge"`
	TimeToFullCharge    int    `json:"time_to_full_charge"`
	InsideTemp          int    `json:"inside_temp"`
	IsClimateOn         bool   `json:"is_climate_on"`
	IsPreconditioning   bool   `json:"is_preconditioning"`
	OutsideTemp         int    `json:"outside_temp"`
	Locked              bool   `json:"locked"`
	Odometer            int    `json:"odometer"`
	ExteriorColor       string `json:"exterior_color"`
	VehicleName         string `json:"vehicle_name"`
	CarType             string `json:"car_type"`
	CarSpecialType      string `json:"car_special_type"`
}

type DataResponse struct {
	Body       VehicleData
	StatusCode int
}

func CallGetVehicleData() (DataResponse, error) {
	dataResponse := &DataResponse{}
	var vehicleData VehicleData

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

	if resp.StatusCode != 200 {
		return DataResponse{
			StatusCode: resp.StatusCode,
			Body:       vehicleData,
		}, nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return DataResponse{
			StatusCode: 500,
			Body:       vehicleData,
		}, err
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
