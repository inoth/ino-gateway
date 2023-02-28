package httpproxymiddleware

import "github.com/gin-gonic/gin"

/*
	获取用户权限信息
*/

func HTTPJwtAuthRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, ok := c.Get("user")
		if !ok {
			c.Next()
		}
		// 处理用户相关权限事宜
	}
}
