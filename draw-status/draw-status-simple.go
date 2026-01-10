package drawstatus

import (
	"fmt"
	"strconv"
	"tfetch/data"
)

func (d *DrawService) DrawStatusSimple(vehicleData *data.VehicleData) error {
	carMap := map[string]string{
		"models": "Model S",
		"model3": "Model 3",
		"modelx": "Model X",
		"modely": "Model Y",
	}

	info := []string{
		"Nickname: " + vehicleData.VehicleName,
		"Model: " + carMap[vehicleData.CarType],
		"Color: " + vehicleData.Color,
		"Miles: " + strconv.Itoa(vehicleData.Odometer),
		"Charge: " + strconv.Itoa(vehicleData.BatteryLevel) + "%",
		"Charge Limit: " + strconv.Itoa(vehicleData.ChargeLimitSoc),
		"Range: " + strconv.FormatFloat(vehicleData.BatteryRange, 'f', 1, 64) + " Miles",
		"Charge State: " + vehicleData.ChargingState,
		"Charge Rate: " + strconv.FormatFloat(vehicleData.ChargeRate, 'f', 1, 64),
		"Climate On: " + strconv.FormatBool(vehicleData.IsClimateOn),
		"Climate Inside: " + strconv.Itoa(convertClimate(vehicleData.InsideTemp)) + "°F",
		"Climate Outside: " + strconv.Itoa(convertClimate(vehicleData.OutsideTemp)) + "°F",
		"Driver Temp Setting: " + strconv.Itoa(convertClimate(int(vehicleData.DriverTempSetting))) + "°F",
		"Passenger Temp Setting: " + strconv.Itoa(convertClimate(int(vehicleData.PassengerTempSetting))) + "°F",
		"Locked: " + strconv.FormatBool(vehicleData.Locked),
	}

	for _, data := range info {
		fmt.Println(data)
	}

	return nil
}
