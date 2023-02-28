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
		// &config.ViperComponent{ConfKeyPrefix: "../config/dev"},
		&config.ViperComponent{},
		&logger.ZapComponent{},
		&cache.RedisComponent{},
		// &servicemanage.ServiceManager{},
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
	fmt.Printf("\n%+v\n", servicemanage.ServiceManage.GetServiceList())
	t.Log("ok")
}

func TestRemoveHttpProxy(t *testing.T) {
	var (
		serviceKey = "job"
		version    = "v1"
		nodes      = []model.ServerNode{{Host: "http://localhost:8081"}}
	)
	initComponents()

	err := servicemanage.ServiceManage.DelService(serviceKey, version, nodes...)
	if err != nil {
		log.Fatalln(err.Error())
	}
	fmt.Printf("\n%+v\n", servicemanage.ServiceManage.GetServiceList())
	t.Log("ok")
}

func TestUtilIntersect(t *testing.T) {
	type TTTTT struct {
		Host string
	}
	var (
		new = []string{"http://localhost:8082", "http://localhost:8081"}
		old = []string{"http://localhost:8082", "http://localhost:8085", "http://localhost:8081"}
	)

	fmt.Println(util.Intersect(new, old))
	fmt.Println(util.Difference(old, new))
}

func TestConfGet(t *testing.T) {
	initComponents()

	// fmt.Println(config.Cfg.GetInt("proxy.http.read_timeout"))
	// fmt.Println(config.Cfg.GetString("zap.err_log"))
	fmt.Println(config.Cfg.GetString("base.server_pord"))
}
