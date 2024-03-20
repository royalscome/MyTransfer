package broadcast

import (
	"MyTransfer/conf"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

// getMyDeviceIP 获取当前设备所有192.168开头的地址
func getMyDeviceIP(filePath string) error {
	// 加载配置文件
	err := conf.LoadConfigFromToml(filePath)
	if err != nil {
		return err
	}
	// 设备ip数量大致在五个以内
	MyDevices = make([]DeviceInfo, 5)

	// 获取数据
	interfaces, err := net.Interfaces()
	if err != nil {
		return err
	}

	for _, i := range interfaces {
		addrs, err := i.Addrs()
		if err != nil {
			return err
		}
		err = processAddresses(addrs)
		if err != nil {
			return err
		}
	}
	if len(MyDevices) == 0 {
		return fmt.Errorf("no IPv4 address starting with 192.168 found")
	}
	return nil
}

func processAddresses(addrs []net.Addr) error {
	for _, addr := range addrs {
		if isLocalIPv4(addr) {
			MyDevices = append(MyDevices, DeviceInfo{
				IP:   addr.(*net.IPNet).IP.String(),
				Port: *conf.Config,
				Tag:  "me",
			})
		}
	}
	return nil
}

// isLocalIPv4 判断是否为可以用的内网ip
func isLocalIPv4(addr net.Addr) bool {
	ipnet, ok := addr.(*net.IPNet)
	return ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil && strings.HasPrefix(ipnet.IP.String(), "192.168")
}

// StartBroadcast 开始广播
func StartBroadcast() error {
	err := getMyDeviceIP("etc/config.toml")
	if err != nil {
		return err
	}
	port, err := strconv.Atoi(conf.Config.Port)
	if err != nil {
		return err
	}
	broadcastAddress := &net.UDPAddr{
		IP:   net.IPv4(255, 255, 255, 255),
		Port: port,
	}
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: port,
	})
	if err != nil {
		return err
	}
	defer func() {
		err = conn.Close()
		if err != nil {
			return
		}
	}()
	go func() {
		for {
			data := make([]byte, 4096)
			n, remoteAddr, err := conn.ReadFromUDP(data)
			if err != nil {
				fmt.Println(err)
				return
			}
			//fmt.Println(remoteAddr, myDeviceIpv4Address)
			// 如果接收到的消息是本机发送的消息，则不处理
			// 循环比对本机的所有IP地址，如果接收到的消息是本机的IP地址，则不处理
			// 如果接收到的消息不是本机的IP地址，则处理
			var isMyDevice bool
			for _, address := range MyDevices {
				if remoteAddr.String() == address.String() {
					isMyDevice = true
				}
			}

			if isMyDevice {
				fmt.Println("接收到本机消息")
			} else {
				fmt.Printf("Received from address: %s data: %s\n", remoteAddr, data[:n])
			}
		}
	}()

	for {
		_, err = conn.WriteToUDP([]byte("Hello from broadcaster"), broadcastAddress)
		if err != nil {
			return err
		}
		time.Sleep(1 * time.Second)
	}
}
