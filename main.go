package main

import (
	"fmt"
	"net"
	"time"
)

var (
	myDeviceIpv4Address string
)

func getMyDeviceIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				fmt.Println("IPv4: ", ipnet.IP.String())
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", fmt.Errorf("your device not has IPv4")
}

func main() {
	device, err := getMyDeviceIP()
	if err != nil {
		return
	}
	myDeviceIpv4Address = fmt.Sprintf("%s:%s", device, "8080")
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
			if remoteAddr.String() == myDeviceIpv4Address {
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
