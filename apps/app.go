package apps

import (
	"MyTransfer/apps/broadcast"
	"fmt"
	"github.com/gin-gonic/gin"
)

// IOC容器层：管理所有的服务的实例

// 1. HostService的实例必须注册过来，HostService才会有具体的实例， 服务启动时注册
// 2. Http 暴露模块，依赖IOC里面的HostService

var (
	broadcastService broadcast.Service
	implApps         map[string]ImplService
	ginApps          map[string]GinService
)

func RegistryImpl(svc ImplService) {
	if _, ok := implApps[svc.Name()]; ok {
		panic(fmt.Sprintf("service %s has registried", svc.Name()))
	}
	implApps[svc.Name()] = svc
	if v, ok := svc.(broadcast.Service); ok {
		broadcastService = v
	}
}

func GetImpl(name string) interface{} {
	for k, v := range implApps {
		if k == name {
			return v
		}
	}
	return nil
}

func RegistryGin(svc GinService) {
	if _, ok := implApps[svc.Name()]; ok {
		panic(fmt.Sprintf("service %s has registried", svc.Name()))
	}
	ginApps[svc.Name()] = svc
}

type ImplService interface {
	Config()
	Name() string
}

type GinService interface {
	Registry(r gin.IRouter)
	Name() string
	Config()
}

// InitImpl 用于初始化 注册到IOC容器里面的所有服务
func InitImpl() {
	for _, v := range implApps {
		v.Config()
	}
}

func LoadedGinApps() (names []string) {
	for k := range ginApps {
		names = append(names, k)
	}
	return names
}

func InitGin(r gin.IRouter) {
	for _, v := range ginApps {
		v.Config()
	}
	for _, v := range ginApps {
		v.Registry(r)
	}
}
