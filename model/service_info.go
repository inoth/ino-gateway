package model

import "encoding/json"

// redis 存储服务基本信息
type ServiceInfo struct {
	ServiceKey  string
	Version     string
	Desc        string
	Hosts       []string
	NeedAuth    bool
	NeedLicense bool // default: false
}

func (si ServiceInfo) String() []byte {
	buf, _ := json.Marshal(si)
	return buf
}
