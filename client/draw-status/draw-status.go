package drawlogo

import (
	"fmt"
	"tesla-app/client/api"
)

type VehicleData struct {
	State               string `json:"state"`
	BatteryLevel        int    `json:"battery_level"`
	ChargeRate          int    `json:"charge_rate"`
	CharginState        string `json:"chargin_state"`
	MinutesToFullCharge int    `json:"minutes_to_full_charge"`
	TimeToFullCharge    int    `json:"time_to_full_charge"`
	InsideTemp          int    `json:"inside_temp"`
	IsClimateOn         bool   `json:"is_climate_on"`
	IsPreconditioning   bool   `json:"is_preconditioning"`
	OutsideTemp         int    `json:"outside_temp"`
	Locked              bool   `json:"locked"`
	Odometer            int    `json:"odometer"`
	ExteriorColor       string `json:"exterior_color"`
	VehicleName         string `json:"vehicle_name"`
}

func DrawStatus(vehicleData api.VehicleData) {

	fmt.Println(vehicleData.VehicleName)

	logo := []string{
		"	                ⣀⣀⣀⣀⣠⣤⣤⣤⣤⣤⣤⣤⣀⣀⣀⣀					",
		"	      ⢀⣀⣤⣤⣴⣶⣶⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣶⣶⣦⣤⣤⣀⡀			",
		"	 ⣀⣤⣴⣶⣿⣿⣿⣿⣿⠿⠟⠛⠛⠛⠋⠉⣉⣉⠀⠀⠀⠀⠀⠀⠀⠀⠀⢈⣉⣉⠉⠙⠛⠛⠛⠿⠿⣿⣿⣿⣿⣿⣶⣤⣤⣀	",
		"	⢻⣿⣿⠿⠛⠛⢉⣁⣤⣤⣴⣶⣶⣿⣿⣿⣿⣿⣷⣄⠀⠀⠀⠀⠀⠀⣠⣿⣿⣿⣿⣿⣿⣶⣶⣤⣤⣤⣈⡉⠛⠛⠿⣿⣿⠏ 	",
		"	 ⠁⣠⣤⣶⣾⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣧⡀⠀⠀⢠⣾⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣷⣶⣤⣄⠈⠀		",
		"	  ⠙⠻⣿⣿⣿⣿⣿⣿⠿⠛⠛⠛⠛⠛⢻⣿⣿⣿⣿⣿⣦⣴⣿⣿⣿⣿⣿⡟⠛⠛⠛⠛⠛⠿⣿⣿⣿⣿⣿⣿⠟⠃⠀		",
		"	    ⠈⠛⠿⣿⠏⠀⠀⠀⠀⠀⠀⠀⠘⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠀⠀⠀⠀⠀⠀⠀⠀⢹⣿⠿⠛⠁⠀⠀⠀		",
		"	                 ⢿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡏⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀		",
		"	                 ⢸⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠃⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀		",
		"	                  ⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀		",
		"	                  ⢿⣿⣿⣿⣿⣿⣿⣿⣿⡟⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀		",
		"	                  ⢸⣿⣿⣿⣿⣿⣿⣿⣿⠇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀		",
		"	                   ⣿⣿⣿⣿⣿⣿⣿⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀		",
		"	                   ⢹⣿⣿⣿⣿⣿⣿⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀		",
		"	                   ⠘⣿⣿⣿⣿⣿⣿⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀		",
		"	                    ⣿⣿⣿⣿⣿⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀		",
		"	                    ⢹⣿⣿⣿⣿⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀		",
		"	                    ⠘⣿⣿⣿⣿⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀		",
		"	                     ⢿⣿⣿⡟⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀		",
		"	                     ⢸⣿⣿⠇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀		",
		"	                      ⣿⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀		",
		"	                      ⢿⡟⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀		",
		"	                      ⢸⠇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀		",
	}

	info := []string{
		fmt.Sprintf("Name:		%s", vehicleData.VehicleName),
		fmt.Sprintf("Color:	%s", vehicleData.ExteriorColor),
		fmt.Sprintf("Miles:		%d", vehicleData.Odometer),
		fmt.Sprintf("Charge:		%d", vehicleData.BatteryLevel),
		fmt.Sprintf("Charge State:	%s", vehicleData.CharginState),
		fmt.Sprintf("Charge Rate:	%d", vehicleData.ChargeRate),
		fmt.Sprintf("Climate On:	%t", vehicleData.IsClimateOn),
		fmt.Sprintf("Climate:		%d", vehicleData.InsideTemp),
		fmt.Sprintf("Locked:		%t", vehicleData.Locked),
	}

	logoSize := len(logo)
	infoSize := len(info)

	centerLogo := logoSize / 2
	centerInfo := infoSize / 2
	infoIdx := 0

	for i, line := range logo {
		fmt.Printf(line)
		if i >= centerLogo-centerInfo && infoIdx < infoSize {
			fmt.Printf(info[infoIdx])
			infoIdx++
		}
		fmt.Printf("\n")
	}
}
