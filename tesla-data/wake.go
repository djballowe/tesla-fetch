package tesla

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"tfetch/model"
	"time"
)

type WakeService struct{}

func (w *WakeService) Wake(token model.Token) (*model.WakeResponse, error) {
	carId := os.Getenv("MY_CAR_ID")
	baseUrl := os.Getenv("TESLA_BASE_URL")

	url := fmt.Sprintf("%s/vehicles/%s/wake_up", baseUrl, carId)

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header = http.Header{
		"Content-Type":  {"application/json"},
		"Authorization": {"Bearer " + token.AccessToken},
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var responseBody model.TeslaVehicleWakeResponse

	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return nil, err
	}

	wakeResponse := &model.WakeResponse{
		State: responseBody.Response.State,
	}

	return wakeResponse, nil
}

func (w *WakeService) PollWake(token model.Token, status model.StatusLoggerMethods) error {
	status.Log("Waking Vehicle")
	state := "offline"
	timeout := time.After(30 * time.Second)
	ticker := time.NewTicker(5 * time.Second)

	for {
		select {
		case <-timeout:
			return errors.New("Timeout, could not wake vehicle")

		case <-ticker.C:
			wakeResponse, err := w.Wake(token)
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
