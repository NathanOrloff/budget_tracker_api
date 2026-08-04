package middleware

import (
	"budget_tracket/constants"
	"context"
	"log"

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
		claims, ok := apiGwContext.Authorizer["claims"].(map[string]interface{})
		if !ok {
			log.Printf("AuthMiddleware: no claims found in authorizer context: %+v", apiGwContext.Authorizer)
			c.AbortWithStatus(401)
			return
		}
		sub, ok := claims["sub"].(string)
		if !ok || sub == "" {
			log.Printf("AuthMiddleware: sub missing from claims: %+v", claims)
			c.AbortWithStatus(401)
			return
		}
		ctx := context.WithValue(c.Request.Context(), constants.USER_ID_KEY, sub)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
