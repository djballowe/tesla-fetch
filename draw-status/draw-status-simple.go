package drawstatus

import (
	"encoding/json"
	"fmt"
	"tesla-app/data"
)

func (d *DrawService) DrawStatusSimple(vehicleData *data.VehicleData) error {
	jsonData, err := json.MarshalIndent(vehicleData, "", " ")
	if err != nil {
		return err
	}

	fmt.Println(string(jsonData))
	return nil
}
