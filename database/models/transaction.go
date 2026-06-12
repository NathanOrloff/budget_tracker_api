package models

import "time"

type Transaction struct {
	PK                      string    `dynamodbav:"PK"`
	SK                      string    `dynamodbav:"SK"`
	GSI1PK                  string    `dynamodbav:"GSI1PK"`
	GSI1SK                  string    `dynamodbav:"GSI1SK"`
	TTL                     time.Time `dynamodbav:"ttl"`
	ID                      string    `dynamodbav:"id"`
	AccountID               string    `dynamodbav:"account_id"`
	Name                    string    `dynamodbav:"name"`
	MerchantName            string    `dynamodbav:"merchant_name"`
	Amount                  float32   `dynamodbav:"amount"`
	Date                    time.Time `dynamodbav:"date"`
	PersonalFinanceCategory string    `dynamodbav:"personal_finance_category"`
	CreatedAt               time.Time `dynamodbav:"created_at"`
	UpdatedAt               time.Time `dynamodbav:"updated_at"`
}
