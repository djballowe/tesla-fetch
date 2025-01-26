package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	awshelpers "tesla-app/server/aws/helpers"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var dynamoDBClient *dynamodb.Client
var tableName = "token_store"

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("unable to load aws config")
	}

	dynamoDBClient = dynamodb.NewFromConfig(cfg)
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

	token := awshelpers.Token{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		IdToken:      tokens.IdToken,
		State:        tokens.State,
		ExpiresIn:    tokens.ExpiresIn,
		TokenType:    tokens.TokenType,
	}
	log.Println(token)
	log.Println("Auth successful")

	now := time.Now().UTC()

	item := map[string]types.AttributeValue{
		"access_token":  &types.AttributeValueMemberS{Value: token.AccessToken},
		"refresh_token": &types.AttributeValueMemberS{Value: token.RefreshToken},
		"id_token":      &types.AttributeValueMemberS{Value: token.IdToken},
		"state":         &types.AttributeValueMemberS{Value: token.State},
		"expire_in":     &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", token.ExpiresIn)},
		"token_type":    &types.AttributeValueMemberS{Value: token.TokenType},
		"created_at":    &types.AttributeValueMemberS{Value: now.Format(time.RFC3339)},
	}

	dbResp, err := dynamoDBClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	})

	log.Println("dbResp: ", dbResp)

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
