package httpproxymiddleware

import (
	servicemanage "github/inoth/ino-gateway/components/service_manage"

	"github.com/inoth/ino-toybox/res"

	"github.com/gin-gonic/gin"
)

/*
	获取服务对应注册信息
*/

// 根据路由约定获取 /[service]/[version]
func HTTPAccessMode() gin.HandlerFunc {
	return func(c *gin.Context) {
		service, err := servicemanage.ServiceManage.HTTPAccessMode(c)
		if err != nil {
			res.ResultErr(c, res.InvalidRequestErrorCode, err)
			c.Abort()
			return
		}
		c.Set("service", service)
		c.Next()
	}
}
