package proxyconfserver

import (
	"context"
	"github/inoth/ino-gateway/components/logger"
	"github/inoth/ino-gateway/model"
	"github/inoth/ino-gateway/model/request"
	"github/inoth/ino-gateway/res"
	"net/http"
	"time"

	servicemanage "github/inoth/ino-gateway/components/service_manage"

	"github.com/gin-gonic/gin"
)

var (
	ProxySrvHandler      *http.Server
	ProxyTestSrvHandler  *http.Server
	ProxyTestASrvHandler *http.Server
)

type HttpProxyServer struct {
	Middlewares []gin.HandlerFunc
}

func (hps *HttpProxyServer) Start() error {
	r := gin.New()
	r.Use(hps.Middlewares...)

	proxy := r.Group("/proxy")
	{
		proxy.GET("/", queryProxyConfig)
		proxy.POST("/add", addProxyConfig)
		proxy.DELETE("/remove", removeProxyConfig)
	}

	ProxySrvHandler = &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	if err := ProxySrvHandler.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Zap.Error(err.Error())
		return err
	}
	return nil
}

func (hps *HttpProxyServer) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := ProxySrvHandler.Shutdown(ctx); err != nil {
		logger.Zap.Error(err.Error())
	}
}

// 添加服务节点
func addProxyConfig(c *gin.Context) {
	req, ok := request.RequestJsonParamHandler[request.ServiceNodeRequests](c)
	if !ok {
		return
	}
	err := servicemanage.ServiceManage.AppendService(&model.ServiceInfo{
		ServiceKey:  req.ServiceKey,
		Version:     req.Version,
		Desc:        req.Desc,
		Hosts:       req.Hosts,
		NeedAuth:    req.NeedAuth,
		NeedLicense: req.NeedLicense,
	})
	if err != nil {
		res.ResultErr(c, res.InvalidRequestErrorCode, err)
		return
	}
	res.ResultOk(c, res.SuccessCode)
}

// 删除节点
func removeProxyConfig(c *gin.Context) {
	req, ok := request.RequestJsonParamHandler[request.ServiceNodeRemoveRequest](c)
	if !ok {
		return
	}
	err := servicemanage.ServiceManage.DelService(req.ServiceKey, req.Version, req.Hosts...)
	if err != nil {
		res.ResultErr(c, res.ParamErrorCode, err)
		return
	}
	res.ResultOk(c, res.SuccessCode)
}

func queryProxyConfig(c *gin.Context) {
	list := servicemanage.ServiceManage.GetServiceList()
	res.ResultOk(c, res.SuccessCode, list)
}

// type HttpTestProxyServer struct {
// 	Middlewares []gin.HandlerFunc
// }

// func (htps *HttpTestProxyServer) Start() error {
// 	fmt.Println("load service A")
// 	r := gin.New()
// 	r.Use(htps.Middlewares...)

// 	r.GET("/job/v1", func(ctx *gin.Context) {
// 		v := ctx.Query("v")
// 		user := map[string]interface{}{
// 			"name":   "inoth",
// 			"app_id": "123123123",
// 		}
// 		token, err := auth.CreateToken(user, int64(time.Hour*24))
// 		if err != nil {
// 			res.ResultErr(ctx, 400, err)
// 			return
// 		}
// 		// ctx.SetCookie("Authorization", token, 0, "", "localhost", true, true)
// 		ctx.Header("Authorization", token)

// 		fmt.Println(ctx.Request.Header.Get("name"))
// 		fmt.Println(ctx.Request.Header.Get("app_id"))
// 		res.ResultOk(ctx, 200, "this's service a say: ok "+v)
// 	})

// 	ProxyTestSrvHandler = &http.Server{
// 		Addr:    ":8081",
// 		Handler: r,
// 	}

// 	if err := ProxyTestSrvHandler.ListenAndServe(); err != nil && err != http.ErrServerClosed {
// 		logger.Zap.Error(err.Error())
// 		return err
// 	}
// 	return nil
// }

// func (htps *HttpTestProxyServer) Stop() {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()
// 	if err := ProxyTestSrvHandler.Shutdown(ctx); err != nil {
// 		logger.Zap.Error(err.Error())
// 	}
// }

// type HttpTestAProxyServer struct {
// 	Middlewares []gin.HandlerFunc
// }

// func (htps *HttpTestAProxyServer) Start() error {
// 	fmt.Println("load service B")
// 	r := gin.New()
// 	r.Use(htps.Middlewares...)

// 	r.GET("/cmdb/v1", func(ctx *gin.Context) {
// 		v := ctx.Query("v")
// 		res.ResultOk(ctx, 200, "this's service b say: ok "+v)
// 	})

// 	ProxyTestASrvHandler = &http.Server{
// 		Addr:    ":8082",
// 		Handler: r,
// 	}

// 	if err := ProxyTestASrvHandler.ListenAndServe(); err != nil && err != http.ErrServerClosed {
// 		logger.Zap.Error(err.Error())
// 		return err
// 	}
// 	return nil
// }

// func (htps *HttpTestAProxyServer) Stop() {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()
// 	if err := ProxyTestASrvHandler.Shutdown(ctx); err != nil {
// 		logger.Zap.Error(err.Error())
// 	}
// }
