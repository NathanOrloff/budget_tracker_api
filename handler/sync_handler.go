package handler

import (
	"budget_tracket/service"
	"context"
	"fmt"
)

type SyncHandler struct {
	syncService *service.SyncService
}

func NewSyncHandler() (*SyncHandler, error) {
	op := "NewSyncHandler"

	service, err := service.NewSyncService()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	handler := SyncHandler{
		syncService: service,
	}

	return &handler, nil
}

func (s *SyncHandler) SyncTransactions(ctx context.Context) error {
	op := "SyncTransactions"

	err := s.syncService.SyncTransactions(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
