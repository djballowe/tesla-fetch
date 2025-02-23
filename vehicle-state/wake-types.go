package vehicle

type WakeResponse struct {
	State string
}

type TeslaVehicleWakeResponse struct {
	Response struct {
		State string `json:"state"`
	}
}
