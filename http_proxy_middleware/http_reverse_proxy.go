package httpproxymiddleware

import (
	"errors"
	"github/inoth/ino-gateway/model"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/inoth/ino-toybox/res"

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
