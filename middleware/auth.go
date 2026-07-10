package middleware

import (
	"budget_tracket/constants"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	region := os.Getenv(constants.ENV_REGION)
	userPoolID := os.Getenv(constants.ENV_USER_POOL_ID)

	jwksURL := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json", region, userPoolID)
	k, err := keyfunc.NewDefault([]string{jwksURL})
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize JWKS keyfunc: %v", err))
	}

	unauthorizedAccessResponse := gin.H{
		"error": "Unauthorized Access",
	}

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, unauthorizedAccessResponse)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, k.Keyfunc)
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, unauthorizedAccessResponse)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, unauthorizedAccessResponse)
			return
		}

		subject, err := claims.GetSubject()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, unauthorizedAccessResponse)
			return
		}

		c.Set(constants.USER_ID_KEY, subject)
		c.Next()
	}
}
