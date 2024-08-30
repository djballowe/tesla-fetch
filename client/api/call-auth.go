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

	notify := make(chan bool, 1)

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
	// app needs to somehow wait for the user to hit the callback route
	openBrowser(redirectUrl)
	buildNotificationServer(notify)
	<-notify

	return *authResponse, nil
}

func buildNotificationServer(notify chan bool) {
	http.HandleFunc("/notify-client", func(writer http.ResponseWriter, resp *http.Request) {
		notify <- true
	})
	go func() {
		err := http.ListenAndServe(":3000", nil)
		if err != nil {
			fmt.Println("Error starting notify server")
		}
	}()

	return
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
