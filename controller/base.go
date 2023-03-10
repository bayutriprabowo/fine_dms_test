package controller

import (
	"enigmacamp.com/fine_dms/model/dto"
	"github.com/gin-gonic/gin"
)

func SuccessJSONResponse(ctx *gin.Context, code int, resp dto.ApiResponse) {
	ctx.JSON(code, resp)
}

func FailedJSONResponse(ctx *gin.Context, code int, resp dto.ApiResponse) {
	ctx.JSON(code, resp)
}
