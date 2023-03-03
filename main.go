package main

import (
	servicemanage "github/inoth/ino-gateway/components/service_manage"

	httpproxyserver "github/inoth/ino-gateway/service/http_proxy_server"
	proxyconfserver "github/inoth/ino-gateway/service/proxy_conf_server"
	"os"

	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"github.com/inoth/ino-toybox/components/cache"
	"github.com/inoth/ino-toybox/components/config"
	"github.com/inoth/ino-toybox/components/logger"
	"github.com/inoth/ino-toybox/components/redis"
	"github.com/inoth/ino-toybox/middleware"
	"github.com/inoth/ino-toybox/register"
)

type UserInfo struct {
	Name string
}

func main() {
	reg := register.NewRegister(
		&cache.CacheComponent{},
		&config.ViperComponent{},
		&logger.ZapComponent{},
		&redis.RedisComponent{},
		&servicemanage.ServiceManager{},
	).Init().SubStart(
		&proxyconfserver.HttpProxyServer{}, // 网关配置接口服务
		&httpproxyserver.HttpProxyServer{
			Middlewares: []gin.HandlerFunc{
				middleware.Recovery(),
				middleware.RequestLog(),
				middleware.Cors(),
			}}, // 代理服务
	)

	// 监听退出信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	reg.Stop() // 关闭注册器内所有服务
}
