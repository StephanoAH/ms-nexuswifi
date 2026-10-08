package middleware

import (
	"encoding/base64"
	"genesis/pos/domain/entities"
	"genesis/pos/presentation/container"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
)

func BasicAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Basic ") {
			c.JSON(http.StatusUnauthorized, entities.ResponseJson{
				Status:  http.StatusUnauthorized,
				Message: "Authorization header missing or malformed",
				Data:    nil,
			})
			c.Abort()
			return
		}

		payload, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(authHeader, "Basic "))
		if err != nil {
			c.JSON(http.StatusUnauthorized, entities.ResponseJson{
				Status:  http.StatusUnauthorized,
				Message: "Invalid Authorization header",
				Data:    nil,
			})
		}

		parts := strings.SplitN(string(payload), ":", 2)
		if len(parts) != 2 {
			c.JSON(http.StatusUnauthorized, entities.ResponseJson{
				Status:  http.StatusUnauthorized,
				Message: "Authorization header invalid or malformed",
				Data:    nil,
			})
		}

		username, password := parts[0], parts[1]

		credentials, err := container.GetCredentialsUseCase.GetCredentialsUseCase(1, username)
		if err != nil {
			c.JSON(http.StatusUnauthorized, entities.ResponseJson{
				Status:  http.StatusUnauthorized,
				Message: "Invalid uthorization",
				Data:    nil,
			})
		}

		if err = bcrypt.CompareHashAndPassword([]byte(credentials.Password), []byte(password)); err != nil {
			c.JSON(http.StatusUnauthorized, entities.ResponseJson{
				Status:  http.StatusUnauthorized,
				Message: "Invalid username or password",
				Data:    nil,
			})
		}

		c.Next()
	}
}
