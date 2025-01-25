package main

import (
	"context"
	"log"
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

	// authStatus := false
	if code == "" {
		log.Println("Missing request data")
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"message: missing required params"}`,
		}, nil
	}

	// storeMutex.Lock()
	// storedState := stateStore
	// fmt.Println(storedState, "||", state)
	// if storedState != state {
	// 	http.Error(writer, "State does not match", http.StatusBadRequest)
	// 	return
	// }
	// storeMutex.Unlock()

	tokens, err := awshelpers.ExchangeCodeForToken(code)
	if err != nil || tokens == nil {
		log.Printf("Error missing tokens: %s", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"message: missing tokens"}`,
		}, nil
	}

	// storeMutex.Lock()
	token := Token{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}
	// tokenStore[state] = token
	// storeMutex.Unlock()

	// fmt.Fprintf(writer, "Auth successful\n")
	log.Println(token)
	log.Println("Auth successful")
	// authStatus = true

	// notifyClientUrl := fmt.Sprintf("http://localhost:3000/notify?auth_status=%t", authStatus)

	// _, err = http.Post(notifyClientUrl, "application/x-www-form-urlencoded", nil)
	// if err != nil {
	// 	http.Error(writer, fmt.Sprintf("Failed to update the client of auth status: %s", err.Error()), http.StatusInternalServerError)
	// 	return
	// }

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: "Success",
	}, nil
}

func main() {
	lambda.Start(authCallback)
}
