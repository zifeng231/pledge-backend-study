package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pledge-backend-study/api/common/statecode"
	"pledge-backend-study/api/models/response"
	"pledge-backend-study/config"
	"pledge-backend-study/db"
	"pledge-backend-study/utils"
)

func CheckToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		res := response.Gin{Res: c}
		token := c.Request.Header.Get("authCode")

		username, err := utils.ParseToken(token, config.Config.Jwt.SecretKey)
		if err != nil {
			res.Response(c, statecode.TokenErr, nil)
			c.Abort()
			return
		}

		if username != config.Config.DefaultAdmin.Username {
			res.Response(c, statecode.TokenErr, nil)
			c.Abort()
			return
		}

		// Judge whether the user logout
		resByteArr, err := db.RedisGet(username)
		if string(resByteArr) != `"login_ok"` {
			res.Response(c, statecode.TokenErr, nil)
			c.Abort()
			return
		}

		c.Set("username", username)

		c.Next()
	}
}
