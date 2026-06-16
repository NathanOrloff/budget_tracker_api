package models

import "time"

type Transaction struct {
	PK                      string    `dynamodbav:"PK"`     // ACCOUNT#<account_id>
	SK                      string    `dynamodbav:"SK"`     // TXN#<transaction_id>
	GSI1PK                  string    `dynamodbav:"GSI1PK"` // USER#<user_id>
	GSI1SK                  string    `dynamodbav:"GSI1SK"` // TXN#<transaction_id>
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
