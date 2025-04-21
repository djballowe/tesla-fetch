package dependencies

import (
	"tesla-app/vehicle-state"
)

// vehicle data service
type VehicleDataService struct {
	VehicleMethods vehicle.VehicleMethods
}

func VehicleDataController(vehicleMethods vehicle.VehicleMethods) *VehicleDataService {
	return &VehicleDataService{
		VehicleMethods: vehicleMethods,
	}
}
