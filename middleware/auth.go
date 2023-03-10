package middleware

import (
	"net/http"
	"strings"

	"enigmacamp.com/fine_dms/controller"
	"enigmacamp.com/fine_dms/model/dto"
	"enigmacamp.com/fine_dms/utils"
	"github.com/gin-gonic/gin"
)

func ValidateToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHdr := ctx.GetHeader("Authorization")
		tokStr := strings.Replace(authHdr, "Bearer ", "", 1)
		tok, err := utils.ValidateToken(tokStr)
		if err != nil {
			controller.FailedJSONResponse(ctx,
				http.StatusUnauthorized,
				dto.NewApiResponseFailed("invalid token"),
			)
		} else {
			ctx.Set("user_id", tok)
		}

		ctx.Next()
	}
}
