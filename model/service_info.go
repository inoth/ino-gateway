package model

import (
	"encoding/json"
	"github/inoth/ino-gateway/util"
	"math/rand"
)

// redis 存储服务基本信息
type ServiceInfo struct {
	ServiceKey  string
	Version     string
	Desc        string
	Hosts       []ServerNode
	NeedAuth    bool
	NeedLicense bool // default: false
}

type ServerNode struct {
	// Weights int
	Host string
}

func (si ServiceInfo) String() []byte {
	buf, _ := json.Marshal(si)
	return buf
}

func (si *ServiceInfo) GetHost() string {
	return si.Hosts[rand.Intn(len(si.Hosts))].Host
}

func (si *ServiceInfo) AddNode(node ...ServerNode) bool {
	old := make([]string, 0, len(si.Hosts))
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
