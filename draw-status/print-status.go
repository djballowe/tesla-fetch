package drawstatus

import (
	"fmt"
	"tesla-app/data"
)

func printJson(vehicleData *data.VehicleData) {
	fmt.Sprintf("Charging State: %s", vehicleData.ChargingState)
	fmt.Sprintf("Charging BatteryLevel: %d", vehicleData.BatteryLevel)
}
