package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/url"
	awshelpers "tesla-app/server/aws/helpers"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type AuthResponse struct {
	CallbackURL string `json:"callback_url"`
}

func getTeslaAuth(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	baseUrl, err := url.Parse("https://auth.tesla.com/oauth2/v3/authorize")
	if err != nil {
		log.Fatal("Malformed auth url", err)
	}

	state := generateState()
	secrets, err := awshelpers.GetSecrets("tesla_creds")
	if err != nil {
		return awshelpers.HandleAwsReturn("could not get secrets", 500, err)
	}

	var config awshelpers.Secrets
	if err := json.Unmarshal([]byte(secrets), &config); err != nil {
		return awshelpers.HandleAwsReturn("could unmarshal secrets", 500, err)
	}

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

	respBody := AuthResponse{
		CallbackURL: authUrl,
	}

	body, err := json.Marshal(respBody)
	if err != nil {
		return awshelpers.HandleAwsReturn("could not marshal response body", 500, err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body:            string(body),
		IsBase64Encoded: false,
	}, nil
}

func generateState() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func main() {
	lambda.Start(getTeslaAuth)
}
