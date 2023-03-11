package controller

import (
	"enigmacamp.com/fine_dms/model/dto"
	"github.com/gin-gonic/gin"
)

func SuccessJSONResponse(ctx *gin.Context, code int, msg string, data any) {
	ctx.JSON(code, dto.NewApiResponseSuccess(msg, data))
}

func FailedJSONResponse(ctx *gin.Context, code int, msg string) {
	ctx.JSON(code, dto.NewApiResponseFailed(msg))
}
