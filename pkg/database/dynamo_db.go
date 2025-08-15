package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type (
	DynamoDBClient interface {
		PutItem(ctx context.Context, item DynoNotation) error
		GetItem(ctx context.Context, key DynoNotation) (DynoNotation, error)
	}

	dynamoDBClient struct {
		tableName      string
		DynamoDBClient *dynamodb.Client
	}

	DynoNotation map[string]types.AttributeValue
)

func NewDynamoDBClient(ctx context.Context, sdkConfig aws.Config, tableName string) (*dynamoDBClient, error) {
	if len(tableName) == 0 {
		return nil, errors.New("dynamo_client: table name cannot be empty")
	}

	return &dynamoDBClient{tableName: tableName, DynamoDBClient: dynamodb.NewFromConfig(sdkConfig)}, nil
}

func (d *dynamoDBClient) PutItem(ctx context.Context, item DynoNotation) error {
	if len(item) == 0 {
		return errors.New("dynamo_client: input item cannot be empty")
	}

	_, err := d.DynamoDBClient.PutItem(ctx, &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(d.tableName),
	})

	if err != nil {
		return fmt.Errorf("dynamo_client: %v", err)
	}
	return nil
}

func (d *dynamoDBClient) GetItem(ctx context.Context, key DynoNotation) (DynoNotation, error) {
	if len(key) == 0 {
		return nil, errors.New("dynamo_client: input key cannot be empty")
	}

	result, err := d.DynamoDBClient.GetItem(ctx, &dynamodb.GetItemInput{
		Key:       key,
		TableName: aws.String(d.tableName),
	})

	if err != nil {
		return nil, fmt.Errorf("dynamo_client: %v", err)
	}

	if result.Item == nil {
		return nil, nil
	}

	return result.Item, nil
}
