package getdata

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"tesla-app/server/common"
)

type Command struct {
	AuthToken string
	Vin       string
	Command   string
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

	client := &http.Client{}
	req, err := http.NewRequest("GET", baseUrl+"/vehicles/{id}/vehicle_data", nil)
	if err != nil {
		http.Error(writer, "Failed to create get vehicles request", http.StatusInternalServerError)
	}

	req.Header = http.Header{
		"Content-Type":  {"application/json"},
		"Authorization": {"Bearer " + tokenStore[state].AccessToken},
	}

	res, err := client.Do(req)
	if err != nil {
		http.Error(writer, "Could not get vehicles", http.StatusInternalServerError)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
	}

	fmt.Println("Response Status: ", res.Status)
	fmt.Println("Response: ", string(body))

	var prettyJSON bytes.Buffer
	error := json.Indent(&prettyJSON, body, "", "\t")
	if error != nil {
		log.Println("JSON parse error: ", error)
		return
	}

	log.Println("Response", string(prettyJSON.Bytes()))

	//	command := common.Request{
	//		AuthToken: tokenStore[state].AccessToken,
	//		Vin:       "5YJ3E1EA8MF854070",
	//		Command:   "honk",
	//	}
	//
	//	common.HandleCommand(command)

}
