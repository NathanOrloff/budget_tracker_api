package handler

import (
	"budget_tracket/service"
	"fmt"
)

type AppHandler struct {
	appService *service.AppService
}

func NewAppHandler() (*AppHandler, error) {
	op := "NewAppHandler"

	service, err := service.NewAppService()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	handler := AppHandler{
		appService: service,
	}

	return &handler, nil
}
