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

	_, err := handler.NewSyncHandler()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func main() {
	lambda.Start(handleRequest)
}
