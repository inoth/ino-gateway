package servicemanage

import (
	"context"
	"errors"
	"fmt"
	"github/inoth/gateway/model"
	"github/inoth/gateway/util"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/inoth/toybox/components/redis"
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
	locker       sync.RWMutex
	ServiceSlice map[string]map[string]*model.ServiceInfo
}

func (sm *ServiceManager) Init() (err error) {
	manageOnce.Do(func() {
		sm.ServiceSlice = make(map[string]map[string]*model.ServiceInfo)

		var serviceStr []string
		if serviceStr, err = redis.Rdc.SMembers(context.Background(), ServiceListCacheKey).Result(); err != nil {
			fmt.Printf("no available service cache found\n")
		}

		serviceList := make([]*model.ServiceInfo, 0, len(serviceStr))
		for i := 0; i < len(serviceStr); i++ {
			var tmp model.ServiceInfo
			if tmp, err = util.JsonMarshal[model.ServiceInfo](serviceStr[i]); err != nil {
				continue
			}
			fmt.Printf("load service %s:%s, hosts: %+v\n", tmp.ServiceKey, tmp.Version, tmp.Hosts)
			serviceList = append(serviceList, &tmp)
		}

		// 拉取服务信息到本地内存
		for _, service := range serviceList {
			if sm.ServiceSlice[service.ServiceKey] == nil {
				sm.ServiceSlice[service.ServiceKey] = make(map[string]*model.ServiceInfo)
			}
			if _, ok := sm.ServiceSlice[service.ServiceKey][service.Version]; ok {
				sm.ServiceSlice[service.ServiceKey][service.Version].AddNode(service.Hosts...)
				continue
			}
			sm.ServiceSlice[service.ServiceKey][service.Version] = service
		}

		ServiceManage = sm
	})
	return nil
}

// TODO: 获取服务优化，优先获取本地，其次查找redis
func (sm *ServiceManager) HTTPAccessMode(c *gin.Context) (*model.ServiceInfo, error) {
	path := c.Request.URL.Path
	prefixs := strings.Split(path, "/")
	if len(prefixs) < 3 {
		return nil, errors.New("does not conform to the agreed routing prefix")
	}
	serviceName := prefixs[1]
	version := prefixs[2]

	if vers, ok := sm.ServiceSlice[serviceName]; ok {
		if svcInfo, ok := vers[version]; ok {
			return svcInfo, nil
		}
	}

	return nil, errors.New("not matched service")
}

// 新增一个服务
func (sm *ServiceManager) AppendService(services ...*model.ServiceInfo) error {
	sm.locker.Lock()
	defer sm.locker.Unlock()
	for _, service := range services {
		if svc, ok := sm.ServiceSlice[service.ServiceKey]; ok {
			if ver, ok := svc[service.Version]; ok {
				// 已存在当前版本，直接新增服务host节点
				ver.AddNode(service.Hosts...)
				changeRedisServiceList(true, service)
				continue
			}
			// 创建新的服务版本号
			sm.ServiceSlice[service.ServiceKey][service.Version] = service
			changeRedisServiceList(true, service)
			continue
		}
		// 创建新的服务
		if sm.ServiceSlice[service.ServiceKey] == nil {
			sm.ServiceSlice[service.ServiceKey] = make(map[string]*model.ServiceInfo)
		}
		sm.ServiceSlice[service.ServiceKey][service.Version] = service
		changeRedisServiceList(true, service)
	}
	return nil
}

// 删除一个服务
func (sm *ServiceManager) DelService(serviceKey, version string, nodes ...model.ServerNode) error {
	sm.locker.Lock()
	defer sm.locker.Unlock()
	if _, ok := sm.ServiceSlice[serviceKey]; ok {
		return sm.DelServiceNodes(serviceKey, version, nodes...)
	}
	//TODO: 进行日志记录优化
	fmt.Println("invalid delete, no deleteable object found.")
	return nil
}

// 删除服务下节点
func (sm *ServiceManager) DelServiceNodes(serviceKey, version string, nodes ...model.ServerNode) error {
	var service *model.ServiceInfo
	var ok bool
	if service, ok = sm.ServiceSlice[serviceKey][version]; !ok {
		return errors.New("invalid delete, no deleteable object found")
	}
	// 获取原有nodes列表
	old := make([]string, 0, len(sm.ServiceSlice[serviceKey][version].Hosts))
	for _, host := range sm.ServiceSlice[serviceKey][version].Hosts {
		old = append(old, host.Host)
	}

	// 长度相等，直接全部删除，无需校验差集
	if len(old) == len(nodes) {
		delete(sm.ServiceSlice[serviceKey], version)
	} else {
		// 获取差集之后，修改内存中存活 node
		dels := make([]string, 0, len(nodes))
		for _, node := range nodes {
			dels = append(dels, node.Host)
		}
		// 保留部分
		reserve := util.Difference(old, dels)
		reserveNode := make([]model.ServerNode, 0, len(reserve))
		for _, host := range reserve {
			reserveNode = append(reserveNode, model.ServerNode{
				Host: host,
			})
		}
		sm.ServiceSlice[serviceKey][version].Hosts = reserveNode
	}

	for _, node := range nodes {
		tmp := service.Copy()
		tmp.AddNode(node)
		changeRedisServiceList(false, tmp)
	}
	return nil
}

func (sm *ServiceManager) GetServiceList() []model.ServiceInfo {
	var serviceStr []string
	var err error
	if serviceStr, err = redis.Rdc.SMembers(context.Background(), ServiceListCacheKey).Result(); err != nil {
		return make([]model.ServiceInfo, 0)
	}
	serviceList := make([]model.ServiceInfo, 0, len(serviceStr))
	for i := 0; i < len(serviceStr); i++ {
		var tmp model.ServiceInfo
		if tmp, err = util.JsonMarshal[model.ServiceInfo](serviceStr[i]); err != nil {
			continue
		}
		serviceList = append(serviceList, tmp)
	}
	return serviceList
}

func changeRedisServiceList(action bool, service *model.ServiceInfo) {
	ctx := context.Background()
	if action {
		err := redis.Rdc.SAdd(ctx, ServiceListCacheKey, service.String()).Err()
		if err != nil {
			return
		}
		// redis.Rdc.Expire(ctx, ServiceListCacheKey, time.Hour*72)
	} else {
		redis.Rdc.SRem(ctx, ServiceListCacheKey, service.String())
	}
}
