package broadcast

import "net"

type Service interface {
	// QueryOnlineDevices 查询在线设备
	QueryOnlineDevices() []DeviceInfo
	// SendMessageUseUDP 利用udp发送消息
	SendMessageUseUDP(*net.UDPConn, *UDPMessage) error
}

type UDPMessage struct {
	Message string `json:"message"` // 消息主题，需要符合json序列化
	Address string `json:"address"` // 需要向谁发送消息，对方的udp地址
}

func NewUDPMessage() *UDPMessage {
	return &UDPMessage{
		Message: "",
		Address: "",
	}
}
