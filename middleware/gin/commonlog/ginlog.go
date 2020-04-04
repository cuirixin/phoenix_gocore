/**
 * @Author: victor
 * @Description:
 * @File:  ginlog
 * @Version: 1.0.0
 * @Date: 2020/4/3 3:45 下午
 */

package commonlog

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

func LoggerWithFormatter() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			//param.Request.UserAgent(),
			param.ErrorMessage,
		)
	})
}