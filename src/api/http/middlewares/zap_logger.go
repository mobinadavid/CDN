package middlewares

import (
	"cdn/src/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

func ZapLogger(c *gin.Context) {
	// Start timer
	start := time.Now()

	// Process request
	c.Next()

	// Calculate duration
	duration := time.Since(start)

	// Log request details
	logger.GetInstance().Info("Incoming request",
		zap.String("type", "request"),
		zap.String("method", c.Request.Method),
		zap.String("url", c.Request.URL.String()),
		zap.Int("status_code", c.Writer.Status()),
		zap.Duration("response_time", duration),
		zap.String("client_ip", c.ClientIP()),
		zap.String("user_agent", c.Request.UserAgent()),
		zap.String("uuid", c.GetString("request-uuid")),
		// add body,header,query param
	)
}
