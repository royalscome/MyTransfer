package broadcast

import (
	"fmt"
)

/*
udp广播，负责不断的向外发送udp报文，表示当前局域网内设备在线
*/

var (
	OnlineDevices []DeviceInfo // 在线设备列表
	MyDevices     []DeviceInfo // 本机设备信息数组
)

// DeviceInfo 设备信息
type DeviceInfo struct {
	IP   string // 设备IP地址
	Port string // 设备端口号
	Tag  string // 设备标识符
}

// 返回本机信息字符串格式
// usage: DeviceInfo{}.String()
func (d *DeviceInfo) String() string {
	return fmt.Sprintf("%s:%s", d.IP, d.Port)
}

// MessageType 接收的消息类型
type MessageType int

const (
	ConfirmType MessageType = iota // 通知对方确认是否接收文件
	AcceptType
	RefuseType
)

type MessageData struct {
	Type    MessageType `json:"type"`
	Message string      `json:"message"`
}

func NewDefaultMessageData() *MessageData {
	return &MessageData{
		Type:    ConfirmType,
		Message: "",
	}
}
