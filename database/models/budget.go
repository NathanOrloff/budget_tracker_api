package models

import (
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Budget struct {
	PK         string           `dynamodbav:"PK"` // USER#<user_id>
	SK         string           `dynamodbav:"SK"` // #BUDGET
	UserID     string           `dynamodbav:"user_id"`
	Categories []BudgetCategory `dynamodbav:"categories"`
}

type BudgetCategory struct {
	PersonalFinanceCategory string  `dynamodbav:"personal_finance_category"`
	Amount                  float64 `dynamodbav:"amount"`
}

func (budget *Budget) MarshalKey() (map[string]types.AttributeValue, error) {
	key := struct {
		PK string
		SK string
	}{
		PK: budget.PK,
		SK: budget.SK,
	}
	return attributevalue.MarshalMap(key)
}
