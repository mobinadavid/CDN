package middlewares

import (
	"cdn/src/api/http/response"
	"cdn/src/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ApiKeyMiddleware struct {
	apiKey string
}

func NewApiKeyMiddleware() *ApiKeyMiddleware {
	apiKey := config.GetInstance().Get("APP_API_KEY")
	return &ApiKeyMiddleware{apiKey: apiKey}
}

func (middleware *ApiKeyMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqApiKey := c.GetHeader("X-API-Key")

		if middleware.apiKey == "" {
			response.Api(c).
				SetStatusCode(http.StatusUnauthorized).
				SetMessage("api key not set").
				Send()
			c.Abort()
			return
		}

		if middleware.apiKey != reqApiKey {
			response.Api(c).
				SetStatusCode(http.StatusUnauthorized).
				SetMessage("invalid API key").
				Send()
			c.Abort()
			return
		}

		c.Next()
	}
}
