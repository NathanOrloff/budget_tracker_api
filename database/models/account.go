package models

import "time"

type Account struct {
	PK        string    `dynamodbav:"PK"`
	SK        string    `dynamodbav:"SK"`
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
