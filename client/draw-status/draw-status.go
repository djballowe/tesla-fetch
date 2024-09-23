package drawlogo

import (
	"fmt"
	"tesla-app/client/api"
)

const REDBOLD = "\033[1;31m"
const ESCAPE = "\033[0m"

func DrawStatus(vehicleData api.VehicleData) {
	carMap := map[string]string{
		"models": "Model S",
		"model3": "Model 3",
		"modelx": "Model X",
		"modely": "Model Y",
	}

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

	header := fmt.Sprintf("\033[1;37m%s, %s\033[0m", vehicleData.VehicleName, carMap[vehicleData.CarType])
	headerLine := ""

	for idx := 0; idx < len(header); idx++ {
		headerLine += "-"
	}

	//	info := []string{
	//		fmt.Sprintf("%s", header),
	//		fmt.Sprintf("%s", headerLine),
	//		fmt.Sprintf("\033[1;31mName\033[0m: \033[1;37m%s\033[0m", test),
	//		fmt.Sprintf("\033[1;31mColor\033[0m: \033[1;37m%s\033[0m", test),
	//		fmt.Sprintf("\033[1;31mMiles\033[0m: \033[1;37m%s\033[0m", test),
	//		fmt.Sprintf("\033[1;31mCharge\033[0m: \033[1;37m%s\033[0m", test),
	//		fmt.Sprintf("\033[1;31mCharge State\033[0m: \033[1;37m%s\033[0m", test),
	//		fmt.Sprintf("\033[1;31mCharge Rate\033[0m: \033[1;37m%s\033[0m", test),
	//		fmt.Sprintf("\033[1;31mClimate On\033[0m: \033[1;37m%s\033[0m", test),
	//		fmt.Sprintf("\033[1;31mClimate\033[0m: \033[1;37m%s\033[0m", test),
	//		fmt.Sprintf("\033[1;31mLocked\033[0m: \033[1;37m%s\033[0m", test),
	//		fmt.Sprintf("\033[1;31mLocked\033[0m: \033[1;37m%s\033[0m", test),
	//		fmt.Sprintf("\033[1;31mLocked\033[0m: \033[1;37m%s\033[0m", test),
	//		fmt.Sprintf("\033[1;31mLocked\033[0m: \033[1;37m%s\033[0m", test),
	//		fmt.Sprintf("\033[1;31mLocked\033[0m: \033[1;37m%s\033[0m", test),
	//		fmt.Sprintf("\033[1;31mLocked\033[0m: \033[1;37m%s\033[0m", test),
	//		fmt.Sprintf("\033[1;31mLocked\033[0m: \033[1;37m%s\033[0m", test),
	//		fmt.Sprintf("\033[1;31mLocked\033[0m: \033[1;37m%s\033[0m", test),
	//		fmt.Sprintf("\033[1;31mLocked\033[0m: \033[1;37m%s\033[0m", test),
	//		fmt.Sprintf("\033[1;31mLocked\033[0m: \033[1;37m%s\033[0m", test),
	//	}

	info := []string{
		fmt.Sprintf("%s", header),
		fmt.Sprintf("%s", headerLine),
		fmt.Sprintf("\033[1;31mName\033[0m: \033[1;37m%s\033[0m", vehicleData.VehicleName),
		fmt.Sprintf("\t\033[1;31mColor\033[0m: \033[1;37m%s\033[0m", vehicleData.ExteriorColor),
		fmt.Sprintf("\t\033[1;31mMiles\033[0m: \033[1;37m%d\033[0m", vehicleData.Odometer),
		fmt.Sprintf("\t\033[1;31mCharge\033[0m: \033[1;37m%d\033[0m", vehicleData.BatteryLevel),
		fmt.Sprintf("\t\033[1;31mCharge State\033[0m: \033[1;37m%s\033[0m", vehicleData.CharginState),
		fmt.Sprintf("\t\033[1;31mCharge Rate\033[0m: \033[1;37m%d\033[0m", vehicleData.ChargeRate),
		fmt.Sprintf("\t\033[1;31mClimate On\033[0m: \033[1;37m%t\033[0m", vehicleData.IsClimateOn),
		fmt.Sprintf("\t\033[1;31mClimate\033[0m: \033[1;37m%d\033[0m", vehicleData.InsideTemp),
		fmt.Sprintf("\t\033[1;31mLocked\033[0m: \033[1;37m%t\033[0m", vehicleData.Locked),
		fmt.Sprintf("\t\033[1;31mLocked\033[0m: \033[1;37m%t\033[0m", vehicleData.Locked),
		fmt.Sprintf("\t\033[1;31mLocked\033[0m: \033[1;37m%t\033[0m", vehicleData.Locked),
		fmt.Sprintf("\t\033[1;31mLocked\033[0m: \033[1;37m%t\033[0m", vehicleData.Locked),
		fmt.Sprintf("\t\033[1;31mLocked\033[0m: \033[1;37m%t\033[0m", vehicleData.Locked),
		fmt.Sprintf("\t\033[1;31mLocked\033[0m: \033[1;37m%t\033[0m", vehicleData.Locked),
		fmt.Sprintf("\t\033[1;31mLocked\033[0m: \033[1;37m%t\033[0m", vehicleData.Locked),
		fmt.Sprintf("\t\033[1;31mLocked\033[0m: \033[1;37m%t\033[0m", vehicleData.Locked),
		fmt.Sprintf("\t\033[1;31mLocked\033[0m: \033[1;37m%t\033[0m", vehicleData.Locked),
		fmt.Sprintf("\t\033[1;31mLocked\033[0m: \033[1;37m%t\033[0m", vehicleData.Locked),
	}

	//info := buildInfo(vehicleData)

	draw(logo, info)

}

//func buildInfo(vehicleData api.VehicleData) []string {
//	infoMap := map[string]string {
//		"vehicle_name": "Name",
//		"exterior_color": "Color",
//		"odometer": "Miles",
//		"battery_level": "Battery Level",
//	}
//
//}

func draw(logo []string, info []string) {
	logoSize := len(logo)
	infoSize := len(info)
	longest := findLongestLine(logo)

	infoIdx := 0
	idx := 0
	loop := max(logoSize, infoSize)

	for idx < loop {
		if idx < len(logo) {
			fmt.Printf("\t\033[1;31m%s\033[0m", logo[idx])
			if infoIdx < infoSize {
				fmt.Printf("\t%s", info[infoIdx])
			}
		} else {
			fmt.Printf("%*s\t\t%s", longest, "", info[infoIdx])
		}
		infoIdx++
		idx++
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
