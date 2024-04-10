package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func IPWhiteList(whiteList []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		allow := false
		for _, whiteIP := range whiteList {
			if whiteIP == ip {
				allow = true
			}
		}
		if !allow {
			c.AbortWithStatusJSON(http.StatusForbidden, "IP is not allowed")
			return
		}
		c.Next()
	}
}
