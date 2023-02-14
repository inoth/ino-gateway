package main

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

// NewProxy takes target host and creates a reverse proxy
func NewProxy(targetHost string) (*httputil.ReverseProxy, error) {
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

	// proxy.ModifyResponse = modifyResponse()
	// proxy.ErrorHandler = errorHandler()
	return proxy, nil
}

func modifyRequest(req *http.Request) {
	req.Header.Set("X-Proxy", "Simple-Reverse-Proxy")
}

func errorHandler() func(http.ResponseWriter, *http.Request, error) {
	return func(w http.ResponseWriter, req *http.Request, err error) {
		fmt.Printf("Got error while modifying response: %v \n", err)
		return
	}
}

func modifyResponse() func(*http.Response) error {
	return func(resp *http.Response) error {
		return errors.New("response body is invalid")
	}
}

// ProxyRequestHandler handles the http request using proxy
func ProxyRequestHandler(proxy *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	}
}

var proxyList = []string{"http://localhost:8887", "http://localhost:8888", "http://localhost:8889"}
var server *http.Server

func main() {

	go func() {
		r := gin.New()
		r.POST("/job/v1", func(c *gin.Context) {
			c.String(200, "machine post ok")
		})
		r.GET("/job/v1", func(c *gin.Context) {
			c.String(200, "machine get ok")
		})
		r.Run(":8887")
	}()
	go func() {
		r := gin.New()
		r.POST("/job/v1", func(c *gin.Context) {
			c.String(200, "notify post ok")
		})
		r.GET("/job/v1", func(c *gin.Context) {
			c.String(200, "notify get ok")
		})
		r.Run(":8888")
	}()
	go func() {
		r := gin.New()
		r.POST("/job/v1", func(c *gin.Context) {
			c.String(200, "job post ok")
		})
		r.GET("/job/v1", func(c *gin.Context) {
			c.String(200, "job get ok")
		})
		r.Run(":8889")
	}()

	proxyMid := func() gin.HandlerFunc {
		return func(c *gin.Context) {
			proxy, err := NewProxy(proxyList[rand.Intn(2)])
			if err != nil {
				c.String(http.StatusBadGateway, "Bad Gateway")
				c.Abort()
			}
			// ProxyRequestHandler(proxy)
			proxy.ServeHTTP(c.Writer, c.Request)
			c.Abort()
		}
	}

	// initialize a reverse proxy and pass the actual backend server url here
	// proxy, err := NewProxy(proxyList[0])
	// if err != nil {
	// 	panic(err)
	// }

	// // handle all requests to your server using the proxy
	// http.HandleFunc("/", ProxyRequestHandler(proxy))
	// log.Fatal(http.ListenAndServe(":8080", nil))

	go func() {
		r2 := gin.New()

		r2.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})

		r2.Use(proxyMid())

		server = &http.Server{
			Addr:    ":8080",
			Handler: r2,
		}
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf(" [ERROR] http_proxy_run err:%v\n", err)
		}
	}()
	<-make(chan struct{})
}
