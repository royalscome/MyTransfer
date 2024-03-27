package impl

import (
	"MyTransfer/apps"
	"MyTransfer/apps/broadcast"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"net"
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

func (i *BroadcastServiceImpl) QueryOnlineDevices() *broadcast.OnlineDevicesData {
	return &broadcast.OnlineDevicesData{
		DeviceList: broadcast.OnlineDevices,
		Total:      len(broadcast.OnlineDevices),
	}
}

func (i *BroadcastServiceImpl) SendMessageUseUDP(conn *net.UDPConn, UDPMessage *broadcast.UDPMessage) error {
	addr, err := net.ResolveUDPAddr("udp", UDPMessage.Address)
	if err != nil {
		return err
	}
	if _, err := conn.WriteToUDP([]byte(UDPMessage.Message), addr); err != nil {
		return err
	}
	return nil
}

func init() {
	apps.RegistryImpl(impl)
}
