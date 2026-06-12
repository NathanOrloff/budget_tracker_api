package models

import "time"

type Item struct {
	PK              string    `dynamodbav:"PK"`
	SK              string    `dynamodbav:"SK"`
	ID              string    `dynamodbav:"id"`
	UserID          string    `dynamodbav:"user_id"`
	AccessToken     string    `dynamodbav:"access_token"`
	Cursor          *string   `dynamodbav:"cursor"`
	InstitutionName string    `dynamodbav:"institution_name"`
	CreatedAt       time.Time `dynamodbav:"created_at"`
	UpdatedAt       time.Time `dynamodbav:"updated_at"`

	Accounts []Account `dynamodbav:"-"`
}
