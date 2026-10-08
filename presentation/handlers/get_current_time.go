package handlers

import (
	"genesis/pos/domain/entities"
	"genesis/pos/presentation/container"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetCurrentTime(ctx *gin.Context) {
	/*
		var requestBody dtos.MovementDto
		if err := ctx.BindJSON(&requestBody); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	*/
	response := container.GetCurrentTimeService.ExecuteService()
	res := &entities.ResponseHttp{}

	res.Code = response.StatusCode
	res.Message = response.Result.String()

	if response.StatusCode != 200 {
		res.Success = false
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res.Success = true
	ctx.JSON(http.StatusOK, res)
	return
}
