package main

import (
	"github/inoth/ino-gateway/components/cache"
	"github/inoth/ino-gateway/components/config"
	"github/inoth/ino-gateway/components/local"
	"github/inoth/ino-gateway/components/logger"
	servicemanage "github/inoth/ino-gateway/components/service_manage"
	"github/inoth/ino-gateway/middleware"
	"github/inoth/ino-gateway/register"
	httpproxyserver "github/inoth/ino-gateway/service/http_proxy_server"
	proxyconfserver "github/inoth/ino-gateway/service/proxy_conf_server"
	"os"

	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

type UserInfo struct {
	Name string
}

func main() {
	reg := register.NewRegister(
		&local.CacheComponent{},
		&config.ViperComponent{},
		&logger.ZapComponent{},
		&cache.RedisComponent{},
		&servicemanage.ServiceManager{},
	).Init().SubStart(
		&proxyconfserver.HttpTestProxyServer{},  // 服务A
		&proxyconfserver.HttpTestAProxyServer{}, // 服务B
		&httpproxyserver.HttpProxyServer{
			Middlewares: []gin.HandlerFunc{
				middleware.RecoveryMiddleware(),
				middleware.RequestLog(),
			}}, // 代理
	)

	// 监听退出信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	reg.Stop()
}
