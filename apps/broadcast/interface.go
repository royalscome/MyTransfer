package broadcast

import (
	"fmt"
	"time"
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
	IP     string      `json:"ip"`     // 设备IP地址
	Port   string      `json:"port"`   // 设备端口号
	Tag    string      `json:"tag"`    // 设备标识符
	Status bool        `json:"status"` // 设备状态
	Timer  *time.Timer `json:"-"`      // 设备定时器
}

// 返回本机信息字符串格式
// usage: DeviceInfo{}.String()
func (d *DeviceInfo) String() string {
	return fmt.Sprintf("%s:%s", d.IP, d.Port)
}

// StartTimer 用于启动一个持续时间为20秒的定时器。
func (d *DeviceInfo) StartTimer() {
	d.Timer = time.AfterFunc(20*time.Second, func() {
		d.Status = false
	})
}

// ResetTimer 用于将定时器重置为20秒。
func (d *DeviceInfo) ResetTimer() {
	if d.Timer != nil {
		d.Timer.Stop()
		d.StartTimer()
	}
}

// MessageType 接收的消息类型
type MessageType int

const (
	ConfirmType MessageType = iota // 通知对方确认是否接收文件
	AcceptType
	RefuseType
	AliveType // 保活
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
