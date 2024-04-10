package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	libredis "github.com/redis/go-redis/v9"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	sredis "github.com/ulule/limiter/v3/drivers/store/redis"
)

/*
	5 reqs/second: "5-S"
	10 reqs/minute: "10-M"
	1000 reqs/hour: "1000-H"
	2000 reqs/day: "2000-D"
*/

func RateLimiter(client *libredis.Client, limitStr string) gin.HandlerFunc {
	return func(context *gin.Context) {
		var (
			rate  limiter.Rate
			err   error
			store limiter.Store
		)
		if len(limitStr) == 0 {
			// 1000-H 每小时1000次
			limitStr = "1000-H"
		} else {
			parts := strings.Split(limitStr, "-")
			if len(parts) != 2 {
				context.AbortWithStatusJSON(http.StatusInternalServerError, "Init rate limiter instance failed")
				return
			}
		}
		rate, err = limiter.NewRateFromFormatted(limitStr)
		if err != nil {
			context.AbortWithStatusJSON(http.StatusInternalServerError, "Init rate limiter instance failed")
			return
		}
		prefix := ""
		if userID, ok := context.Get("userID"); ok {
			prefix = fmt.Sprintf("%s_%s_", context.FullPath(), userID)
		} else {
			prefix = fmt.Sprintf("%s_", context.FullPath())
		}
		store, err = sredis.NewStoreWithOptions(client, limiter.StoreOptions{
			Prefix: prefix,
		})
		if err != nil {
			context.AbortWithStatusJSON(http.StatusInternalServerError, "Init rate limiter instance failed")
			return
		}
		middleware := &mgin.Middleware{
			Limiter:     limiter.New(store, rate),
			KeyGetter:   mgin.DefaultKeyGetter,
			ExcludedKey: nil,
		}
		key := middleware.KeyGetter(context)
		if middleware.ExcludedKey != nil && middleware.ExcludedKey(key) {
			context.Next()
			return
		}
		ctx, err := middleware.Limiter.Get(context, key)
		if err != nil {
			context.AbortWithStatusJSON(http.StatusInternalServerError, "Get the client ip failed")
			return
		}
		if ctx.Reached {
			context.AbortWithStatusJSON(http.StatusTooManyRequests, "Reached the rate limit of request")
			return
		}
		context.Next()
	}
}
