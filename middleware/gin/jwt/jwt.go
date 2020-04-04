package jwt

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/cuirixin/phoenix_gocore/libs/jwtoken"
)

func Auth(tokenKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get(tokenKey)
		println("auth-token: ", tokenString)
		if tokenString == "" {
			c.AbortWithError(401, errors.New("token required"))
			return
		}
		succ, uid := jwtoken.JWTVerify(tokenString)
		if succ == false {
			c.AbortWithError(401, errors.New("token auth failed"))
			return
		}
		c.Set("uid", uid)
	}
}
