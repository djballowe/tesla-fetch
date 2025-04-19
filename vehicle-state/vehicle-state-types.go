package vehicle

import (
	"tesla-app/auth"
	"tesla-app/ui"
)

type VehicleApi interface {
	VehicleState(token auth.Token) (*VehicleStateResponse, error)
	PollWake(token auth.Token, status chan ui.ProgressUpdate) error
}

type VehicleStateResponse struct {
	State string
	Vin   string
}

type TeslaVehicleApiResponse struct {
	Response struct {
		ID             int64  `json:"id"`
		UserID         int64  `json:"user_id"`
		VehicleID      int64  `json:"vehicle_id"`
		Vin            string `json:"vin"`
		Color          any    `json:"color"`
		AccessType     string `json:"access_type"`
		GranularAccess struct {
			HidePrivate bool `json:"hide_private"`
		} `json:"granular_access"`
		Tokens                 any    `json:"tokens"`
		State                  string `json:"state"`
		InService              bool   `json:"in_service"`
		IDS                    string `json:"id_s"`
		CalendarEnabled        bool   `json:"calendar_enabled"`
		APIVersion             int    `json:"api_version"`
		BackseatToken          any    `json:"backseat_token"`
		BackseatTokenUpdatedAt any    `json:"backseat_token_updated_at"`
		BleAutopairEnrolled    bool   `json:"ble_autopair_enrolled"`
	} `json:"response"`
}
