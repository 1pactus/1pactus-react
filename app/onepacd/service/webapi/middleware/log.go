package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func Log(logger zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		end := time.Now()
		latency := end.Sub(start)
		if raw != "" {
			path = path + "?" + raw
		}

		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()
		bodySize := c.Writer.Size()

		event := logger.Info()
		if statusCode >= 500 {
			event = logger.Error()
		} else if statusCode >= 400 {
			event = logger.Warn()
		} else {
			event = logger.Info()
		}

		event.
			Str("client_ip", clientIP).
			Str("method", method).
			Int("status_code", statusCode).
			Int("body_size", bodySize).
			Str("path", path).
			Str("latency", latency.String()).
			Str("error_message", errorMessage).
			Msg("request")
	}
}
