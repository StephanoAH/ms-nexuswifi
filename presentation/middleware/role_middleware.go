package middleware

import (
	"genesis/pos/domain/entities"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RoleMiddleware(allowedRoles ...int) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, exists := c.Get("claims")
		if !exists {
			c.JSON(http.StatusUnauthorized, entities.ResponseJson{
				Status:  http.StatusBadRequest,
				Message: "Unauthorized",
				Data:    nil,
			})
			c.Abort()
			return
		}

		claimsMap, ok := claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusForbidden, entities.ResponseJson{
				Status:  http.StatusBadRequest,
				Message: "Invalid claims",
				Data:    nil,
			})
			c.Abort()
			return
		}

		userRole, ok := claimsMap["role"].(float64)
		if !ok {
			c.JSON(http.StatusForbidden, entities.ResponseJson{
				Status:  http.StatusBadRequest,
				Message: "Invalid role",
				Data:    nil,
			})
			c.Abort()
			return
		}

		for _, role := range allowedRoles {
			if int(userRole) == role {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, entities.ResponseJson{
			Status:  http.StatusForbidden,
			Message: "Access denied",
			Data:    nil,
		})
		c.Abort()
	}
}
