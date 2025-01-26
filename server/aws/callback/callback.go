package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	awshelpers "tesla-app/server/aws/helpers"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func authCallback(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("QueryStringParameters: %+v", event.QueryStringParameters)
	code := event.QueryStringParameters["code"]

	if code == "" {
		log.Println("Missing request data")
		return handleReturn("missing required params", 500, nil)
	}

	tokens, err := awshelpers.ExchangeCodeForToken(code)
	if err != nil || tokens == nil {
		log.Printf("Error missing tokens: %s", err)
		return handleReturn("missing tokens", 500, err)
	}

	token := Token{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}
	log.Println(token)
	log.Println("Auth successful")

	passTokenUrl := os.Getenv("TOKEN_URL")
	log.Println(passTokenUrl)

	client := &http.Client{}

	payload, err := json.Marshal(token)
	if err != nil {
		log.Printf("Could not marshal JSON payload: %s", err)
		return handleReturn("could not marshal JSON payload", 500, err)
	}

	req, err := http.NewRequest("POST", passTokenUrl, bytes.NewBuffer(payload))
	if err != nil {
		log.Printf("Could not create token url: %s", err)
		return handleReturn("could not create token url", 500, err)
	}

	response, err := client.Do(req)
	if err != nil {
		log.Printf("could not get token post response: %s", err)
		return handleReturn("could not get token post response", 500, err)
	}
	defer response.Body.Close()

	return handleReturn("Auth successful", 200, err)
}

func handleReturn(message string, statusCode int, err error) (events.APIGatewayProxyResponse, error) {
	msg, _ := json.Marshal(map[string]string{"message": message})
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(msg),
	}, err
}

func main() {
	lambda.Start(authCallback)
}
