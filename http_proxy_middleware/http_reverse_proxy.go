package httpproxymiddleware

import (
	"errors"
	"fmt"
	servicemanage "github/inoth/gateway/components/service_manage"
	"github/inoth/gateway/model"
	flowcount "github/inoth/gateway/util/flow_count"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/inoth/toybox/components/config"
	"github.com/inoth/toybox/res"

	"github.com/gin-gonic/gin"
)

// 获取服务key
// 拿到服务注册信息
// 负载均衡选择代理目标
// 初始化反向代理
func HTTPReverseProxy() gin.HandlerFunc {
	return func(c *gin.Context) {
		service, ok := c.Get("service")
		if !ok {
			res.ResultErr(c, http.StatusBadGateway, errors.New("service not found"))
			c.Abort()
			return
		}
		serviceInfo := service.(*model.ServiceInfo)

		proxy, err := newProxy(serviceInfo.GetHost())
		if err != nil {
			res.ResultErr(c, http.StatusBadGateway, errors.New("bad status bad gateway"))
			c.Abort()
			return
		}
		proxy.ServeHTTP(c.Writer, c.Request)
		c.Abort()
		return
	}
}

// NewProxy takes target host and creates a reverse proxy
func newProxy(targetHost string) (*httputil.ReverseProxy, error) {
	url, err := url.Parse(targetHost)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(url)

	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		modifyRequest(req)
	}

	// proxy.ModifyResponse = modifyResponse
	// proxy.ErrorHandler = errorHandler
	return proxy, nil
}

func modifyRequest(req *http.Request) {
	req.Header.Set("X-Proxy", "Cmdb-Reverse-Proxy")
}

// 错误端口计数器
func errRequestFlowCount(serverName, version, host string) error {
	hostCount, err := flowcount.FlowCounterHandler.GetCounter(flowcount.FlowTotalErrHost + host)
	if err != nil {
		return err
	}
	// 每次错误新增一次记录
	hostCount.Increase()
	// TODO:配置文件中设定每日允许最大错误量
	count, _ := hostCount.GetHourData(time.Now())
	if count > int64(config.Cfg.GetInt("proxy.error_von_hour")) {
		// 删除当前服务端口, 等待服务重启后重新注册
		fmt.Printf("服务：%v/%v[%v]累计错误数超出设定阈值:%d\n", serverName, version, host, count)
		return servicemanage.ServiceManage.DelService(serverName, version, model.ServerNode{Host: host})
	}
	return nil
}

// func modifyResponse(resp *http.Response) error {
// 	// 获取响应的 body
// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return err
// 	}

// 	// 修改响应的 body
// 	newBody := []byte("Modified Response: " + string(body))
// 	resp.Body = io.NopCloser(bytes.NewBuffer(newBody))
// 	resp.ContentLength = int64(len(newBody))
// 	resp.Header.Set("Content-Length", fmt.Sprint(len(newBody)))

//	// TODO:插入缓存的数据？

// 	return nil
// }
