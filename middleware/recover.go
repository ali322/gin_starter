package middleware

import (
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Recovery(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						seErr := strings.ToLower(se.Error())
						if strings.Contains(seErr, "broken pipe") ||
							strings.Contains(seErr, "connection reset by peer") {
							brokenPipe = true
						}
					}
				}
				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					logger.Error(c.Request.URL.Path, zap.Any("error", err))
					c.Error(err.(error))
					c.Abort()
					return
				}
				logger.Error("[Recovery from panic]",
					zap.Time("time", time.Now()),
					zap.Any("error", err),
					zap.String("request", string(httpRequest)),
				)
				c.JSON(http.StatusOK, &gin.H{
					"code": -1, "msg": err,
				})
			}
		}()
		c.Next()
	}
}
