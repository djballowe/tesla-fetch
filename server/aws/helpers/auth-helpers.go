package awshelpers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Secrets struct {
	ClientSecret string `json:"client_secret"`
	ClientId     string `json:"client_id"`
	Audience     string `json:"audience"`
	RedirectUri  string `json:"redirect_uri"`
	Scope        string `json:"scope"`
}

func GetSecrets(secret string) (string, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return "", fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := secretsmanager.NewFromConfig(cfg)

	output, err := client.GetSecretValue(context.TODO(), &secretsmanager.GetSecretValueInput{
		SecretId: &secret,
	})
	if err != nil {
		return "", fmt.Errorf("failed to retrieve secret: %w", err)
	}

	return *output.SecretString, nil
}

func ExchangeCodeForToken(code string) (*Token, error) {
	log.Println("Exchanging code for token...")

	baseUrl, err := url.Parse("https://auth.tesla.com/oauth2/v3/token")
	if err != nil {
		return nil, err
	}
	tokenUrl := baseUrl.String()

	secrets, err := GetSecrets("tesla_creds")
	if err != nil {
		return nil, err
	}

	var config Secrets
	if err := json.Unmarshal([]byte(secrets), &config); err != nil {
		return nil, err
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
