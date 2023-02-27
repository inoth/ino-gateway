package test

import (
	"fmt"
	"github/inoth/ino-gateway/components/cache"
	"github/inoth/ino-gateway/components/config"
	"github/inoth/ino-gateway/components/local"
	"github/inoth/ino-gateway/components/logger"
	servicemanage "github/inoth/ino-gateway/components/service_manage"
	"github/inoth/ino-gateway/model"
	"github/inoth/ino-gateway/register"
	"github/inoth/ino-gateway/util"
	"log"
	"os"
	"testing"
)

var reg *register.Register

func initComponents() {
	os.Setenv("GORUNEVN", "dev")
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
		Hosts:       []model.ServerNode{{Host: "http://localhost:8081"}},
		NeedAuth:    true,
		NeedLicense: false,
	}, &model.ServiceInfo{
		ServiceKey:  "job",
		Version:     "v1",
		Desc:        "job machine service",
		Hosts:       []model.ServerNode{{Host: "http://localhost:8083"}},
		NeedAuth:    true,
		NeedLicense: false,
	}, &model.ServiceInfo{
		ServiceKey:  "cmdb",
		Version:     "v1",
		Desc:        "cmdb service",
		Hosts:       []model.ServerNode{{Host: "http://localhost:8082"}},
		NeedAuth:    false,
		NeedLicense: false,
	})
	if err != nil {
		log.Fatalln(err.Error())
	}
	fmt.Printf("%d", (1 << uint(config.Cfg.GetInt("proxy.http.max_header_bytes"))))
	// for _, v := range servicemanage.ServiceManage.ServiceSlice {
	// 	for _, val := range v {
	// 		fmt.Printf("%+v\n", val.Hosts)
	// 	}
	// }
	fmt.Printf("%+v", servicemanage.ServiceManage.GetServiceList())
	t.Log("ok")
}

func TestUtilIntersect(t *testing.T) {
	var (
		new = []string{"http://localhost:8081", "http://localhost:8083"}
		old = []string{"http://localhost:8082", "http://localhost:8085", "http://localhost:8081"}
	)

	fmt.Println(util.Intersect(new, old))
	fmt.Println(util.Difference(new, old))
}
