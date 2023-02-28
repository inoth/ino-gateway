package httpproxymiddleware

import (
	"errors"
	"fmt"
	"github/inoth/ino-gateway/model"
	"github/inoth/ino-gateway/res"
	flowcount "github/inoth/ino-gateway/util/flow_count"
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
	对服务模块进行请求上限限制
*/

func HTTPFlowLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		service, ok := c.Get("service")
		if !ok {
			res.ResultErr(c, http.StatusBadGateway, errors.New("service not found"))
			c.Abort()
			return
		}
		serviceInfo := service.(*model.ServiceInfo)

		// 单个服务请求流量限制
		if serviceInfo.MaxQps > 0 {
			serviceLimiter, err := flowcount.FlowLimiterHandler.GetLimiter(flowcount.FlowTotalService, float64(serviceInfo.MaxQps))
			if err != nil {
				res.ResultErr(c, http.StatusBadGateway, err)
				c.Abort()
				return
			}
			if !serviceLimiter.Allow() {
				res.ResultErr(c, http.StatusBadGateway, fmt.Errorf("service flow limit %v", serviceInfo.MaxQps))
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
