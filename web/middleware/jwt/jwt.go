package jwt

import (
	"miner_proxy/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		var data interface{}

		code = 200
		token := c.Request.Header.Get("token")
		if token == "" {
			code = 301
		} else {
			claims, err := utils.ParseToken(token)
			if err != nil {
				code = 302
			} else if time.Now().Unix() > claims.ExpiresAt {
				code = 303
			}
		}

		if code != 200 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": code,
				"msg":  "Token 不正确",
				"data": data,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
