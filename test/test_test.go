package test

import (
	"github/inoth/ino-gateway/components/cache"
	"github/inoth/ino-gateway/components/config"
	"github/inoth/ino-gateway/components/local"
	"github/inoth/ino-gateway/components/logger"
	servicemanage "github/inoth/ino-gateway/components/service_manage"
	"github/inoth/ino-gateway/model"
	"github/inoth/ino-gateway/register"
	"log"
	"testing"
)

var reg *register.Register

func initComponents() {
	reg = register.NewRegister(
		&local.CacheComponent{},
		&config.ViperComponent{Path: "../config"},
		&logger.ZapComponent{},
		&cache.RedisComponent{},
		&servicemanage.ServiceManager{},
	).Init()
}

func TestAddHttpProxy(t *testing.T) {
	initComponents()
	err := servicemanage.ServiceManage.AppendService(&model.ServiceInfo{
		ServiceKey:  "job",
		Version:     "v1",
		Desc:        "job machine service",
		Hosts:       []string{"http://localhost:8081"},
		NeedAuth:    false,
		NeedLicense: false,
	}, &model.ServiceInfo{
		ServiceKey:  "cmdb",
		Version:     "v1",
		Desc:        "cmdb service",
		Hosts:       []string{"http://localhost:8082"},
		NeedAuth:    false,
		NeedLicense: false,
	})
	if err != nil {
		log.Fatalln(err.Error())
	}
	t.Log("ok")
}
