package api

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
)

type AuthResponse struct {
	CallbackUrl string `json:"callback_url"`
}

func CallAuth() error {
	baseUrl := os.Getenv("BASE_URL")
	authUrl := fmt.Sprintf("%s/auth", baseUrl)
	key := os.Getenv("API_KEY")

	// go back to localhost auth but close it as soon as you're done

	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/auth", nil)
	if err != nil {
		return err
	}

	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	var authResponse AuthResponse

	if err := json.NewDecoder(resp.Body).Decode(&authResponse); err != nil {
		return err
	}

	fmt.Println(authResponse.CallbackUrl)

	openBrowser(authResponse.CallbackUrl)

	// callbackResp, err := http.Get(authResponse.CallbackUrl)
	// if err != nil {
	// 	return err
	// }
	//
	// body, err := io.ReadAll(callbackResp.Body)
	// if err != nil {
	// 	return err
	// }

	// fmt.Println(string(body))

	return nil
}

func generateState() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func buildNotificationServer(notify chan bool) {
	http.HandleFunc("/notify", func(writer http.ResponseWriter, resp *http.Request) {
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
