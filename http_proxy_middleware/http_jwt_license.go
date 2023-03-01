package httpproxymiddleware

import "github.com/gin-gonic/gin"

/*
	获取租户对应服务模块使用限制
*/

func HttpJwtLicense() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
