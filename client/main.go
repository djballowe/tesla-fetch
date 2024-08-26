package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const BASE_URL = "localhost:8080"

type DataResponse struct {
	Response string
}

type Routes struct {
	Data string
	Honk string
	Auth string
}

var routes = Routes{
	Data: "/data",
	Honk: "/honk",
	Auth: "/auth",
}

func main() {
	req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/data", nil)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}

	req.Header.Add("Accept", "aplication/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if resp == nil || err != nil {
		fmt.Println(err.Error())
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
		return
	}

	responseBody := &DataResponse{}

	json.Unmarshal(body, responseBody)

	fmt.Println(responseBody)
	return
}
