package main

import (
	"fmt"
	"net"
	"strings"
	"time"
)

var (
	myDeviceIpv4Address []string
)

//	func getMyDeviceIP() (string, error) {
//		//addrs, err := net.InterfaceAddrs()
//		//if err != nil {
//		//	fmt.Println(err)
//		//	return "", err
//		//}
//		//
//		//for _, addr := range addrs {
//		//	fmt.Println(addr)
//		//	if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
//		//		if ipnet.IP.To4() != nil {
//		//			fmt.Println("IPv4: ", ipnet.IP.String())
//		//			//return ipnet.IP.String(), nil
//		//		}
//		//	}
//		//}
//		interfaces, err := net.Interfaces()
//		if err != nil {
//			return "", err
//		}
//
//		for _, i := range interfaces {
//			fmt.Println(i.Name)
//			if i.Name != "eth0" {
//				continue
//			}
//
//			addrs, err := i.Addrs()
//			if err != nil {
//				return "", err
//			}
//
//			for _, addr := range addrs {
//				if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
//					if ipnet.IP.To4() != nil {
//						return ipnet.IP.String(), nil
//					}
//				}
//			}
//		}
//		return "", fmt.Errorf("your device not has IPv4")
//	}
func getMyDeviceIP() ([]string, error) {
	var ips []string
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, i := range interfaces {
		addrs, err := i.Addrs()
		if err != nil {
			return nil, err
		}

		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					ip := ipnet.IP.String()
					if strings.HasPrefix(ip, "192.168") {
						ips = append(ips, ip)
					}
				}
			}
		}
	}

	if len(ips) == 0 {
		return nil, fmt.Errorf("no IPv4 address starting with 192.168 found")
	}

	return ips, nil
}

func main() {
	devices, err := getMyDeviceIP()
	if err != nil {
		return
	}
	for _, dev := range devices {
		myDeviceIpv4Address = append(myDeviceIpv4Address, fmt.Sprintf("%s:8080", dev))
	}
	if err != nil {
		panic(err)
	}
	broadcastAddress := &net.UDPAddr{
		IP:   net.IPv4(255, 255, 255, 255),
		Port: 8080,
	}

	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 8080,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

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
			for _, address := range myDeviceIpv4Address {
				if remoteAddr.String() == address {
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
			fmt.Println(err)
			return
		}
		time.Sleep(1 * time.Second)
	}
}

// GOOS=windows GOARCH=amd64 go build -o myTransfer.exe main.go
