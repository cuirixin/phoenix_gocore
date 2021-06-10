package log

import (
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

func GinLogger() gin.HandlerFunc {

	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		if raw != "" {
			path = path + "?" + raw
		}

		// log
		Info("[GIN]",
			Int("status", c.Writer.Status()),
			Duration("latency", time.Now().Sub(start)),
			String("ip", c.ClientIP()),
			String("method", c.Request.Method),
			Int64("content_length", c.Request.ContentLength),
			Int("size", c.Writer.Size()),
			String("path", path),
			String("error_message", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			String("user-agent", c.Request.UserAgent()))
	}
}

func GinRecoveryLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				path := c.Request.URL.Path
				raw := c.Request.URL.RawQuery
				if raw != "" {
					path = path + "?" + raw
				}
				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					Error("[Recovery]",
						String("error", err.(string)),
						String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				Error("[Recovery]",
					String("error", err.(string)),
					String("request", string(httpRequest)),
					String("stack", string(debug.Stack())),
				)
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
