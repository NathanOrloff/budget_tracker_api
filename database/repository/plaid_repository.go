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

func NewPlaidRepository(tableName string) (*PlaidRepository, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	client := dynamodb.NewFromConfig(cfg)

	repo := PlaidRepository{
		Client:    client,
		TableName: tableName,
	}
	return &repo, nil
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
			transaction.FillKey(ctx)
			transaction.CreatedAt = time.Now()
			transaction.UpdatedAt = time.Now()
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
		TableName:                 &plaidRepository.TableName,
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

func (plaidRepository *PlaidRepository) CreateItem(ctx context.Context, itm models.Item) error {
	op := "CreateItem"

	itm.CreatedAt = time.Now()
	itm.UpdatedAt = time.Now()

	item, err := attributevalue.MarshalMap(itm)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = plaidRepository.Client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(plaidRepository.TableName), Item: item,
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (plaidRepository *PlaidRepository) UpdateItemCursor(ctx context.Context, itm models.Item) error {
	op := "UpdateItem"

	key, err := itm.MarshalKey()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	updateExpr := expression.Set(
		expression.Name("cursor"),
		expression.Value(itm.Cursor),
	).Set(
		expression.Name("updated_at"),
		expression.Value(time.Now()),
	)

	expr, err := expression.NewBuilder().WithUpdate(updateExpr).Build()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	input := &dynamodb.UpdateItemInput{
		TableName:                 aws.String(plaidRepository.TableName),
		Key:                       key,
		UpdateExpression:          expr.Update(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ReturnValues:              types.ReturnValueUpdatedNew,
	}

	_, err = plaidRepository.Client.UpdateItem(ctx, input)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (plaidRepository *PlaidRepository) ListAllItems(ctx context.Context) ([]models.Item, error) {
	op := "ListAllItems"
	keyEx := expression.Name("SK").BeginsWith("ITEM#")
	expr, err := expression.NewBuilder().WithFilter(keyEx).Build()
	if err != nil {
		return []models.Item{}, fmt.Errorf("%s: %w", op, err)
	}

	var items []models.Item
	var lastEvaluatedKey map[string]types.AttributeValue

	for {
		out, err := plaidRepository.Client.Scan(ctx, &dynamodb.ScanInput{
			TableName:                 aws.String(plaidRepository.TableName),
			FilterExpression:          expr.Filter(),
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
		})
		if err != nil {
			return []models.Item{}, fmt.Errorf("%s: %w", op, err)
		}

		var page []models.Item
		if err = attributevalue.UnmarshalListOfMaps(out.Items, &page); err != nil {
			return []models.Item{}, fmt.Errorf("%s: %w", op, err)
		}

		items = append(items, page...)

		lastEvaluatedKey = out.LastEvaluatedKey
		if lastEvaluatedKey == nil {
			break
		}
	}

	return items, nil
}

func (plaidRepository *PlaidRepository) ListItemsByUserID(ctx context.Context, userID string) ([]models.Item, error) {
	op := "ListItemsByUserID"
	keyEx := expression.Key("PK").Equal(expression.Value("USER#" + userID))

	expr, err := expression.NewBuilder().
		WithKeyCondition(keyEx).
		Build()
	if err != nil {
		return []models.Item{}, fmt.Errorf("%s: %w", op, err)
	}

	var items []models.Item
	paginator := dynamodb.NewQueryPaginator(plaidRepository.Client, &dynamodb.QueryInput{
		TableName:                 &plaidRepository.TableName,
		IndexName:                 aws.String("PK"),
		KeyConditionExpression:    expr.KeyCondition(),
		FilterExpression:          expr.Filter(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return []models.Item{}, fmt.Errorf("%s: %w", op, err)
		}
		var page_items []models.Item
		if err := attributevalue.UnmarshalListOfMaps(page.Items, &page_items); err != nil {
			return []models.Item{}, fmt.Errorf("%s: %w", op, err)
		}
		items = append(items, page_items...)
	}

	return items, nil

}

func (plaidRepository *PlaidRepository) CreateAccount(ctx context.Context, account models.Account) error {
	op := "CreateAccount"

	item, err := attributevalue.MarshalMap(account)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = plaidRepository.Client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(plaidRepository.TableName), Item: item,
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if len(account.Transactions) > 0 {
		err = plaidRepository.BulkCreateTransactions(ctx, account.Transactions)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil
}
