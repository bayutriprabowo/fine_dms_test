package middleware

import (
	"net/http"
	"strings"

	"enigmacamp.com/fine_dms/model/dto"
	"enigmacamp.com/fine_dms/utils"
	"github.com/gin-gonic/gin"
)

func ValidateToken(secret []byte) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHdr := ctx.GetHeader("Authorization")
		tokStr := strings.Replace(authHdr, "Bearer ", "", 1)

		user_id, err := utils.ValidateToken(tokStr, secret)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized,
				dto.NewApiResponseFailed("invalid token"))
			ctx.Abort()
			return
		} else {
			ctx.Set("user_id", user_id)
		}

		ctx.Next()
	}
}
