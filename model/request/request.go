package request

import (
	"github/inoth/gateway/model"
)

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
