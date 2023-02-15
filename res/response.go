package res

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
)

// 自定义错误码
const (
	SuccessCode int = 100 * iota
	UndefErrorCode
	ValidErrorCode
	InternalErrorCode

	InvalidRequestErrorCode = 401
)

type Result struct {
	ErrorCode int         `json:"errno"`
	ErrorMsg  string      `json:"errmsg"`
	Data      interface{} `json:"data"`
	TraceId   interface{} `json:"trace_id"`
}

func (r Result) String() []byte {
	buf, _ := json.Marshal(r)
	return buf
}

func ResultErr(c *gin.Context, code int, err error) {
	traceId, _ := c.Get("trace_id")
	resp := &Result{ErrorCode: code, ErrorMsg: err.Error(), Data: "", TraceId: traceId}
	c.Set("result", string(resp.String()))

	c.JSON(200, resp)
	c.AbortWithError(200, err)
}

func ResultOk(c *gin.Context, code int, data interface{}) {
	traceId, _ := c.Get("trace_id")
	resp := &Result{ErrorCode: code, ErrorMsg: "ok", Data: data, TraceId: traceId}

	c.JSON(200, resp)
	c.Set("result", string(resp.String()))
}
