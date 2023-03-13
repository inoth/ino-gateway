package httpproxymiddleware

import (
	"errors"
	"github/inoth/gateway/model"
	"net/http"
	"regexp"
	"strings"

	jwtauth "github.com/inoth/toybox/utils/jwt_auth"

	"github.com/inoth/toybox/res"

	"github.com/gin-gonic/gin"
)

/*
	获取客户端 jwt 信息，解析
	后续步骤：
	租户服务请求计数统计
	角色权限
	租户license
*/

func HTTPJwtAuthToken() gin.HandlerFunc {
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

			// 判断是否为login相关接口
			r := regexp.MustCompile("login|Login|registry")
			if r.MatchString(c.Request.URL.Path) {
				c.Next()
				return
			}

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
			if len(serviceInfo.JwtSignKey) <= 0 {
				serviceInfo.JwtSignKey = jwtauth.DEFAULT_SIGNKEY
			}
			customerInfo, err := jwtauth.ParseToken(serviceInfo.JwtSignKey, token)
			if err != nil {
				res.ResultErr(c, res.InvalidRequestErrorCode, err)
				c.Abort()
				return
			}
			c.Set("user", customerInfo.UserInfo)
			for key, val := range customerInfo.UserInfo {
				c.Request.Header.Set(key, val.(string))
			}
		}
		c.Next()
	}
}
