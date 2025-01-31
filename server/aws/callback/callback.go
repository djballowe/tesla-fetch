package main

import (
	"context"
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

	log.Println(token)

	return awshelpers.HandleAwsReturn("tokens store successful", 200, err)
}

func main() {
	lambda.Start(authCallback)
}
