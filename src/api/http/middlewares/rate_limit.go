package middlewares

import (
	"cdn/src/api/http/response"
	"cdn/src/config"
	"cdn/src/service/redis"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RateLimit(redisService *redis.RedisService) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.GetHeader("x-forwarded-for")   // Get the client's IP address
		userAgent := c.GetHeader("User-Agent") // Get the User-Agent header
		configs := config.GetInstance()

		// Check rate limit
		allowed, err := redisService.CheckAndIncrementRateLimit(ip, userAgent)
		if err != nil {
			response.Api(c).SetMessage(err.Error()).SetStatusCode(http.StatusInternalServerError).Send()
			c.Abort()
			return
		}
		rateLimitStr := configs.Get("RATE_LIMIT")
		periodStr := configs.Get("RATE_LIMITER_PERIOD_PER_SECOND")
		if !allowed {
			response.Api(c).SetMessage("You can't put more than " + rateLimitStr + " Objects in " + periodStr + " seconds").SetStatusCode(http.StatusTooManyRequests).Send()
			c.Abort()
			return
		}

		c.Next()
	}
}
