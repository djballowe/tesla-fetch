package drawstatus

import "tesla-app/data"

type DrawService struct{}

type DrawMethods interface {
	DrawStatus(vehicleData *data.VehicleData)
	DrawStatusSimple(vehicleData *data.VehicleData) error
}
