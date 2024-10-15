package vehicle

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"tesla-app/server/common"
	"time"
)

type WakeResponse struct {
	State string
}

type TeslaVehicleWakeResponse struct {
	State string `json:"state"`
}

func Wake() (WakeResponse, error) {
	var wakeResponse = WakeResponse{}

	tokenStore, state := common.GetTokenStore()
	carId := os.Getenv("MY_CAR_ID")
	baseUrl := os.Getenv("TESLA_BASE_URL")

	url := fmt.Sprintf("%s/vehicles/%s/wake_up", baseUrl, carId)

	client := &http.Client{}
	vehicleStateRequest, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return wakeResponse, err
	}

	vehicleStateRequest.Header = http.Header{
		"Content-Type":  {"application/json"},
		"Authorization": {"Bearer " + tokenStore[state].AccessToken},
	}

	res, err := client.Do(vehicleStateRequest)
	if err != nil {
		return wakeResponse, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	fmt.Println(res.Status)
	if err != nil {
		return wakeResponse, err
	}

	var responseBody TeslaVehicleWakeResponse

	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return wakeResponse, err
	}

	wakeResponse.State = responseBody.State

	return wakeResponse, nil
}

func PollWake() error {
	state := "offline"
	timeout := time.After(30 * time.Second)

	for {
		select {
			case <- timeout:
			return errors.New("Timeout, could not wake vehicle")

		default:
			wakeResponse, _ := Wake()
			state = wakeResponse.State
			if state == "online" {
				return nil
			}
		}
	}
}
