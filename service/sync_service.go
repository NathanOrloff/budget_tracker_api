package service

import (
	"budget_tracket/client"
	"budget_tracket/constants"
	"budget_tracket/database/repository"
	"context"
	"fmt"
	"os"
)

type SyncService struct {
	plaidClient     *client.PlaidClient
	plaidRepository *repository.PlaidRepository
}

func NewSyncService() (*SyncService, error) {
	op := "NewSyncService"

	plaidRepsitory, err := repository.NewPlaidRepository(os.Getenv(constants.PLAID_DB_NAME_KEY))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	service := SyncService{
		plaidClient:     client.NewPlaidClient(),
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
			err = s.plaidRepository.BulkCreateTransactions(ctx, added)
			if err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}
		}

	}

	return nil
}
