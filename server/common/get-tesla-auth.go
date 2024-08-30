package common

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

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Config struct {
	ClientId     string
	ClientSecret string
	Audience     string
	RedirectUri  string
	Scope        string
}

var (
	stateStore string
	tokenStore = make(map[string]Token)
	storeMutex sync.Mutex
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

func AuthCallBack(writer http.ResponseWriter, req *http.Request) {
	code := req.URL.Query().Get("code")
	state := req.URL.Query().Get("state")
	if code == "" || state == "" {
		http.Error(writer, "Invalid request in auth callback", http.StatusBadRequest)
		return
	}

	storeMutex.Lock()
	storedState := stateStore
	fmt.Println(storedState, "||", state)
	if storedState != state {
		http.Error(writer, "State does not match", http.StatusBadRequest)
		return
	}
	storeMutex.Unlock()

	tokens, err := exchangeCodeForToken(code)
	if err != nil {
		http.Error(writer, "Failed to get auth token", http.StatusInternalServerError)
		return
	}

	storeMutex.Lock()
	token := Token{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}
	tokenStore[state] = token
	storeMutex.Unlock()

	fmt.Fprintf(writer, "Auth successful token stored")
}

func exchangeCodeForToken(code string) (*Token, error) {
	fmt.Println("Exchanging code for token...")

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

func GetTokenStore() (map[string]Token, string) {
	fmt.Println("Getting Token Store")
	storeMutex.Lock()
	defer storeMutex.Unlock()
	copyStore := make(map[string]Token)
	stateCopy := stateStore

	for k, v := range tokenStore {
		copyStore[k] = v
	}

	return copyStore, stateCopy
}
