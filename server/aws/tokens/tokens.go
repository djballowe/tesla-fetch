package main

import (
	"context"
	"log"
	awshelpers "tesla-app/server/aws/helpers"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
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

func pollTokensHandler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("QueryStringParameters: %+v", event)
	state := event.QueryStringParameters["state"]

	token := awshelpers.Token{}

	result, err := dynamoDBClient.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"state": &types.AttributeValueMemberS{Value: state},
		},
	})
	if err != nil {
		return awshelpers.HandleAwsReturn("error finding state", 500, err)
	}
	if result.Item == nil {
		return awshelpers.HandleAwsReturn("could not find state", 200, nil)
	}

	err = attributevalue.UnmarshalMap(result.Item, &token)
	log.Println("item: ", result.Item)

	return awshelpers.HandleAwsReturn("success", 200, nil)
}

func main() {
	lambda.Start(pollTokensHandler)
}
