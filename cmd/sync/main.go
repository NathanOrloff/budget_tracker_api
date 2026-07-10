package main

import (
	"budget_tracket/handler"
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

func handleRequest(ctx context.Context, event json.RawMessage) error {
	op := "handleRequest"

	handler, err := handler.NewSyncHandler()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	handler.SyncTransactions(ctx)

	return nil
}

func main() {
	lambda.Start(handleRequest)
}
