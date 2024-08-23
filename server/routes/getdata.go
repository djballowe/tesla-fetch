package getdata

import (
	//	"bytes"
	"encoding/json"
	"fmt"
	"io"
	//	"log"
	"net/http"
	"os"
	"tesla-app/server/common"
)

type Command struct {
	AuthToken string
	Vin       string
	Command   string
}

type Return struct {
	State               string `json:"state"`
	BatteryLevel        int    `json:"battery_level"`
	ChargeRate          int    `json:"charge_rate"`
	ChargingState       string `json:"chargin_state"`
	MinutesToFullCharge int    `json:"minutes_to_full_charge"`
	InsideTemp          int    `json:"inside_temp"`
	IsClimateOn         bool   `json:"is_climate_on"`
	IsPreconditioning   bool   `json:"is_preconditioning"`
	OutsideTemp         int    `json:"outside_temp"`
	Locked              bool   `json:"locked"`
	Odometer            int    `json:"odometer"`
	Color               string `json:"exterior_color"`
	VehicleName         string `json:"vehicle_name"`
}

func GetCarStatus(writer http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/car" {
		http.Error(writer, "404 not found", http.StatusNotFound)
	}

	if req.Method != "GET" {
		http.Error(writer, "This method is not supported", http.StatusNotFound)
	}
	tokenStore, state := common.GetTokenStore()
	baseUrl := os.Getenv("TESLA_BASE_URL")
	carId := os.Getenv("MY_CAR_ID")

	url := fmt.Sprintf("%s/vehicles/%s/vehicle_data", baseUrl, carId)

	fmt.Println(url)

	client := &http.Client{}
	getReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		http.Error(writer, "Failed to create get vehicles request", http.StatusInternalServerError)
	}

	getReq.Header = http.Header{
		"Content-Type":  {"application/json"},
		"Authorization": {"Bearer " + tokenStore[state].AccessToken},
	}

	res, err := client.Do(getReq)
	if err != nil {
		http.Error(writer, "Could not get vehicles", http.StatusInternalServerError)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
	}

	fmt.Println("Response Status: ", res.Status)
	//	fmt.Println("Response: ", string(body))

	//var prettyJSON bytes.Buffer
	//error := json.Indent(&prettyJSON, body, "", "\t")
	//if error != nil {
	//	log.Println("JSON parse error: ", error)
	//	return
	//}

	//log.Println("Response", string(prettyJSON.Bytes()))

	var returnBody map[string]interface{}

	err = json.Unmarshal(body, &returnBody)
	if err != nil {
		http.Error(writer, "Could not unmarshal response body", http.StatusInternalServerError)
	}
	// this is erroring fix

	response := returnBody["response"].(map[string]interface{})
	test := response["state"].(string)

	fmt.Println(test)

}
