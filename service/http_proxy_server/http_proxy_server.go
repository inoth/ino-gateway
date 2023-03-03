package httpproxyserver

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	httpproxymiddleware "github/inoth/ino-gateway/http_proxy_middleware"

	"github.com/inoth/ino-toybox/components/config"
	"github.com/inoth/ino-toybox/components/logger"
)

var (
	httpSrvHandler *http.Server
)

type HttpProxyServer struct {
	Middlewares []gin.HandlerFunc
}

func (hps *HttpProxyServer) Start() error {
	fmt.Println("load service proxy")
	r := gin.New()
	r.Use(hps.Middlewares...)
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.String(200, "ok")
	})

	r.Use(
		httpproxymiddleware.HTTPAccessMode(),
		httpproxymiddleware.HTTPFlowCount(),
		httpproxymiddleware.HTTPFlowLimit(),
		httpproxymiddleware.HTTPJwtAuthToken(),
		httpproxymiddleware.HttpJwtFlowCount(),
		httpproxymiddleware.HttpJwtFlowLimit(),
		httpproxymiddleware.HTTPReverseProxy(),
	)

	httpSrvHandler = &http.Server{
		Addr:           config.Cfg.GetString("base.server_port"),
		Handler:        r,
		ReadTimeout:    time.Duration(config.Cfg.GetInt("proxy.http.read_timeout")) * time.Second,
		WriteTimeout:   time.Duration(config.Cfg.GetInt("proxy.http.write_timeout")) * time.Second,
		MaxHeaderBytes: 1 << uint(config.Cfg.GetInt("proxy.http.max_header_bytes")),
	}

	if err := httpSrvHandler.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Zap.Error(err.Error())
		return err
	}
	return nil
}

func (hps *HttpProxyServer) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := httpSrvHandler.Shutdown(ctx); err != nil {
		logger.Zap.Error(err.Error())
	}
}
