package callback 

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
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

	tokens, err := callAuth(code)
	if err != nil || tokens == nil {
		log.Println("Missing tokens")
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

func callAuth(code string) (*Token, error) {
	appUrl := os.Getenv("APP_BASE_URL")
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	client := http.Client{}
	url := fmt.Sprintf("%s/auth?code=%s", appUrl, code)

	fmt.Println(url)
	authRequest, err := http.NewRequest("POST", url, nil)
	authRequest.Header.Set("Content-Type", "application/json")

	signer := v4.NewSigner()
	creds, err := cfg.Credentials.Retrieve(context.TODO())
	if err != nil {
		return nil, err
	}

	// hash for an empty payload
	hash := "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	err = signer.SignHTTP(context.TODO(), creds, authRequest, hash, "execute-api", cfg.Region, time.Now())
	if err != nil {
		return nil, err
	}

	response, err := client.Do(authRequest)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, _ := io.ReadAll(response.Body)
	fmt.Println("Response:", string(body))

	var tokens Token

	err = json.Unmarshal(body, &tokens)
	if err != nil {
		return nil, err
	}

	return &tokens, nil
}

func main() {
	lambda.Start(authCallback)
}
