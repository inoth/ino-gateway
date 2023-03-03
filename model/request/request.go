package request

import (
	"github/inoth/ino-gateway/model"

	"github.com/inoth/ino-toybox/res"

	"github.com/gin-gonic/gin"
)

func RequestJsonParamHandler[T interface{}](c *gin.Context) (T, bool) {
	var req T
	if err := c.ShouldBindJSON(&req); err != nil {
		res.ResultErr(c, res.ParamErrorCode, err)
		return req, false
	}
	return req, true
}

func RequestQueryParamHandler[T interface{}](c *gin.Context) (T, bool) {
	var req T
	if err := c.ShouldBindQuery(&req); err != nil {
		res.ResultErr(c, res.ParamErrorCode, err)
		return req, false
	}
	return req, true
}

func RequestXMLParamHandler[T interface{}](c *gin.Context) (T, bool) {
	var req T
	if err := c.ShouldBindXML(&req); err != nil {
		res.ResultErr(c, res.ParamErrorCode, err)
		return req, false
	}
	return req, true
}

type ServiceNodeRequests struct {
	ServiceKey string `json:"service_key"`
	Version    string `json:"version"`
	Desc       string `json:"desc"`
	// Host        model.ServerNode `json:"host"`
	Host        string `json:"host"`
	NeedAuth    bool   `json:"need_auth"`
	NeedLicense bool   `json:"need_lic"` // default: false
}

type ServiceNodeRemoveRequest struct {
	ServiceKey string             `json:"server_key"`
	Version    string             `json:"version"`
	Hosts      []model.ServerNode `json:"hosts"`
}
