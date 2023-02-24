package httpproxyserver

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github/inoth/ino-gateway/components/logger"
	httpproxymiddleware "github/inoth/ino-gateway/http_proxy_middleware"
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
		httpproxymiddleware.HTTPAccessModeMiddleware(),
		httpproxymiddleware.HTTPJwtAuthTokenMiddleware(),
		httpproxymiddleware.HTTPReverseProxyMiddleware(),
	)

	httpSrvHandler = &http.Server{
		Addr:    ":80",
		Handler: r,
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