package models

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Item struct {
	PK              string    `dynamodbav:"PK"` // USER#<user_id>
	SK              string    `dynamodbav:"SK"` // #ITEM#<item_id>
	ID              string    `dynamodbav:"id"`
	UserID          string    `dynamodbav:"user_id"`
	AccessToken     string    `dynamodbav:"access_token"`
	Cursor          *string   `dynamodbav:"cursor"`
	InstitutionName string    `dynamodbav:"institution_name"`
	CreatedAt       time.Time `dynamodbav:"created_at"`
	UpdatedAt       time.Time `dynamodbav:"updated_at"`

	Accounts []Account `dynamodbav:"-"`
}

func (item *Item) MarshalKey() (map[string]types.AttributeValue, error) {
	key := struct {
		PK string
		SK string
	}{
		PK: item.PK,
		SK: item.SK,
	}
	return attributevalue.MarshalMap(key)
}
