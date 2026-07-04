package models

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	plaid "github.com/plaid/plaid-go"
)

type Transaction struct {
	PK                      string    `dynamodbav:"PK"`     // ACCOUNT#<account_id>
	SK                      string    `dynamodbav:"SK"`     // TXN#<transaction_id>
	GSI1PK                  string    `dynamodbav:"GSI1PK"` // USER#<user_id>
	GSI1SK                  string    `dynamodbav:"GSI1SK"` // TXN#<transaction_id>
	TTL                     int64     `dynamodbav:"ttl"`
	ID                      string    `dynamodbav:"id"`
	AccountID               string    `dynamodbav:"account_id"`
	Name                    string    `dynamodbav:"name"`
	MerchantName            string    `dynamodbav:"merchant_name"`
	Amount                  float64   `dynamodbav:"amount"`
	Date                    time.Time `dynamodbav:"date"`
	PersonalFinanceCategory string    `dynamodbav:"personal_finance_category"`
	CreatedAt               time.Time `dynamodbav:"created_at"`
	UpdatedAt               time.Time `dynamodbav:"updated_at"`
}

func MarshalTransaction(plaidTransaction plaid.Transaction) (Transaction, error) {
	op := "Marshal"

	pfc := ""
	if plaidTransaction.PersonalFinanceCategory.IsSet() {
		if cat := plaidTransaction.PersonalFinanceCategory.Get(); cat != nil {
			pfc = cat.Primary
		}
	}

	date, err := time.Parse("2006-01-02", plaidTransaction.Date)
	if err != nil {
		return Transaction{}, fmt.Errorf("%s: %w", op, err)
	}

	merchantName := ""
	if plaidTransaction.MerchantName.IsSet() {
		if v := plaidTransaction.MerchantName.Get(); v != nil {
			merchantName = *v
		}
	}

	return Transaction{
		ID:                      plaidTransaction.TransactionId,
		AccountID:               plaidTransaction.AccountId,
		Name:                    plaidTransaction.Name,
		MerchantName:            merchantName,
		Amount:                  plaidTransaction.Amount,
		Date:                    date,
		PersonalFinanceCategory: pfc,
	}, nil
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
