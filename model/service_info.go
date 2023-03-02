package model

import (
	"encoding/json"
	"github/inoth/ino-gateway/util"
	"math/rand"
)

// redis 存储服务基本信息
type ServiceInfo struct {
	ServiceKey  string       `json:"server_key"`
	Version     string       `json:"version"`
	Desc        string       `json:"desc"`
	Hosts       []ServerNode `json:"hosts"`
	NeedAuth    bool         `json:"need_auth"`
	NeedLicense bool         `json:"need_lic"` // default: false

	MaxQps     int64  `json:"max_qps"`
	JwtSignKey string `json:"jwt_signkey"`
}

type ServerNode struct {
	// Weights int
	Host string `json:"host"`
}

func (si ServiceInfo) String() []byte {
	buf, _ := json.Marshal(si)
	return buf
}

func (si *ServiceInfo) GetHost() string {
	return si.Hosts[rand.Int()%len(si.Hosts)].Host
}

func (si *ServiceInfo) AddNode(node ...ServerNode) bool {
	n := len(si.Hosts)
	if n <= 0 {
		si.Hosts = append(si.Hosts, node...)
		return true
	}

	old := make([]string, 0, n)
	new := make([]string, 0, len(node))

	for _, v1 := range si.Hosts {
		old = append(old, v1.Host)
	}
	for _, v2 := range node {
		new = append(new, v2.Host)
	}

	adds := util.Difference(new, old)
	for i := 0; i < len(adds); i++ {
		si.Hosts = append(si.Hosts, ServerNode{
			Host: adds[i],
		})
	}
	return true
}

func (si *ServiceInfo) Copy() *ServiceInfo {
	r := &ServiceInfo{
		ServiceKey:  si.ServiceKey,
		Version:     si.Version,
		Desc:        si.Desc,
		NeedAuth:    si.NeedAuth,
		NeedLicense: si.NeedLicense,
	}
	return r
}
