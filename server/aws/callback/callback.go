package main

import (
	"context"
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
		return awshelpers.HandleAwsReturn("missing required params", 500, nil)
	}

	tokens, err := awshelpers.ExchangeCodeForToken(code)
	if err != nil || tokens == nil {
		log.Printf("Error missing tokens: %s", err)
		return awshelpers.HandleAwsReturn("missing tokens", 500, err)
	}

	accessTokenEncrypt, err := awshelpers.EncryptKey(tokens.AccessToken)
	refreshTokenEncrypt, err := awshelpers.EncryptKey(tokens.RefreshToken)
	idTokenEncrypt, err := awshelpers.EncryptKey(tokens.IdToken)
	if err != nil {
		return awshelpers.HandleAwsReturn("could not encrypt tokens", 500, err)
	}

	token := awshelpers.Token{
		AccessToken:  accessTokenEncrypt,
		RefreshToken: refreshTokenEncrypt,
		IdToken:      idTokenEncrypt,
		State:        tokens.State,
		ExpiresIn:    tokens.ExpiresIn,
		TokenType:    tokens.TokenType,
	}

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

	_, err = dynamoDBClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	})

	return awshelpers.HandleAwsReturn("tokens store successful", 200, err)
}

func main() {
	lambda.Start(authCallback)
}
