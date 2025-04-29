package middlewares

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/gin-gonic/gin"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		body, _ := io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		c.Next()
		latency := time.Since(start)
		fmt.Printf("\n[GIN] %v | %3d | %13v | %-7s %s\n",
			time.Now().Format("2006/01/02 - 15:04:05"),
			c.Writer.Status(),
			latency,
			c.Request.Method,
			c.Request.URL.Path)
		fmt.Println("[Request Header]")
		for k, v := range c.Request.Header {
			fmt.Printf("%s: %v\n", k, v)
		}
		if len(c.Request.URL.Query()) != 0 {
			fmt.Println("[Request Param]")
			fmt.Println(c.Request.URL.Query())
		}
		if len(body) != 0 {
			fmt.Println("[Request Body]")
			fmt.Println(string(body))
		}
	}
}
