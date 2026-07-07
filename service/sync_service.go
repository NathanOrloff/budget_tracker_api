package service

import (
	"budget_tracket/client"
	"budget_tracket/constants"
	"budget_tracket/database/repository"
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
