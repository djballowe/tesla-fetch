package main

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func callbackHandler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("QueryStringParameters: %+v", event.QueryStringParameters)

	code := event.QueryStringParameters["code"]
	state := event.QueryStringParameters["state"]
	prevState := event.QueryStringParameters["prevState"]

	//	authStatus := false
	if code == "" || state == "" {
		log.Println("Missing request data")
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"message: missing required params"}`,
		}, nil
	}

	if prevState != state {
		log.Println("State does not match")
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"message: state does not match"}`,
		}, nil
	}

	_, err := exchangeCodeForToken(code)
	if err != nil {
		log.Println("Could not gather token: ", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"message: could not gather token data"}`,
		}, nil
	}
//	log.Println("token: ", tokens.AccessToken)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       `{"message: success"}`,
	}, nil
}

func exchangeCodeForToken(code string) (*Token, error) {
	log.Println("Exchanging code for token...")

	baseUrl, err := url.Parse("https://auth.tesla.com/oauth2/v3/token")
	if err != nil {
		return nil, err
	}
	tokenUrl := baseUrl.String()

	secrets, err := getSecrets("tesla_creds")
	if err != nil {
		return nil, err
	}

	log.Println("Token url: ", tokenUrl)
	log.Println("clientId: ", secrets)

	return nil, nil
}

func getSecrets(secret string) (string, error) {
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

func main() {
	lambda.Start(callbackHandler)
}
