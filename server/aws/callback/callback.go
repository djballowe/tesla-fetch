package main

import (
	"context"
	"encoding/json"
	"log"
	awshelpers "tesla-app/server/aws/helpers"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func authCallback(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("QueryStringParameters: %+v", event.QueryStringParameters)
	code := event.QueryStringParameters["code"]

	if code == "" {
		log.Println("Missing request data")
		return awshelpers.HandleAwsReturn("missing required params", 500, nil)
	}

	tokens, err := awshelpers.ExchangeCodeForToken(code)
	if err != nil || tokens == nil {
		log.Printf("Error missing tokens: %s", err)
		return awshelpers.HandleAwsReturn("missing tokens", 500, err)
	}

	token := awshelpers.Token{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		IdToken:      tokens.IdToken,
		State:        tokens.State,
		ExpiresIn:    tokens.ExpiresIn,
		TokenType:    tokens.TokenType,
	}

	body, err := json.Marshal(token)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Error creating response",
		}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type":              "application/json",
			"Cache-Control":             "no-store",
			"Strict-Transport-Security": "max-age=31536000; includeSubDomains",
			"X-Content-Type-Options":    "nosniff",
			"X-Frame-Options":           "DENY",
			"Content-Security-Policy":   "default-src 'none'",
			"X-XSS-Protection":          "1; mode=block",
		},
		Body:            string(body),
		IsBase64Encoded: false,
	}, nil
}

func main() {
	lambda.Start(authCallback)
}
