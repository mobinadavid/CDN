package middlewares

import (
	"cdn/src/api/http/response"
	"cdn/src/config"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type IPMiddleware struct {
	allowOrigin []string
}

func NewIPMiddleware() *IPMiddleware {
	allowOriginStr := config.GetInstance().Get("ALLOW_ORIGINS")
	allowOrigins := strings.Split(allowOriginStr, ",")

	return &IPMiddleware{
		allowOrigin: allowOrigins,
	}
}

// Middleware checks if the client's IP is in the allowed list
func (middleware *IPMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.Request.RemoteAddr
		if forwarded := c.GetHeader("x-forwarded-for"); forwarded != "" {
			ip = strings.Split(forwarded, ",")[0]
		}

		allowed := false
		for _, allowedIP := range middleware.allowOrigin {
			if ip == allowedIP {
				allowed = true
				break
			}
		}

		if !allowed {
			response.Api(c).
				SetStatusCode(http.StatusForbidden).
				SetMessage("forbidden: Your IP is not allowed").
				Send()
			c.Abort()
			return
		}

		c.Next()
	}
}
