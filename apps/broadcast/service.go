package broadcast

type Service interface {
	// QueryOnlineDevices 查询在线设备
	QueryOnlineDevices() []DeviceInfo
}
