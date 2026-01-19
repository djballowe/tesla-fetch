package model

import (
	"errors"
	"time"
)

var ErrNoStateFile = errors.New("no state file exists")

type Token struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	IdToken      string    `json:"id_token"`
	State        string    `json:"state"`
	TokenType    string    `json:"token_type"`
	ExpiresIn    int       `json:"expires_in"`
	CreateAt     time.Time `json:"created_at"`
}

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

type VehicleStateResponse struct {
	State string
	Vin   string
}

type WakeResponse struct {
	State string
}

type CommandRequest struct {
	AuthToken string `json:"authToken"`
	Vin       string `json:"vin"`
	Command   string `json:"command"`
}

type CommandResponse struct {
	Success bool   `json:"status"`
	Message string `json:"message"`
}
