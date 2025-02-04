package api

import (
	"encoding/json"
	"fmt"
	//	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
)

type AuthResponse struct {
	StatusCode int
}

type CallbackResponse struct {
	CallbackUrl string `json:"callback_url"`
}

func CallAuth() (AuthResponse, error) {
	authResponse := &AuthResponse{}

	baseUrl := os.Getenv("BASE_URL")
	authUrl := fmt.Sprintf("%s/auth", baseUrl)
	key := os.Getenv("API_KEY")

	req, err := http.NewRequest("GET", authUrl, nil)
	if err != nil {
		return *authResponse, err
	}
	req.Header.Add("x-api-key", key)

	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return *authResponse, err
	}

	var callbackResponse CallbackResponse

	if err := json.NewDecoder(resp.Body).Decode(&callbackResponse); err != nil {
		return *authResponse, err
	}

	fmt.Println(callbackResponse.CallbackUrl)

	return AuthResponse{
		StatusCode: 200,
	}, nil
}

func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	// windows can eat sand
	default:
		err = fmt.Errorf("unsupported platform")
	}

	if err != nil {
		fmt.Println("Failed to open browser:", err)
	}
}
