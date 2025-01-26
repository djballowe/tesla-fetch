package common

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	// "context"
	// "encoding/json"
	// "time"
	// "io"
	// "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	// "github.com/aws/aws-sdk-go-v2/config"
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

// func AuthCallBack(writer http.ResponseWriter, req *http.Request) {
// 	code := req.URL.Query().Get("code")
// 	state := req.URL.Query().Get("state")
// 	authStatus := false
// 	if code == "" || state == "" {
// 		http.Error(writer, "Invalid request in auth callback", http.StatusBadRequest)
// 		return
// 	}
//
// 	storeMutex.Lock()
// 	storedState := stateStore
// 	fmt.Println(storedState, "||", state)
// 	if storedState != state {
// 		http.Error(writer, "State does not match", http.StatusBadRequest)
// 		return
// 	}
// 	storeMutex.Unlock()
//
// 	tokens, err := callAuth(code)
// 	if err != nil || tokens == nil {
// 		http.Error(writer, "Failed to get auth token", http.StatusInternalServerError)
// 		return
// 	}
//
// 	storeMutex.Lock()
// 	token := Token{
// 		AccessToken:  tokens.AccessToken,
// 		RefreshToken: tokens.RefreshToken,
// 	}
// 	tokenStore[state] = token
// 	storeMutex.Unlock()
//
// 	fmt.Fprintf(writer, "Auth successful\n")
// 	authStatus = true
//
// 	notifyClientUrl := fmt.Sprintf("http://localhost:3000/notify?auth_status=%t", authStatus)
//
// 	_, err = http.Post(notifyClientUrl, "application/x-www-form-urlencoded", nil)
// 	if err != nil {
// 		http.Error(writer, fmt.Sprintf("Failed to update the client of auth status: %s", err.Error()), http.StatusInternalServerError)
// 		return
// 	}
// }
//
// func callAuth(code string) (*Token, error) {
// 	appUrl := os.Getenv("APP_BASE_URL")
// 	cfg, err := config.LoadDefaultConfig(context.TODO())
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	client := http.Client{}
// 	url := fmt.Sprintf("%s/auth?code=%s", appUrl, code)
//
// 	fmt.Println(url)
// 	authRequest, err := http.NewRequest("POST", url, nil)
// 	authRequest.Header.Set("Content-Type", "application/json")
//
// 	signer := v4.NewSigner()
// 	creds, err := cfg.Credentials.Retrieve(context.TODO())
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	// hash for an empty payload
// 	hash := "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
// 	err = signer.SignHTTP(context.TODO(), creds, authRequest, hash, "execute-api", cfg.Region, time.Now())
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	response, err := client.Do(authRequest)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer response.Body.Close()
//
// 	body, _ := io.ReadAll(response.Body)
// 	fmt.Println("Response:", string(body))
//
// 	var tokens Token
//
// 	err = json.Unmarshal(body, &tokens)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	return &tokens, nil
// }

func generateState() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
