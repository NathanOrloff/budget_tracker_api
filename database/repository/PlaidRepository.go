package repository

import (
	"budget_tracket/database/models"
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type PlaidRepository struct {
	Client    *dynamodb.Client
	TableName string
}

func NewPlaidClient(ctx context.Context, profile string) (*dynamodb.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	client := dynamodb.NewFromConfig(cfg)
	return client, nil
}

func (plaidRepository *PlaidRepository) BulkCreateTransactions(ctx context.Context, transactions []models.Transaction) error {
	op := "BulkCreateTransactions"
	batchSize := 25 // DynamoDB max allowed batch size

	start := 0
	end := start + batchSize
	for start < len(transactions) {
		var writeReqs []types.WriteRequest
		if end > len(transactions) {
			end = len(transactions)
		}
		for _, transaction := range transactions[start:end] {
			item, err := attributevalue.MarshalMap(transaction)
			if err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}
			writeReqs = append(writeReqs, types.WriteRequest{PutRequest: &types.PutRequest{Item: item}})
		}

		_, err := plaidRepository.Client.BatchWriteItem(ctx, &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]types.WriteRequest{plaidRepository.TableName: writeReqs},
		})
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		start = end
		end += batchSize
	}
	return nil
}

func (plaidRepository *PlaidRepository) ListTransactionsByUserID(ctx context.Context, userID string, startDate *time.Time, endDate *time.Time) ([]models.Transaction, error) {
	op := "ListTransactionsByUserID"
	keyEx := expression.Key("GSI1PK").Equal(expression.Value("USER#" + userID))

	var filter expression.ConditionBuilder
	if startDate != nil && endDate != nil {
		filter = expression.Name("date").Between(
			expression.Value(startDate.Format(time.RFC3339)),
			expression.Value(endDate.Format(time.RFC3339)),
		)
	}

	expr, err := expression.NewBuilder().
		WithKeyCondition(keyEx).
		WithFilter(filter).
		Build()
	if err != nil {
		return []models.Transaction{}, fmt.Errorf("%s: %w", op, err)
	}

	var transactions []models.Transaction
	paginator := dynamodb.NewQueryPaginator(plaidRepository.Client, &dynamodb.QueryInput{
		TableName:                 plaidRepository.TableName,
		IndexName:                 aws.String("GSI1"),
		KeyConditionExpression:    expr.KeyCondition(),
		FilterExpression:          expr.Filter(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return []models.Transaction{}, fmt.Errorf("%s: %w", op, err)
		}
		var page_txns []models.Transaction
		if err := attributevalue.UnmarshalListOfMaps(page.Items, &page_txns); err != nil {
			return []models.Transaction{}, fmt.Errorf("%s: %w", op, err)
		}
		transactions = append(transactions, page_txns...)
	}

	return transactions, nil
}
