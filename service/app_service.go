package service

import (
	"budget_tracket/client"
	"budget_tracket/constants"
	"budget_tracket/database/models"
	"budget_tracket/database/repository"
	"budget_tracket/frontend"
	"budget_tracket/utils"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type AppService struct {
	plaidClient     *client.PlaidClient
	plaidRepository *repository.PlaidRepository
}

func NewAppService() (*AppService, error) {
	op := "NewAppService"
	region := os.Getenv(constants.ENV_REGION)

	config, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	svc := secretsmanager.NewFromConfig(config)

	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(constants.PLAID_SECRET),
	}

	result, err := svc.GetSecretValue(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var secretsManagerMap map[string]string
	err = json.Unmarshal([]byte(*result.SecretString), &secretsManagerMap)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	plaidRepsitory, err := repository.NewPlaidRepository(os.Getenv(constants.PLAID_TABLE_NAME))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	plaidClient := client.NewPlaidClient(secretsManagerMap[constants.PLAID_SECRET_KEY])

	service := AppService{
		plaidRepository: plaidRepsitory,
		plaidClient:     plaidClient,
	}

	return &service, nil
}

func (a *AppService) CreateLinkToken(ctx context.Context) (string, error) {
	op := "CreateLinkToken"

	userID := utils.GetUIDFromCtx(ctx)
	if userID == "" {
		return "", fmt.Errorf("%s: Invalid userID", op)
	}

	token, err := a.plaidClient.CreateLinkToken(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

func (a *AppService) ExchangePublicToken(ctx context.Context, publicToken string, institution_name string) error {
	op := "ExchangePublicToken"

	token, itemID, err := a.plaidClient.ExchangePublicToken(ctx, publicToken)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	userID := utils.GetUIDFromCtx(ctx)
	if userID == "" {
		return fmt.Errorf("%s: Invalid userID", op)
	}

	newItem := models.Item{
		PK:              "USER#" + userID,
		SK:              "ITEM#" + itemID,
		ID:              itemID,
		UserID:          userID,
		AccessToken:     token,
		Cursor:          nil,
		InstitutionName: institution_name,
	}

	err = a.plaidRepository.CreateItem(ctx, newItem)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *AppService) ListTransactionsSinceDate(ctx context.Context, fromDate time.Time) ([]frontend.TransactionOutput, error) {
	op := "ListTransactionsSinceDate"

	userID := utils.GetUIDFromCtx(ctx)
	if userID == "" {
		return []frontend.TransactionOutput{}, fmt.Errorf("%s: Invalid userID", op)
	}

	currentDate := time.Now()
	dbTransactions, err := a.plaidRepository.ListTransactionsByUserID(ctx, userID, &fromDate, &currentDate)
	if err != nil {
		return []frontend.TransactionOutput{}, fmt.Errorf("%s: %w", op, err)
	}

	var transactions []frontend.TransactionOutput
	for _, dbTransaction := range dbTransactions {
		transactions = append(transactions, frontend.MarshalTransaction(dbTransaction))
	}

	return transactions, nil
}

func (a *AppService) AccountIsRegistered(ctx context.Context) (bool, error) {
	op := "AccountIsRegistered"

	userID := utils.GetUIDFromCtx(ctx)
	if userID == "" {
		return false, fmt.Errorf("%s: Invalid userID", op)
	}

	items, err := a.plaidRepository.ListItemsByUserID(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	if len(items) == 0 {
		return false, nil
	}
	return true, nil
}
