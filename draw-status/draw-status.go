package drawstatus

import (
	"fmt"
	"math"
	apitypes "tfetch/api/types"
)

const (
	red   = "\033[1;31m"
	white = "\033[1;37m"
	reset = "\033[0m"
)

func (d *DrawService) DrawStatus(vehicleData *apitypes.VehicleData) {
	logo := []string{
		"⠀⠀⠀⠀⠀⠀⠀⠀⢀⣀⣀⣀⣀⣀⣀⡀⠀⠀⠀⠀⠀⠀⠀⠀",
		"⢀⣀⣤⣤⣶⡶⠿⠟⠛⠛⠛⠛⠛⠛⠛⠛⠻⠿⢶⣶⣤⣤⣀⡀",
		"⠘⢛⣫⣭⣴⣶⣾⣿⣿⣿⣦⠀⢀⣴⣿⣿⣿⣷⣶⣦⣭⣝⡛⠁",
		"⠀⠀⠙⠿⡿⠛⠉⠉⠙⣿⣿⣷⣾⣿⣿⠉⠉⠉⠛⢿⠿⠋⠀⠀",
		"⠀⠀⠀⠀⠀⠀⠀⠀⠀⢿⣿⣿⣿⣿⡟⠀⠀⠀⠀⠀⠀⠀⠀ ",
		"⠀⠀⠀⠀⠀⠀⠀⠀⠀⠸⣿⣿⣿⣿⠇⠀⠀⠀⠀⠀⠀⠀⠀ ",
		"⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣿⣿⣿⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀ ",
		"⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢿⣿⣿⡿⠀⠀⠀⠀⠀⠀⠀⠀⠀ ",
		"⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢸⣿⣿⠇⠀⠀⠀⠀⠀⠀⠀⠀⠀ ",
		"⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣿⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀ ",
		"⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢸⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀ ",
		"⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀ ",
	}

	info := d.buildInfo(vehicleData)

	d.draw(logo, info)
}

func (d *DrawService) buildInfo(vehicleData *apitypes.VehicleData) []string {
	carMap := map[string]string{
		"models": "Model S",
		"model3": "Model 3",
		"modelx": "Model X",
		"modely": "Model Y",
	}

	header := fmt.Sprintf("%s%s, %s%s", white, vehicleData.VehicleName, carMap[vehicleData.CarType], reset)
	headerLine := ""

	for i := 0; i < len([]rune(header)); i++ {
		headerLine += "-"
	}

	info := []string{
		fmt.Sprintf("%s", header),
		fmt.Sprintf("%s", headerLine),
		fmt.Sprintf("%sColor%s: %s%s%s", red, reset, white, vehicleData.Color, reset),
		fmt.Sprintf("%sMiles%s: %s%d%s", red, reset, white, vehicleData.Odometer, reset),
		fmt.Sprintf("%sCharge%s: %s%d%%%s", red, reset, white, vehicleData.BatteryLevel, reset),
		fmt.Sprintf("%sCharge Limit%s: %s%d%s", red, reset, white, vehicleData.ChargeLimitSoc, reset),
		fmt.Sprintf("%sRange%s: %s%.1f Miles%s", red, reset, white, vehicleData.BatteryRange, reset),
		fmt.Sprintf("%sCharge State%s: %s%s%s", red, reset, white, vehicleData.ChargingState, reset),
		fmt.Sprintf("%sCharge Rate%s: %s%.1f%s", red, reset, white, vehicleData.ChargeRate, reset),
		fmt.Sprintf("%sClimate On%s: %s%t%s", red, reset, white, vehicleData.IsClimateOn, reset),
		fmt.Sprintf("%sClimate Inside%s: %s%d\u00B0F%s", red, reset, white, convertClimate(vehicleData.InsideTemp), reset),
		fmt.Sprintf("%sClimate Outside%s: %s%d\u00B0F%s", red, reset, white, convertClimate(vehicleData.OutsideTemp), reset),
		fmt.Sprintf("%sDriver Temp Setting%s: %s%d\u00B0F%s", red, reset, white, convertClimate(int(vehicleData.DriverTempSetting)), reset),
		fmt.Sprintf("%sPassenger Temp Setting%s: %s%d\u00B0F%s", red, reset, white, convertClimate(int(vehicleData.PassengerTempSetting)), reset),
		fmt.Sprintf("%sLocked%s: %s%t%s", red, reset, white, vehicleData.Locked, reset),
	}

	return info
}

func (d *DrawService) draw(logo []string, info []string) {
	logoSize := len(logo)
	infoSize := len(info)
	longest := findLongestLine(logo)

	infoIdx := 0
	loop := max(logoSize, infoSize)

	fmt.Printf("\n")
	for i := 0; i < loop; i++ {
		if i < len(logo) {
			fmt.Printf("\t%s%s%s", red, logo[i], reset)
			if infoIdx < infoSize {
				fmt.Printf("\t%s", info[infoIdx])
			}
		} else {
			fmt.Printf("%*s\t\t%s", longest, "", info[infoIdx])
		}
		infoIdx++
		fmt.Printf("\n")
	}

}

func findLongestLine(logo []string) int {
	longest := 0
	for _, line := range logo {
		longest = max(len([]rune(line)), longest)
	}
	return longest
}

func max(a int, b int) int {
	if a < b {
		return b
	}
	return a
}

func convertClimate(temp int) int {
	return int(math.Round(float64(temp)*1.8 + 32))
}
