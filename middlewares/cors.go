package middlewares

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type CorsConfig struct {
	AllowOrigins     []string `json:"allowOrigins"`
	AllowMethods     []string `json:"allowMethods"`
	AllowHeaders     []string `json:"allowHeaders"`
	AllowCredentials bool     `json:"allowCredentials"`
	MaxAge           int      `json:"maxAge"`
}

func Cors(conf CorsConfig) gin.HandlerFunc {
	corsConfig := cors.Config{
		AllowMethods:     conf.AllowMethods,
		AllowHeaders:     conf.AllowHeaders,
		AllowCredentials: conf.AllowCredentials,
	}
	if len(conf.AllowOrigins) == 0 {
		corsConfig.AllowAllOrigins = true
	} else {
		corsConfig.AllowOrigins = conf.AllowOrigins
	}
	if conf.MaxAge != 0 {
		corsConfig.MaxAge = time.Duration(conf.MaxAge) * time.Hour
	}
	return cors.New(corsConfig)
}
