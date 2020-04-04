package jwt

import (
	//"errors"
	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/cuirixin/phoenix_gocore/libs/jwtoken"
)

func Auth(tokenKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get(tokenKey)
		println("auth-token: ", tokenString)
		if tokenString == "" {
			// c.AbortWithError(http.StatusBadRequest, errors.New("token required"))

			c.JSON(http.StatusBadRequest, gin.H{
				"code": http.StatusBadRequest,
				"msg": "token required",
			})
			c.Abort()
			return
		}
		succ, uid := jwtoken.JWTVerify(tokenString)
		if succ == false {
			// c.AbortWithError(401, errors.New("token auth failed"))
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": http.StatusBadRequest,
				"msg": "token auth failed",
			})
			c.Abort()
			return
		}
		c.Set("uid", uid)
	}
}
