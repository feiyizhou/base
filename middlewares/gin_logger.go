package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GinLogger 接收gin框架默认的日志
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		// 挂起当前中间件，执行下一个中间件
		c.Next()
		cost := time.Since(start)
		zap.L().Info(path,
			zap.Int("status", c.Writer.Status()),                                 // 状态码 eg: 200
			zap.String("method", c.Request.Method),                               // 请求方法类型 eg: GET
			zap.String("path", path),                                             // 请求路径 eg: /test
			zap.String("query", query),                                           // 请求参数 eg: name=1&password=2
			zap.String("ip", c.ClientIP()),                                       // 返回真实的客户端IP eg: ::1（这个就是本机IP，ipv6地址）
			zap.String("user-agent", c.Request.UserAgent()),                      // 返回客户端的用户代理。 eg: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.0.0 Safari/537.36
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()), // 返回Errors 切片中ErrorTypePrivate类型的错误
			zap.Duration("cost", cost),                                           // 返回花费时间
		)
	}
}
