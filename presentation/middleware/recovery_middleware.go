package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Puedes agregar aquí tu lógica de logging
				fmt.Printf("Panic recovered: %v\n", err)

				// Respuesta a la solicitud cuando ocurre un panic
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Internal Server Error",
				})

				// Finaliza la solicitud para evitar continuar con los siguientes handlers
				c.Abort()
			}
		}()

		// Continúa con la ejecución del resto de la cadena de middleware y handlers
		c.Next()
	}
}
