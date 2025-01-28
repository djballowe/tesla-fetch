package main

import (
	"context"
	"log"
	awshelpers "tesla-app/server/aws/helpers"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
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

func pollTokensHandler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("QueryStringParameters: %+v", event)

	return awshelpers.HandleAwsReturn("success", 200, nil)
}

func main() {
	lambda.Start(pollTokensHandler)
}
