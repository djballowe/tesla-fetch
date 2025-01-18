package main

import (
	"context"
	"encoding/json"
	"log"
	awshelpers "tesla-app/server/aws/secrets"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func callbackHandler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("QueryStringParameters: %+v", event.QueryStringParameters)

	code := event.QueryStringParameters["code"]
	state := event.QueryStringParameters["state"]
	prevState := event.QueryStringParameters["prevState"]

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

	tokens, err := awshelpers.ExchangeCodeForToken(code)
	if err != nil || tokens.AccessToken == "" || tokens.RefreshToken == "" {
		log.Println("Could not gather token: ", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"message: could not gather token data"}`,
		}, nil
	}

	response, err := json.Marshal(tokens)
	if err != nil {
		log.Println("Could not marshal response body: ", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"message: could not gather token data"}`,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(response),
	}, nil
}

func main() {
	lambda.Start(callbackHandler)
}
