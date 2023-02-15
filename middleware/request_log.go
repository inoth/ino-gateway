package middleware

import (
	"bytes"
	"github/inoth/ino-gateway/util"
	"io"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func RequestInLog(c *gin.Context) {
	c.Set("startExecTime", time.Now())
	traceId := util.GetTraceId()
	c.Set("trace_id", traceId)

	bodyBytes, _ := io.ReadAll(c.Request.Body)
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	// TODO: 请求日志等日志库完善后添加
	req := map[string]interface{}{
		"trace_id": traceId,
		"uri":      c.Request.RequestURI,
		"method":   c.Request.Method,
		"args":     c.Request.PostForm,
		"body":     string(bodyBytes),
		"from":     c.ClientIP(),
	}
	log.Fatalln(req)
}

func RequestOutLog(c *gin.Context) {
	endExecTime := time.Now()
	traceId, _ := c.Get("trace_id")
	response, _ := c.Get("response")
	st, _ := c.Get("startExecTime")
	startExecTime, _ := st.(time.Time)
	resp := map[string]interface{}{
		"trace_id":  traceId,
		"uri":       c.Request.RequestURI,
		"method":    c.Request.Method,
		"args":      c.Request.PostForm,
		"from":      c.ClientIP(),
		"response":  response,
		"proc_time": endExecTime.Sub(startExecTime).Seconds(),
	}
	log.Fatalln(resp)
}

func RequestLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		RequestInLog(c)
		defer RequestOutLog(c)
		c.Next()
	}
}
