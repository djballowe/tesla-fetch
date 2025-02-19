package data

import (
	"os"
	"tesla-app/client/api"
	"tesla-app/client/auth"
	"tesla-app/client/ui"
)

func GetVehicleData(status chan ui.ProgressUpdate) (*api.VehicleData, error) {
	passphrase := os.Getenv("PASSPHRASE")
	store := auth.TokenStore{}
	token, err := store.LoadTokens(passphrase)
	if err != nil {
		if err.Error() == "No token stored" {
			token, err = auth.CallAuth()
			if err != nil {
				return nil, err
			}
		}
	}

	if store.IsExpired(token.CreateAt, token.ExpiresIn) {
		// go through refresh token flow
	}

	carDataResponse, err := api.CallGetVehicleData(*token, status)
	if err != nil {
		return nil, err
	}

	return carDataResponse, nil
}
