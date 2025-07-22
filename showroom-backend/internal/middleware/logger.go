package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggerMiddleware creates a logger middleware with custom format
func LoggerMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	})
}

// ErrorHandlerMiddleware handles panics and errors
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return gin.Recovery()
}