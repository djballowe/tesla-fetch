package drawlogo

import (
	"fmt"
	"tesla-app/client/api"
)

func DrawStatus(vehicleData api.VehicleData) {
	carMap := map[string]string{
		"models": "Model S",
		"model3": "Model 3",
		"modelx": "Model X",
		"modely": "Model Y",
	}

	logo := []string{
		"\t⠀⠀⠀⠀⠀⠀⠀⠀⢀⣀⣀⣀⣀⣀⣀⡀⠀⠀⠀⠀⠀⠀⠀⠀\t",
		"\t⢀⣀⣤⣤⣶⡶⠿⠟⠛⠛⠛⠛⠛⠛⠛⠛⠻⠿⢶⣶⣤⣤⣀⡀\t",
		"\t⠘⢛⣫⣭⣴⣶⣾⣿⣿⣿⣦⠀⢀⣴⣿⣿⣿⣷⣶⣦⣭⣝⡛⠁\t",
		"\t⠀⠀⠙⠿⡿⠛⠉⠉⠙⣿⣿⣷⣾⣿⣿⠉⠉⠉⠛⢿⠿⠋⠀⠀\t",
		"\t⠀⠀⠀⠀⠀⠀⠀⠀⠀⢿⣿⣿⣿⣿⡟⠀⠀⠀⠀⠀⠀⠀⠀\t",
		"\t⠀⠀⠀⠀⠀⠀⠀⠀⠀⠸⣿⣿⣿⣿⠇⠀⠀⠀⠀⠀⠀⠀⠀\t",
		"\t⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣿⣿⣿⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀\t",
		"\t⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢿⣿⣿⡿⠀⠀⠀⠀⠀⠀⠀⠀⠀\t",
		"\t⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢸⣿⣿⠇⠀⠀⠀⠀⠀⠀⠀⠀⠀\t",
		"\t⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣿⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀\t",
		"\t⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢸⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀\t",
		"\t⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀\t",
	}

	header := fmt.Sprintf("\033[1;37mDavids %s\033[0m", carMap[vehicleData.CarType])
	headerLine := ""

	for i := 0; i < len(header); i++ {
		headerLine += "-"
	}

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
	}

	logoSize := len(logo)
	infoSize := len(info)

	centerLogo := logoSize / 2
	centerInfo := infoSize / 2
	infoIdx := 0

	for i, line := range logo {
		fmt.Printf("\033[1;31m%s\033[0m", line)
		if i >= centerLogo-centerInfo && infoIdx < infoSize {
			fmt.Printf(info[infoIdx])
			infoIdx++
		}
		fmt.Printf("\n")
	}
	fmt.Printf("\n")
}
