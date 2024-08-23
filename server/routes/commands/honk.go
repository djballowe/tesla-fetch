package commands 

import (
	"fmt"
	"io"
	"net/http"
	"os"
	getTeslaAuth "tesla-app/server/common"
)

func Honk(writer http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		http.Error(writer, "404 not found", http.StatusNotFound)
	}

	tokenStore, state := getTeslaAuth.GetTokenStore()
	baseUrl := os.Getenv("TESLA_BASE_URL")

	client := &http.Client{}
	req, err := http.NewRequest("POST", baseUrl+"/vehicles/1493089035103228/command/honk_horn", nil)
	if err != nil {
		http.Error(writer, "Failed to create honk horn command", http.StatusInternalServerError)
	}

	req.Header = http.Header{
		"Content-Type":  {"application/json"},
		"Authorization": {"Bearer " + tokenStore[state].AccessToken},
	}

	res, err := client.Do(req)
	if err != nil {
		http.Error(writer, "Could not issue honk horn command", http.StatusInternalServerError)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
	}
	fmt.Println("Response Status: ", res.Status)
	fmt.Println("Response: ", string(body))
}

