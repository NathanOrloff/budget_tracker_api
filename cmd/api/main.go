package main

import (
	"budget_tracket/handler"
	"budget_tracket/middleware"
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
)

var ginLambda *ginadapter.GinLambda

func init() {

	handler, err := handler.NewAppHandler()
	if err != nil {
		log.Fatalf("failed to initialize app handler: %v", err)
	}

	r := gin.Default()
	r.Use(middleware.CorsMiddleware())
	r.Use(middleware.AuthMiddleware())

	r.GET("/create-link-token", handler.CreateLinkToken)
	r.GET("/transactions", handler.ListTransactionsSinceDate)
	r.GET("/account-is-registered", handler.AccountIsRegistered)

	r.POST("/exchange-public-token", handler.ExchangePublicToken)

	ginLambda = ginadapter.New(r)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(Handler)
}
