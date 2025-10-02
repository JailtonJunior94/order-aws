package dynamo

import (
	"context"
	"fmt"
	"strconv"

	"github.com/jailtonjunior94/order-aws/internal/domain/entities"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type sequenceRepository struct {
	tableName string
	client    *dynamodb.Client
}

func NewSequenceRepository(tableName string, cfg aws.Config) *sequenceRepository {
	return &sequenceRepository{
		tableName: tableName,
		client:    dynamodb.NewFromConfig(cfg),
	}
}

func (r *sequenceRepository) NextValue(ctx context.Context, sequence *entities.Sequence) (int64, error) {
	input := &dynamodb.UpdateItemInput{
		TableName:    aws.String(r.tableName),
		ReturnValues: types.ReturnValueAllNew,
		Key: map[string]types.AttributeValue{
			"date": &types.AttributeValueMemberS{Value: sequence.Date},
			"code": &types.AttributeValueMemberS{Value: sequence.Code},
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":q":         &types.AttributeValueMemberN{Value: strconv.FormatInt(sequence.Sequence, 10)},
			":expire_at": &types.AttributeValueMemberN{Value: strconv.FormatInt(sequence.ExpireAt, 10)},
		},
		UpdateExpression: aws.String(fmt.Sprintf("ADD %s :q SET %s = :expire_at", "sequence_number", "expire_at")),
	}

	output, err := r.client.UpdateItem(ctx, input)
	if err != nil {
		return 0, fmt.Errorf("sequence_repository: %v", err)
	}

	attr, ok := output.Attributes["sequence_number"].(*types.AttributeValueMemberN)
	if !ok {
		return 0, fmt.Errorf("sequence attribute missing or wrong type")
	}

	nextValue, err := strconv.ParseInt(attr.Value, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parse sequence value: %w", err)
	}
	return nextValue, nil
}
