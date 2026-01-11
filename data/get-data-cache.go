package data

import (
	"encoding/json"
	"os"
	"path/filepath"
	"tfetch/vehicle-state"
)

func getDataCache(vehicleService vehicle.VehicleMethods, vehicleData *VehicleData) error {
	err := getState(vehicleData)
	if err != nil {
		return err
	}
	return nil
}

func getStateFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	localDir := filepath.Join(homeDir, ".local", "share", "tfetch")

	err = os.MkdirAll(localDir, 0700)
	if err != nil {
		return "", err
	}

	statePath := filepath.Join(localDir, "vehicle-state.json")

	return statePath, nil
}

func getState(vehicleData *VehicleData) error {
	statePath, err := getStateFilePath()
	if err != nil {
		return err
	}

	data, err := os.ReadFile(statePath)
	if err != nil {
		if os.IsNotExist(err) {
			return ErrNoStateFile
		}
		return err
	}

	err = json.Unmarshal(data, vehicleData)
	if err != nil {
		return err
	}

	return nil
}
