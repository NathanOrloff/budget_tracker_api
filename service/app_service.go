package service

import (
	"budget_tracket/constants"
	"budget_tracket/database/repository"
	"fmt"
	"os"
)

type AppService struct {
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
	}

	return &service, nil
}
