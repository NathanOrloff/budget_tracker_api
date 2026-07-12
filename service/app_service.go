package service

import (
	"budget_tracket/client"
	"budget_tracket/constants"
	"budget_tracket/database/models"
	"budget_tracket/database/repository"
	"budget_tracket/frontend"
	"budget_tracket/utils"
	"context"
	"fmt"
	"os"
	"time"
)

type AppService struct {
	plaidClient     *client.PlaidClient
	plaidRepository *repository.PlaidRepository
}

func NewAppService() (*AppService, error) {
	op := "NewAppService"

	plaidRepsitory, err := repository.NewPlaidRepository(os.Getenv(constants.PLAID_DB_NAME_KEY))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	service := AppService{
		plaidRepository: plaidRepsitory,
		plaidClient:     client.NewPlaidClient(),
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
