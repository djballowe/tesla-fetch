package drawstatus

import (
	apitypes "tfetch/api/types"
)

type DrawService struct{}

type DrawMethods interface {
	DrawStatus(vehicleData *apitypes.VehicleData)
	DrawStatusSimple(vehicleData *apitypes.VehicleData) error
}
