package vehicle

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"tesla-app/auth"
)

func (v *VehicleService) VehicleState(token auth.Token) (*VehicleStateResponse, error) {
	carId := os.Getenv("MY_CAR_ID")
	baseUrl := os.Getenv("TESLA_BASE_URL")

	url := fmt.Sprintf("%s/vehicles/%s", baseUrl, carId)

	client := &http.Client{}
	vehicleStateRequest, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	vehicleStateRequest.Header = http.Header{
		"Content-Type":  {"application/json"},
		"Authorization": {"Bearer " + token.AccessToken},
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

	if res.StatusCode != 200 {
		err = errors.New(fmt.Sprintf("vehicle state response failed with status code %s", res.Status))
		return nil, err
	}

	var responseBody TeslaVehicleApiResponse

	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return nil, err
	}

	vehicleState := &VehicleStateResponse{
		State: responseBody.Response.State,
		Vin:   responseBody.Response.Vin,
	}

	return vehicleState, nil
}
