package proxyconfserver

import (
	"context"
	"github/inoth/ino-gateway/model"
	"github/inoth/ino-gateway/model/request"
	"net/http"
	"time"

	"github.com/inoth/ino-toybox/components/config"
	"github.com/inoth/ino-toybox/components/logger"
	"github.com/inoth/ino-toybox/res"

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
		Addr:           ":9001",
		Handler:        r,
		ReadTimeout:    time.Duration(config.Cfg.GetInt("proxy.http.read_timeout")) * time.Second,
		WriteTimeout:   time.Duration(config.Cfg.GetInt("proxy.http.write_timeout")) * time.Second,
		MaxHeaderBytes: 1 << uint(config.Cfg.GetInt("proxy.http.max_header_bytes")),
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
	req, ok := request.RequestJsonParamHandler[[]request.ServiceNodeRequests](c)
	if !ok {
		return
	}
	adds := make([]*model.ServiceInfo, 0, len(req))
	for _, val := range req {
		tmp := &model.ServiceInfo{
			ServiceKey:  val.ServiceKey,
			Version:     val.Version,
			Desc:        val.Desc,
			NeedAuth:    val.NeedAuth,
			NeedLicense: val.NeedLicense,
		}
		tmp.Hosts = append(tmp.Hosts, model.ServerNode{Host: val.Host})
		adds = append(adds, tmp)
	}
	err := servicemanage.ServiceManage.AppendService(adds...)
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
