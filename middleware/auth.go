package middleware

import (
	"github.com/gin-gonic/gin"
)

func ValidateToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// TODO

		ctx.Next()
	}
}
