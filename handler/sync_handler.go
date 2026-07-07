package handler

import (
	"budget_tracket/service"
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
