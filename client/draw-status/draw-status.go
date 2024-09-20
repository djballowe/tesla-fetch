package drawlogo

import (
	"fmt"
	"tesla-app/client/api"
)

func DrawStatus(vehicleData api.VehicleData) {

	logo := []string{
		"\033[38;5;88m",
		"	████████╗███████╗███████╗██╗      █████╗",
		"	╚══██╔══╝██╔════╝██╔════╝██║     ██╔══██╗",
		"	   ██║   █████╗  ███████╗██║     ███████║",
		"	   ██║   ██╔══╝  ╚════██║██║     ██╔══██║",
		"	   ██║   ███████╗███████║███████╗██║  ██║",
		"	   ╚═╝   ╚══════╝╚══════╝╚══════╝╚═╝  ╚═╝",
		"\033[0m",
	}

	info := []string{
		fmt.Sprintf("\t\033[1;31mName\033[0m: \033[1;37m%s\033[0m", vehicleData.VehicleName),
		fmt.Sprintf("\t\033[1;31mColor\033[0m: \033[1;37m%s\033[0m", vehicleData.ExteriorColor),
		fmt.Sprintf("\t\033[1;31mMiles\033[0m: \033[1;37m%d\033[0m", vehicleData.Odometer),
		fmt.Sprintf("\t\033[1;31mCharge\033[0m: \033[1;37m%d\033[0m", vehicleData.BatteryLevel),
		fmt.Sprintf("\t\033[1;31mCharge State\033[0m: \033[1;37m%s\033[0m", vehicleData.CharginState),
		fmt.Sprintf("\t\033[1;31mCharge Rate\033[0m: \033[1;37m%d\033[0m", vehicleData.ChargeRate),
		fmt.Sprintf("\t\033[1;31mClimate On\033[0m: \033[1;37m%t\033[0m", vehicleData.IsClimateOn),
		fmt.Sprintf("\t\033[1;31mClimate\033[0m: \033[1;37m%d\033[0m", vehicleData.InsideTemp),
		fmt.Sprintf("\t\033[1;31mLocked\033[0m: \033[1;37m%t\033[0m", vehicleData.Locked),
	}

	for _, line := range logo {
		fmt.Printf(line)
		fmt.Printf("\n")
	}

	for _, info := range info {
		fmt.Printf(info)
		fmt.Printf("\n")
	}
}
