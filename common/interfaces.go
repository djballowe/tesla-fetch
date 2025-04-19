package interfacs

import (
	"tesla-app/auth"
	"tesla-app/vehicle-state"
)

// vehicle data service
type VehicleDataServcice struct {
	VehicleMethods vehicle.VehicleMethods
	AuthMethods    auth.AuthMethods
}

func VehicleDataController(vehicleMethods vehicle.VehicleMethods, authMethods auth.AuthMethods) *VehicleDataServcice {
	return &VehicleDataServcice{
		VehicleMethods: vehicleMethods,
		AuthMethods:    authMethods,
	}
}
