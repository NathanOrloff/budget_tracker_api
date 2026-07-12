package handler

import (
	"budget_tracket/frontend"
	"budget_tracket/service"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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

func (a *AppHandler) CreateLinkToken(c *gin.Context) {
	op := "CreateLinkToken"
	ctx := c.Request.Context()

	token, err := a.appService.CreateLinkToken(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Errorf("%s: %w", op, err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"link_token": token,
	})
}

func (a *AppHandler) ExchangePublicToken(c *gin.Context) {
	op := "ExchangePublicToken"
	ctx := c.Request.Context()

	var input frontend.ExchangePublicTokenInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Errorf("%s: %w", op, err),
		})
		return
	}

	err = a.appService.ExchangePublicToken(ctx, input.PublicToken, input.InstitutionName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Errorf("%s: %w", op, err),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}

func (a *AppHandler) ListTransactionsSinceDate(c *gin.Context) {
	op := "ListTransactionsSinceDate"
	ctx := c.Request.Context()

	fromDateInput := c.Param("from_date")
	fromDate, err := time.Parse("2006-01-02", fromDateInput)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Errorf("%s: %w", op, err),
		})
		return
	}

	transactions, err := a.appService.ListTransactionsSinceDate(ctx, fromDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Errorf("%s: %w", op, err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": transactions,
	})
}
