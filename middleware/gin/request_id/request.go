/**
 * @Author: victor
 * @Description:
 * @File:  request
 * @Version: 1.0.0
 * @Date: 2020/4/4 9:53 上午
 */

package request_id

import "github.com/gin-gonic/gin"
import "github.com/satori/go.uuid"

func RequestIdMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("X-Request-Id", uuid.NewV4().String())
		c.Next()
		// c.Writer.Header().Get("X-Request-Id")
	}
}