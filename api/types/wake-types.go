package apitypes

import (
	"tfetch/auth"
	"tfetch/ui"
)

type WakeMethods interface {
	Wake(token auth.Token) (*WakeResponse, error)
	PollWake(token auth.Token, status ui.StatusLoggerMethods) (error)
}

type WakeResponse struct {
	State string
}

type TeslaVehicleWakeResponse struct {
	Response struct {
		State string `json:"state"`
	}
}
