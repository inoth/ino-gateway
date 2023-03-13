package httpproxymiddleware

import (
	"errors"
	flowcount "github/inoth/gateway/util/flow_count"
	"net/http"

	"github.com/inoth/toybox/res"

	"github.com/gin-gonic/gin"
)

func HttpJwtFlowCount() gin.HandlerFunc {
	return func(c *gin.Context) {
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
		// 租户计数器
		tenantCount, err := flowcount.FlowCounterHandler.GetCounter(flowcount.FlowTotalTenant + tenantId)
		if err != nil {
			res.ResultErr(c, http.StatusBadRequest, err)
			c.Abort()
			return
		}
		tenantCount.Increase()

		c.Next()
	}
}
