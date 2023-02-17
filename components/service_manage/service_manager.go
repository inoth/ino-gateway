package servicemanage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github/inoth/ino-gateway/model"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/inoth/ino-toybox/components/redis"
)

const (
	ServiceListCacheKey = "_gateway_:servicelist"
)

var (
	ServiceManage *ServiceManager
	manageOnce    sync.Once
)

// map[serviceName]map[version]*ServiceInfo
// 运行在 redis 初始化之后
type ServiceManager struct {
	m            sync.RWMutex
	ServiceSlice map[string]map[string]*model.ServiceInfo
}

func (sm *ServiceManager) Init() error {
	manageOnce.Do(func() {
		ServiceManage = &ServiceManager{
			ServiceSlice: make(map[string]map[string]*model.ServiceInfo),
		}
		var serviceStr []string
		if err := redis.Rdc.SMembers(context.Background(), ServiceListCacheKey).ScanSlice(&serviceStr); err != nil {
			fmt.Printf("no available service found")
		}

		serviceList := make([]*model.ServiceInfo, 0, len(serviceStr))
		var err error
		for i := 0; i < len(serviceStr); i++ {
			err = json.Unmarshal([]byte(serviceStr[i]), serviceList[i])
			if err != nil {
				continue
			}
		}

		// 拉取服务信息到本地内存
		for _, service := range serviceList {
			if _, ok := sm.ServiceSlice[service.ServiceKey][service.Version]; ok {
				sm.ServiceSlice[service.ServiceKey][service.Version].Hosts = append(sm.ServiceSlice[service.ServiceKey][service.Version].Hosts, service.Hosts...)
				continue
			}
			sm.ServiceSlice[service.ServiceKey][service.Version] = service
		}
	})
	return nil
}

// TODO: 获取服务优化，优先获取本地，其次查找redis
func (sm *ServiceManager) HTTPAccessMode(c *gin.Context) (*model.ServiceInfo, error) {
	path := c.Request.URL.Path
	prefixs := strings.Split(path, "/")
	if len(prefixs) < 2 {
		return nil, errors.New("does not conform to the agreed routing prefix")
	}
	serviceName := prefixs[0]
	version := prefixs[1]

	if vers, ok := sm.ServiceSlice[serviceName]; ok {
		if svcInfo, ok := vers[version]; ok {
			return svcInfo, nil
		}
	}

	return nil, errors.New("not matched service")
}

// 新增一个服务
func (sm *ServiceManager) AppendService(service *model.ServiceInfo) error {
	sm.m.Lock()
	defer sm.m.Unlock()
	if svc, ok := sm.ServiceSlice[service.ServiceKey]; ok {
		if ver, ok := svc[service.Version]; ok {
			// 已存在当前版本，直接新增服务host节点
			ver.Hosts = append(ver.Hosts, service.Hosts...)
			changeRedisServiceList(true, service)
			return nil
		}
		// 创建新的服务版本号
		sm.ServiceSlice[service.ServiceKey][service.Version] = service
		changeRedisServiceList(true, service)
		return nil
	}
	// 创建新的服务
	sm.ServiceSlice[service.ServiceKey][service.Version] = service
	changeRedisServiceList(true, service)
	return nil
}

// 删除一个服务
func (sm *ServiceManager) DelService(service *model.ServiceInfo) error {
	sm.m.Lock()
	defer sm.m.Unlock()
	if _, ok := sm.ServiceSlice[service.ServiceKey]; ok {
		delete(sm.ServiceSlice[service.ServiceKey], service.Version)
		changeRedisServiceList(false, service)
		return nil
	}
	//TODO: 进行日志记录优化
	fmt.Println("invalid delete, no deleteable object found.")
	return nil
}

func changeRedisServiceList(action bool, service *model.ServiceInfo) {
	if action {
		redis.Rdc.SAdd(context.Background(), ServiceListCacheKey, service.String())
	} else {
		redis.Rdc.SRem(context.Background(), ServiceListCacheKey, service.String())
	}
}
