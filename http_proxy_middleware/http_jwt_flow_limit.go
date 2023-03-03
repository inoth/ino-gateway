package httpproxymiddleware

import (
	"errors"
	"fmt"
	flowcount "github/inoth/ino-gateway/util/flow_count"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/inoth/ino-toybox/res"
)

func HttpJwtFlowLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取租户信息
		user, ok := c.Get("user")
		if !ok {
			// 不存在用户信息直接跳过
			c.Next()
			return
		}
		userInfo := user.(map[string]interface{})

		var tenantId string
		t, ok := userInfo["tenant_id"]
		if !ok {
			res.ResultErr(c, http.StatusUnauthorized, errors.New("invalid tenant"))
			c.Abort()
			return
		}
		tenantId = t.(string)

		// 获取租户计数器
		tenantCount, err := flowcount.FlowCounterHandler.GetCounter(flowcount.FlowTotalTenant + tenantId)
		if err != nil {
			res.ResultErr(c, http.StatusForbidden, err)
			c.Abort()
			return
		}

		if tenantCount.TotalCount > int64(100000) {
			res.ResultErr(c, http.StatusForbidden, fmt.Errorf("租户日请求量限流 limit:%v current:%v", 100000, tenantCount.TotalCount))
			c.Abort()
			return
		}
		c.Next()
		return
	}
}
