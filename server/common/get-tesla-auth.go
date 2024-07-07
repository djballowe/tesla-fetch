package getTeslaAuth

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"
)

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

var (
	stateStore = make(map[string]string)
	tokenStore = make(map[string]TokenResponse)
	storeMutex sync.Mutex
)

var (
	clientId     = os.Getenv("CLIENT_ID")
	clientSecret = os.Getenv("CLIENT_SECRTET")
	audience     = os.Getenv("AUDIENCE")
	redirectUri  = os.Getenv("REDIRECT_URI")
	scope        = os.Getenv("SCOPES")
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
	storeMutex.Lock()
	stateStore[state] = state
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

	http.Redirect(writer, req, authUrl, http.StatusFound)
}

func AuthCallBack(writer http.ResponseWriter, req *http.Request) {
	code := req.URL.Query().Get("code")
	state := req.URL.Query().Get("state")
	if code == "" || state == "" {
		http.Error(writer, "Invalid request in auth callback", http.StatusBadRequest)
		return
	}

	storeMutex.Lock()
	storedState, exists := stateStore[state]
	if !exists || storedState != state {
		http.Error(writer, "State does not match", http.StatusBadRequest)
		return
	}

	token, err := exchangeCodeForToken(code)
	if err != nil {
		http.Error(writer, "Failed to get auth token", http.StatusInternalServerError)
		return
	}

	storeMutex.Lock()
	tokenStore[state] = *token
	storeMutex.Unlock()
}

func exchangeCodeForToken(code string) (*TokenResponse, error) {
	baseUrl, err := url.Parse("https://auth.tesla.com/oauth2/v3/token")
	if err != nil {
		return nil, err
	}
	tokenUrl := baseUrl.String()

	config, err := loadEnvConfig()
	if err != nil {
		return nil, err
	}

	exchangeData := map[string]string{
		"grant_type":    "authorization_code",
		"client_id":     config.ClientId,
		"client_secret": config.ClientSecret,
		"code":          code,
		"audience":      config.Audience,
		"redirect_uri":  config.RedirectUri,
		"scope":         config.Scope,
	}

	jsonPayload, err := json.Marshal(exchangeData)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(tokenUrl, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var tokenResponse TokenResponse
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

//func GetTokenStore() (*sync.Mutex, map[string]TokenResponse, string) {
//	return &storeMutex, tokenStore, stateStore[state]
//}
