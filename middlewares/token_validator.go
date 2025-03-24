package middlewares

import (
	"net/http"
	"strings"

	"github.com/feiyizhou/base/utils"
	"github.com/gin-gonic/gin"
)

// TokenValidate validate user's token is valid or not
func TokenValidate() gin.HandlerFunc {
	return func(context *gin.Context) {
		auth := context.Request.Header.Get("Authorization")
		if len(auth) == 0 {
			context.AbortWithStatusJSON(http.StatusUnauthorized, "Token is empty")
			return
		}
		if len(strings.Fields(auth)) != 2 || !strings.EqualFold(strings.Fields(auth)[0], "Bearer") {
			context.AbortWithStatusJSON(http.StatusUnauthorized, "Token is invalid")
			return
		}
		userID, _, err := utils.ParseToken(strings.Fields(auth)[1])
		if err != nil {
			context.AbortWithStatusJSON(http.StatusUnauthorized, "Token oauth failed")
			return
		}
		context.Set("userID", userID)
		context.Next()
		return
	}
}
