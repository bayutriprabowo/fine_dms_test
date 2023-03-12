package middleware

import (
	"net/http"
	"strings"

	"enigmacamp.com/fine_dms/controller"
	"enigmacamp.com/fine_dms/utils"
	"github.com/gin-gonic/gin"
)

func ValidateToken(secret []byte) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHdr := ctx.GetHeader("Authorization")
		tokStr := strings.Replace(authHdr, "Bearer ", "", 1)

		user_id, err := utils.ValidateToken(tokStr, secret)
		if err != nil {
			controller.FailedJSONResponse(ctx, http.StatusUnauthorized,
				"invalid token")
		} else {
			ctx.Set("user_id", user_id)
		}

		ctx.Next()
	}
}
