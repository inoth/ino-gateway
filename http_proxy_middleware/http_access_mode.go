package httpproxymiddleware

import (
	servicemanage "github/inoth/ino-gateway/components/service_manage"
	"github/inoth/ino-gateway/res"

	"github.com/gin-gonic/gin"
)

// 根据路由约定获取 /[service]/[version]
func HTTPAccessModeMiddleware() gin.HandlerFunc {
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
