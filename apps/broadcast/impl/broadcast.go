package impl

import "MyTransfer/apps/broadcast"

// 接口实现的静态检查
var _ broadcast.Service = (*BroadcastServiceImpl)(nil)

type BroadcastServiceImpl struct {
}

func (i *BroadcastServiceImpl) QueryOnlineDevices() []broadcast.DeviceInfo {
	return broadcast.OnlineDevices
}
