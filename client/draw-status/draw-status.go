package drawlogo

import (
	"fmt"
	"math"
	"tesla-app/client/api"
)

func DrawStatus(vehicleData api.VehicleData) {
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

	info := buildInfo(vehicleData)

	draw(logo, info)
}

func buildInfo(vehicleData api.VehicleData) []string {
	carMap := map[string]string{
		"models": "Model S",
		"model3": "Model 3",
		"modelx": "Model X",
		"modely": "Model Y",
	}

	header := fmt.Sprintf("\033[1;37m%s, %s\033[0m", vehicleData.VehicleName, carMap[vehicleData.CarType])
	headerLine := ""

	for i := 0; i < len([]rune(header)); i++ {
		headerLine += "-"
	}

	info := []string{
		fmt.Sprintf("%s", header),
		fmt.Sprintf("%s", headerLine),
		fmt.Sprintf("\033[1;31mColor\033[0m: \033[1;37m%s\033[0m", vehicleData.ExteriorColor),
		fmt.Sprintf("\033[1;31mMiles\033[0m: \033[1;37m%d\033[0m", vehicleData.Odometer),
		fmt.Sprintf("\033[1;31mCharge\033[0m: \033[1;37m%d\033[0m", vehicleData.BatteryLevel),
		fmt.Sprintf("\033[1;31mCharge State\033[0m: \033[1;37m%s\033[0m", vehicleData.ChargingState),
		fmt.Sprintf("\033[1;31mCharge Rate\033[0m: \033[1;37m%f\033[0m", vehicleData.ChargeRate),
		fmt.Sprintf("\033[1;31mClimate On\033[0m: \033[1;37m%t\033[0m", vehicleData.IsClimateOn),
		fmt.Sprintf("\033[1;31mClimate Inside\033[0m: \033[1;37m%d\033[0m", convertClimate(vehicleData.InsideTemp)),
		fmt.Sprintf("\033[1;31mClimate Outside\033[0m: \033[1;37m%d\033[0m", convertClimate(vehicleData.OutsideTemp)),
		fmt.Sprintf("\033[1;31mLocked\033[0m: \033[1;37m%t\033[0m", vehicleData.Locked),
	}

	return info
}

func draw(logo []string, info []string) {
	logoSize := len(logo)
	infoSize := len(info)
	longest := findLongestLine(logo)

	infoIdx := 0
	loop := max(logoSize, infoSize)

	fmt.Printf("\n")
	for i := 0; i < loop; i++ {
		if i < len(logo) {
			fmt.Printf("\t\033[1;31m%s\033[0m", logo[i])
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
