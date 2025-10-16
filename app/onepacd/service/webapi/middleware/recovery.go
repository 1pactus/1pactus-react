package middleware

import (
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"

	"github.com/frimin/1pactus-react/log"
	"github.com/gin-gonic/gin"
)

func CustomRecovery(logger log.ILogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				httpRequest, _ := httputil.DumpRequest(c.Request, false)

				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				logger.WithFields(map[string]interface{}{
					"panic":   err,
					"request": string(httpRequest),
					"stack":   string(debug.Stack()),
				}).Error("http server panic")

				if brokenPipe {
					c.Error(err.(error)) //nolint: errcheck
					c.Abort()
					return
				}

				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"code":  http.StatusInternalServerError,
					"error": "Internal Server Error",
				})
			}
		}()
		c.Next()
	}
}
