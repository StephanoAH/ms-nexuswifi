package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Validatable interface {
	Validate() error
}

func ValidateBody[T Validatable]() gin.HandlerFunc {
	return func(c *gin.Context) {
		var obj T
		if err := c.ShouldBindJSON(&obj); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "JSON invalid", "details": err.Error(), "success": false, "status": http.StatusBadRequest})
			c.Abort()
			return
		}
		if err := obj.Validate(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed validation", "details": err.Error(), "success": false, "status": http.StatusBadRequest})
			c.Abort()
			return
		}
		c.Set("body", &obj)
		c.Next()
	}
}

func ValidateQuery[T Validatable]() gin.HandlerFunc {
	return func(c *gin.Context) {
		var obj T
		if err := c.ShouldBindQuery(&obj); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Query invalid", "details": err.Error(), "success": false, "status": http.StatusBadRequest})
			c.Abort()
			return
		}
		if err := obj.Validate(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed query validation", "details": err.Error(), "success": false, "status": http.StatusBadRequest})
			c.Abort()
			return
		}
		c.Set("query", &obj)
		c.Next()
	}
}

func ValidateParams[T Validatable]() gin.HandlerFunc {
	return func(c *gin.Context) {
		var obj T
		if err := c.ShouldBindUri(&obj); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parameters", "details": err.Error(), "success": false, "status": http.StatusBadRequest})
			c.Abort()
			return
		}
		if err := obj.Validate(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Parameter validation failed", "details": err.Error(), "success": false, "status": http.StatusBadRequest})
			c.Abort()
			return
		}
		c.Set("params", &obj)
		c.Next()
	}
}
