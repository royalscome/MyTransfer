package broadcast

import "net"

func GetPort() {

}

// getMyDeviceIP 获取当前设备所有192.168开头的地址
func getMyDeviceIP() error {
	interfaces, err := net.Interfaces()
	if err != nil {
		return err
	}

	for _, i := range interfaces {
		addrs, err := i.Addrs()
		if err != nil {
			return err
		}
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {

			}
		}
	}

	return nil
}
