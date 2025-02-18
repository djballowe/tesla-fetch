package vehicle

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"tesla-app/client/auth"
	"tesla-app/client/ui"
	"time"
)

func Wake() (*WakeResponse, error) {
	tokenStore, state := auth.GetTokenStore()
	carId := os.Getenv("MY_CAR_ID")
	baseUrl := os.Getenv("TESLA_BASE_URL")

	url := fmt.Sprintf("%s/vehicles/%s/wake_up", baseUrl, carId)

	client := &http.Client{}
	vehicleStateRequest, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	vehicleStateRequest.Header = http.Header{
		"Content-Type":  {"application/json"},
		"Authorization": {"Bearer " + tokenStore[state].AccessToken},
	}

	res, err := client.Do(vehicleStateRequest)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var responseBody TeslaVehicleWakeResponse

	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return nil, err
	}

	wakeResponse := &WakeResponse{
		State: responseBody.Response.State,
	}

	return wakeResponse, nil
}

func PollWake(status chan ui.ProgressUpdate) error {
	status <- ui.ProgressUpdate{Message: "Waking vehicle"}
	state := "offline"
	timeout := time.After(30 * time.Second)
	ticker := time.NewTicker(5 * time.Second)

	for {
		select {
		case <-timeout:
			return errors.New("Timeout, could not wake vehicle")

		case <-ticker.C:
			wakeResponse, err := Wake()
			if err != nil {
				return err
			}
			state = wakeResponse.State
			if state == "online" {
				return nil
			}
		}
	}
}
