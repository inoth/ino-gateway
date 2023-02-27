package httpproxymiddleware

import (
	"errors"
	"github/inoth/ino-gateway/model"
	"github/inoth/ino-gateway/res"
	"github/inoth/ino-gateway/util/auth"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

/*
	获取客户端 jwt 信息，解析
*/

func HTTPJwtAuthTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		service, ok := c.Get("service")
		if !ok {
			res.ResultErr(c, http.StatusBadGateway, errors.New("service not found"))
			c.Abort()
			return
		}
		serviceInfo := service.(*model.ServiceInfo)
		// 该服务是否需要登录信息
		if serviceInfo.NeedAuth {
			token := strings.ReplaceAll(c.GetHeader("Authorization"), "Bearer ", "")
			if len(token) <= 0 {
				var err error
				token, err = c.Cookie("Authorization")
				if err != nil {
					res.ResultErr(c, res.InvalidRequestErrorCode, errors.New("session not found"))
					c.Abort()
					return
				}
			}
			if len(token) <= 0 {
				res.ResultErr(c, res.InvalidRequestErrorCode, errors.New("session not found"))
				c.Abort()
				return
			}
			customerInfo, err := auth.ParseToken(token)
			if err != nil {
				res.ResultErr(c, res.InvalidRequestErrorCode, err)
				c.Abort()
				return
			}
			for key, val := range customerInfo.UserInfo {
				c.Request.Header.Set(key, val.(string))
			}
		}
		c.Next()
	}
}
