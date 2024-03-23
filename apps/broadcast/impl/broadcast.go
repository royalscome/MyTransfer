package impl

import (
	"MyTransfer/apps"
	"MyTransfer/apps/broadcast"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
)

// 接口实现的静态检查
var (
	_    broadcast.Service = (*BroadcastServiceImpl)(nil)
	impl                   = &BroadcastServiceImpl{}
)

type BroadcastServiceImpl struct {
	l logger.Logger
}

func NewBroadcastServiceImpl() *BroadcastServiceImpl {
	return &BroadcastServiceImpl{
		l: zap.L().Named("broadcast"),
	}
}

func (i *BroadcastServiceImpl) Config() {
	i.l = zap.L().Named("broadcast")
}

func (i *BroadcastServiceImpl) Name() string {
	return broadcast.AppName
}

func (i *BroadcastServiceImpl) QueryOnlineDevices() []broadcast.DeviceInfo {
	return broadcast.OnlineDevices
}

func init() {
	apps.RegistryImpl(impl)
}
