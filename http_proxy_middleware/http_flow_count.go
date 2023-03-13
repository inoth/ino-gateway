package httpproxymiddleware

import (
	"errors"
	"github/inoth/gateway/model"
	flowcount "github/inoth/gateway/util/flow_count"
	"net/http"

	"github.com/inoth/toybox/res"

	"github.com/gin-gonic/gin"
)

/*
	所有模块总计数统计
	对服务模块请求计数统计
*/

func HTTPFlowCount() gin.HandlerFunc {
	return func(c *gin.Context) {
		service, ok := c.Get("service")
		if !ok {
			res.ResultErr(c, http.StatusBadGateway, errors.New("service not found"))
			c.Abort()
			return
		}
		serviceInfo := service.(*model.ServiceInfo)

		// 所有服务请求计数
		totalCount, err := flowcount.FlowCounterHandler.GetCounter(flowcount.FlowTotal)
		if err != nil {
			res.ResultErr(c, http.StatusBadRequest, err)
			c.Abort()
			return
		}
		totalCount.Increase()

		// 单个服务计数
		serviceTotalCount, err := flowcount.FlowCounterHandler.GetCounter(flowcount.FlowTotalService + serviceInfo.ServiceKey)
		if err != nil {
			res.ResultErr(c, http.StatusBadRequest, err)
			c.Abort()
			return
		}
		serviceTotalCount.Increase()

		c.Next()
	}
}
