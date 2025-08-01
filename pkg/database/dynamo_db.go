package database

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type DynamoDBClient struct {
	TableName      string
	DynamoDBClient *dynamodb.Client
}

func NewDynamoDBClient(ctx context.Context, sdkConfig aws.Config, tableName string) (*DynamoDBClient, error) {
	return &DynamoDBClient{
		TableName:      tableName,
		DynamoDBClient: dynamodb.NewFromConfig(sdkConfig),
	}, nil
}
