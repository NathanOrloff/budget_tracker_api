package models

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Account struct {
	PK        string    `dynamodbav:"PK"` // ITEM#<item_id>
	SK        string    `dynamodbav:"SK"` // ACCOUNT#<account_id>
	ID        string    `dynamodbav:"id"`
	ItemID    string    `dynamodbav:"item_id"`
	Name      string    `dynamodbav:"name"`
	Type      string    `dynamodbav:"type"`
	Subtype   string    `dynamodbav:"subtype"`
	CreatedAt time.Time `dynamodbav:"created_at"`
	UpdatedAt time.Time `dynamodbav:"updated_at"`

	// associations
	Transactions []Transaction `dynamodbav:"-"`
}

func (account *Account) MarshalKey() (map[string]types.AttributeValue, error) {
	key := struct {
		PK string
		SK string
	}{
		PK: account.PK,
		SK: account.SK,
	}
	return attributevalue.MarshalMap(key)
}
