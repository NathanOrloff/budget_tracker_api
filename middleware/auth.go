package middleware

import (
	"budget_tracket/constants"
	"context"

	"github.com/awslabs/aws-lambda-go-api-proxy/core"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiGwContext, ok := core.GetAPIGatewayContextFromContext(c.Request.Context())
		if !ok {
			c.AbortWithStatus(401)
			return
		}
		sub, _ := apiGwContext.Authorizer["claims"].(map[string]interface{})["sub"].(string)
		ctx := context.WithValue(c.Request.Context(), constants.USER_ID_KEY, sub)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
