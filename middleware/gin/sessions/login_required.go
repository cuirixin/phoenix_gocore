package sessions

import (
	"github.com/gin-gonic/gin"
	"github.com/cuirixin/phoenix_gocore/middleware/gin/secure_cookie"
	"net/http"
)

const (
	COOKIE_MAX_AGE = 1999999999
	COOKIE_DOMAIN = ".phoenix.com" // 可以访问该Cookie的域名。如果设置为“.google.com”，则所有以“google.com”结尾的域名都可以访问该Cookie。注意第一个字符必须为“.”
	COOKIE_PATH = "/"
)


// set secure cookie user_token
func AuthLogin(c *gin.Context, uid string)  {
	secure_cookie.SetSecureCookie(
		c,
		"user_token",
		uid,
		COOKIE_MAX_AGE,
		COOKIE_PATH,
		COOKIE_DOMAIN,
		true,true)
}

// delete cookie user_token
func AuthLogout(c *gin.Context)  {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "user_token",
		Value:    "",
		MaxAge:   -1,
		Path:     COOKIE_PATH,
		Domain:   COOKIE_DOMAIN,
		Secure:   true,
		HttpOnly: true,
	})
}

// Login Require Decorator
func LoginRequired(handle gin.HandlerFunc) gin.HandlerFunc {

	return func(c *gin.Context) {
		userToken, cookie_err := secure_cookie.GetSecureCookie(c,"user_token",1)

		var is_login  bool = true

		if cookie_err != nil{
			is_login = false
		}

		//Tudo 添加查数据库逻辑

		if is_login == false{
			c.JSON(http.StatusUnauthorized,
				gin.H{
					"code":  -2,
					"message": "login requierd",
				})
		}else {
			handle(c)
			c.Set("currentUserId",userToken)
			c.Set("currentUser", userToken)
		}
	}
}