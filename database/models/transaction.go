package models

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

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

func (transaction *Transaction) MarshalKey() (map[string]types.AttributeValue, error) {
	key := struct {
		PK     string
		SK     string
		GSI1PK string
		GSI1SK string
	}{
		PK:     transaction.PK,
		SK:     transaction.SK,
		GSI1PK: transaction.GSI1PK,
		GSI1SK: transaction.GSI1SK,
	}
	return attributevalue.MarshalMap(key)
}
