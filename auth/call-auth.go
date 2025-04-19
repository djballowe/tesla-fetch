package auth

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"time"
)

func loadEnvConfig() (*Config, error) {
	config := &Config{
		ClientId:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		Audience:     os.Getenv("AUDIENCE"),
		RedirectUri:  os.Getenv("REDIRECT_URI"),
		Scope:        os.Getenv("SCOPES"),
		Passphrase:   os.Getenv("PASSPHRASE"),
	}
	if config.ClientId == "" || config.ClientSecret == "" || config.Audience == "" || config.RedirectUri == "" || config.Scope == "" {
		return nil, fmt.Errorf("Missing environment variables")
	}
	return config, nil
}

func (a *AuthService) CallAuth() (*Token, error) {
	baseUrl, err := url.Parse("https://auth.tesla.com/oauth2/v3/authorize")
	if err != nil {
		log.Fatalf("Malformed auth url: %s", err)
		return nil, err
	}

	config, err := loadEnvConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %s", err)
		return nil, err
	}

	state := generateState()
	StateStore = state

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
		return nil, err
	}

	store, err := a.NewTokenStore(config.Passphrase)
	if err != nil {
		return nil, err
	}

	err = store.SaveTokens(tokens, store.salt)
	if err != nil {
		return nil, err
	}

	tokenStore, err := store.LoadTokens(config.Passphrase)
	if err != nil {
		return nil, err
	}

	return tokenStore, nil
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

	if state != StateStore || state == "" {
		http.Error(w, "Internal auth error", http.StatusBadRequest)
		errChan <- fmt.Errorf("Internal auth error state missing or mismatch")
	}

	if code == "" {
		http.Error(w, "Internal auth error", http.StatusBadRequest)
		errChan <- fmt.Errorf("Internal auth error missing code")
	}

	tokens, err := exchangeCodeForToken(code)
	if err != nil {
		http.Error(w, "Internal auth error", http.StatusBadRequest)
		errChan <- err
		return
	}

	w.Write([]byte("Authentication successful! You can close this window."))
	tokenChan <- tokens
	return
}

// helpers

func exchangeCodeForToken(code string) (*Token, error) {
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
	tokenResponse.CreateAt = time.Now()

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
