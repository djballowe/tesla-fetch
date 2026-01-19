package tesla

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"tfetch/model"
)

type VehicleService struct{}

func (v *VehicleService) VehicleState(token model.Token) (*model.VehicleStateResponse, error) {
	carId := os.Getenv("MY_CAR_ID")
	baseUrl := os.Getenv("TESLA_BASE_URL")

	url := fmt.Sprintf("%s/vehicles/%s", baseUrl, carId)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
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

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("vehicle state response failed with status code %s", res.Status)
	}

	var responseBody model.TeslaVehicleApiResponse

	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return nil, err
	}

	vehicleState := &model.VehicleStateResponse{
		State: responseBody.Response.State,
		Vin:   responseBody.Response.Vin,
	}

	return vehicleState, nil
}
