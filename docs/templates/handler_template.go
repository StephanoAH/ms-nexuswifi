package templates

import (
    "genesis/pos/domain/entities"
    "genesis/pos/presentation/container"
    "net/http"

    "github.com/gin-gonic/gin"
)

// [Nombre] maneja el request para [Descripción]
func [Nombre](ctx *gin.Context) {
    // Opcional: Obtener body validado del middleware
    // body, _ := ctx.Get("body").(*schema.[Nombre]Body)
    
    // Opcional: Obtener params validados del middleware
    // params, _ := ctx.Get("params").(*schema.[Nombre]Params)
    
    // Opcional: Obtener query validados del middleware
    // query, _ := ctx.Get("query").(*schema.[Nombre]Query)

    response := container.[Nombre]Service.ExecuteService()
    res := &entities.ResponseHttp{}

    res.Code = response.StatusCode
    res.Message = response.[CampoMensaje]

    if response.StatusCode != 200 {
        res.Success = false
        ctx.JSON(http.StatusBadRequest, res)
        return
    }

    res.Success = true
    ctx.JSON(http.StatusOK, res)
    return
}

