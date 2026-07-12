package frontend

import (
	"budget_tracket/database/models"
	"time"
)

type TransactionOutput struct {
	ID                      string    `json:"id"`
	AccountID               string    `json:"account_id"`
	Name                    string    `json:"name"`
	MerchantName            string    `json:"merchant_name"`
	Amount                  float64   `json:"amount"`
	Date                    time.Time `json:"date"`
	PersonalFinanceCategory string    `json:"personal_finance_category"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
}

func MarshalTransaction(transaction models.Transaction) TransactionOutput {
	return TransactionOutput{
		ID:                      transaction.ID,
		AccountID:               transaction.AccountID,
		Name:                    transaction.Name,
		MerchantName:            transaction.MerchantName,
		Amount:                  transaction.Amount,
		Date:                    transaction.Date,
		PersonalFinanceCategory: transaction.PersonalFinanceCategory,
		CreatedAt:               transaction.CreatedAt,
		UpdatedAt:               transaction.UpdatedAt,
	}
}
