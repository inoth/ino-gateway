package httpproxymiddleware

import "github.com/gin-gonic/gin"

/*
	用户角色权限过滤器
*/

func HTTPJwtAuthRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, ok := c.Get("user")
		if !ok {
			c.Next()
			return
		}
		// 处理用户相关权限事宜
		c.Next()
	}
}
