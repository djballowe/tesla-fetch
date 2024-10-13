package vehicle 

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"tesla-app/server/common"
)

type VehicleStateResponse struct {
	Error error
	State string
}

func VehicleState() VehicleStateResponse {
	var vehicleStateResponse = VehicleStateResponse{}

	tokenStore, state := common.GetTokenStore()
	baseUrl := os.Getenv("TESLA_BASE_URL")
	carId := os.Getenv("MY_CAR_ID")

	url := fmt.Sprintf("%s/vehicles/%s/wake_up", baseUrl, carId)

	client := &http.Client{}
	vehicleStateRequest, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return vehicleStateResponse
	}

	vehicleStateRequest.Header = http.Header{
		"Content-Type":  {"application/json"},
		"Authorization": {"Bearer " + tokenStore[state].AccessToken},
	}

	res, err := client.Do(vehicleStateRequest)
	if err != nil {
		return vehicleStateResponse
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	fmt.Println(res.Status)
	if err != nil {
		return vehicleStateResponse
	}

	fmt.Println(string(body))

	return vehicleStateResponse
}
