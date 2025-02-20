package auth

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"time"
)

func RefreshToken(refreshToken string) (*Token, error) {
	config, err := loadEnvConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %s", err)
		return nil, err
	}

	baseUrl, err := url.Parse("https://auth.tesla.com/oauth2/v3/token")
	if err != nil {
		return nil, err
	}
	tokenUrl := baseUrl.String()

	exchangeData := map[string]string{
		"grant_type":    "refresh_token",
		"client_id":     config.ClientId,
		"refresh_token": refreshToken,
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

	store, err := NewTokeStore(config.Passphrase)
	if err != nil {
		return nil, err
	}

	err = store.SaveTokens(&tokenResponse, store.salt)
	if err != nil {
		return nil, err
	}

	return &tokenResponse, nil
}
