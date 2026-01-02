package drawstatus

import "tfetch/data"

type DrawService struct{}

type DrawMethods interface {
	DrawStatus(vehicleData *data.VehicleData)
	DrawStatusSimple(vehicleData *data.VehicleData) error
}
