package api

import (
	"fmt"
	"net/http"
	"os/exec"
	"runtime"
)

type AuthResponse struct {
	Body       string
	StatusCode int
}

func CallAuth() (AuthResponse, error) {
	authResponse := &AuthResponse{}

	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/auth", nil)
	if err != nil {
		return *authResponse, err
	}

	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return *authResponse, err
	}

	redirectUrl := resp.Header.Get("Location")
	openBrowser(redirectUrl)
	fmt.Println("redirect url: ", redirectUrl)

	return *authResponse, nil
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
