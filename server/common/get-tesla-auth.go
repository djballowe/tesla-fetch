package common

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
)

type Config struct {
	ClientId     string
	ClientSecret string
	Audience     string
	RedirectUri  string
	Scope        string
}

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

func GetTeslaAuth(writer http.ResponseWriter, req *http.Request) {
	baseUrl, err := url.Parse("https://auth.tesla.com/oauth2/v3/authorize")
	if err != nil {
		log.Fatal("Malformed auth url", err)
	}

	config, err := loadEnvConfig()
	if err != nil {
		log.Fatal("Failed to load config: ", err)
	}

	state := generateState()
	fmt.Println(state)
	storeMutex.Lock()
	stateStore = state
	storeMutex.Unlock()

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

	fmt.Println("State stored redirecting...")

	writer.Header().Set("Location", authUrl)
	writer.WriteHeader(http.StatusFound)
}

func generateState() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
