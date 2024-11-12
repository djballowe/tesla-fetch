package vehicle

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"tesla-app/server/common"
)

type VehicleStateResponse struct {
	State      string
	Vin        string
	StatusCode int
}

type TeslaVehicleApiResponse struct {
	Response struct {
		ID             int64  `json:"id"`
		UserID         int64  `json:"user_id"`
		VehicleID      int64  `json:"vehicle_id"`
		Vin            string `json:"vin"`
		Color          any    `json:"color"`
		AccessType     string `json:"access_type"`
		GranularAccess struct {
			HidePrivate bool `json:"hide_private"`
		} `json:"granular_access"`
		Tokens                 any    `json:"tokens"`
		State                  string `json:"state"`
		InService              bool   `json:"in_service"`
		IDS                    string `json:"id_s"`
		CalendarEnabled        bool   `json:"calendar_enabled"`
		APIVersion             int    `json:"api_version"`
		BackseatToken          any    `json:"backseat_token"`
		BackseatTokenUpdatedAt any    `json:"backseat_token_updated_at"`
		BleAutopairEnrolled    bool   `json:"ble_autopair_enrolled"`
	} `json:"response"`
}

func VehicleState() (VehicleStateResponse, error) {
	var vehicleStateResponse = VehicleStateResponse{}

	tokenStore, state := common.GetTokenStore()
	carId := os.Getenv("MY_CAR_ID")
	baseUrl := os.Getenv("TESLA_BASE_URL")

	url := fmt.Sprintf("%s/vehicles/%s", baseUrl, carId)

	client := &http.Client{}
	vehicleStateRequest, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return returnGenericError(), err
	}

	vehicleStateRequest.Header = http.Header{
		"Content-Type":  {"application/json"},
		"Authorization": {"Bearer " + tokenStore[state].AccessToken},
	}

	res, err := client.Do(vehicleStateRequest)
	if err != nil {
		return returnGenericError(), err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return returnGenericError(), err
	}

	if res.StatusCode != 200 {
		err = errors.New(fmt.Sprintf("vehicle state response failed with status code %s", res.Status))
		return VehicleStateResponse{
			State:      "error",
			Vin:        "",
			StatusCode: res.StatusCode,
		}, err
	}

	var responseBody TeslaVehicleApiResponse

	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return returnGenericError(), err
	}

	vehicleStateResponse.State = responseBody.Response.State
	vehicleStateResponse.Vin = responseBody.Response.Vin
	vehicleStateResponse.StatusCode = res.StatusCode

	return vehicleStateResponse, nil
}

func returnGenericError() VehicleStateResponse {
	return VehicleStateResponse{
		State:      "error",
		Vin:        "",
		StatusCode: 500,
	}
}
