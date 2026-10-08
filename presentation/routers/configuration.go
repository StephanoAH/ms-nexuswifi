package routers

import (
	"genesis/pos/presentation/handlers"
	"genesis/pos/presentation/middleware"
	"genesis/pos/presentation/schema"

	"github.com/gin-gonic/gin"
)

func CurrentTime(router *gin.RouterGroup) {
	// example
	router.POST("/:id",
		middleware.ValidateParams[schema.CurrentTimeParams](),
		middleware.ValidateQuery[schema.CurrentTimeQuery](),
		middleware.ValidateBody[schema.CurrentTimeBody](),
		handlers.GetCurrentTime)
}
