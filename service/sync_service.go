package service

import (
	"budget_tracket/client"
	"budget_tracket/constants"
	"budget_tracket/database/repository"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type SyncService struct {
	plaidClient     *client.PlaidClient
	plaidRepository *repository.PlaidRepository
}

func NewSyncService() (*SyncService, error) {
	op := "NewSyncService"
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

	service := SyncService{
		plaidClient:     plaidClient,
		plaidRepository: plaidRepsitory,
	}

	return &service, nil
}

func (s *SyncService) SyncTransactions(ctx context.Context) error {
	op := "SyncTransactions"

	items, err := s.plaidRepository.ListAllItems(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	for _, item := range items {
		added, _, _, cursor, err := s.plaidClient.SyncTransactions(ctx, item.AccessToken, item.Cursor)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		item.Cursor = cursor
		err = s.plaidRepository.UpdateItemCursor(ctx, item)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		if len(added) > 0 {
			err = s.plaidRepository.BulkCreateTransactions(ctx, added, &item.UserID)
			if err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}
		}

	}

	return nil
}
