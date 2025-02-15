package auth

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net"
	"time"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
)

func loadEnvConfig() (*Config, error) {
	config := &Config{
		ClientId:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		Audience:     os.Getenv("AUDIENCE"),
		RedirectUri:  os.Getenv("REDIRECT_URI"),
		Scope:        os.Getenv("SCOPES"),
	}
	if config.ClientId == "" || config.ClientSecret == "" || config.Audience == "" || config.RedirectUri == "" || config.Scope == "" {
		return nil, fmt.Errorf("Missing environment variables")
	}
	return config, nil
}

func CallAuth() error {
	baseUrl, err := url.Parse("https://auth.tesla.com/oauth2/v3/authorize")
	if err != nil {
		log.Fatalf("Malformed auth url: %s", err)
		return err
	}

	config, err := loadEnvConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %s", err)
		return err
	}

	state := generateState()
	StoreMutex.Lock()
	StateStore = state
	StoreMutex.Unlock()

	authData := map[string]string{
		"response_type": "code",
		"client_id":     config.ClientId,
		"redirect_uri":  config.RedirectUri,
		"scope":         config.Scope,
		"state":         state,
	}
	params := url.Values{}

	for key, value := range authData {
		params.Add(key, value)
	}

	baseUrl.RawQuery = params.Encode()
	authUrl := baseUrl.String()
	tokens, err := startServer(authUrl)
	if err != nil {
		log.Fatalf("Could not start callback server: %s", err)
		return err
	}

	TokenStore[state] = *tokens

	log.Println("Tokens stored successful")

	return nil
}

var tokenChan = make(chan *Token)
var errChan = make(chan error)

func startServer(authUrl string) (*Token, error) {
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		return nil, err
	}
	defer listener.Close()

	server := &http.Server{
		Handler: http.HandlerFunc(callback),
	}

	go server.Serve(listener)
	log.Println("callback server started...")
	err = openBrowser(authUrl)
	if err != nil {
		return nil, err
	}

	select {
	case tokens := <-tokenChan:
		server.Shutdown(context.Background())
		return tokens, nil
	case err := <-errChan:
		server.Shutdown(context.Background())
		return nil, err
	case <-time.After(5 * time.Minute):
		server.Shutdown(context.Background())
		return nil, fmt.Errorf("authentication timed out")
	}
}

func callback(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/callback" {
		http.Error(w, "path not found", http.StatusNotFound)
		return
	}

	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")
	log.Println(state)

	tokens, err := exchangeCodeForToken(code)
	if err != nil {
		log.Fatalf("error exchanging code for token: %s", err)
		errChan <- err
		return
	}

	w.Write([]byte("Authentication successful! You can close this window."))
	tokenChan <- tokens
	return
}

// helpers

func exchangeCodeForToken(code string) (*Token, error) {
	log.Println("Exchanging code for token...")

	baseUrl, err := url.Parse("https://auth.tesla.com/oauth2/v3/token")
	if err != nil {
		return nil, err
	}
	tokenUrl := baseUrl.String()

	config, err := loadEnvConfig()
	if err != nil {
		log.Fatal("Failed to load config: ", err)
	}

	exchangeData := map[string]string{
		"grant_type":    "authorization_code",
		"client_id":     config.ClientId,
		"client_secret": config.ClientSecret,
		"audience":      config.Audience,
		"redirect_uri":  config.RedirectUri,
		"scope":         config.Scope,
		"code":          code,
	}

	payload, err := json.Marshal(exchangeData)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(tokenUrl, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var tokenResponse Token
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return nil, err
	}

	return &tokenResponse, nil
}

func generateState() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func openBrowser(url string) error {
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

	return nil
}
