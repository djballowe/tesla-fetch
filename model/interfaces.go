package model

type StatusLoggerMethods interface {
	Log(message string)
	Done()
}

type VehicleMethods interface {
	VehicleState(token Token) (*VehicleStateResponse, error)
}

type WakeMethods interface {
	Wake(token Token) (*WakeResponse, error)
	PollWake(token Token, status StatusLoggerMethods) error
}

type DrawMethods interface {
	DrawStatus(vehicleData *VehicleData)
	DrawStatusSimple(vehicleData *VehicleData) error
}
