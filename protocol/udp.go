package protocol

import (
	"MyTransfer/apps/broadcast"
	"MyTransfer/conf"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"net"
	"strconv"
	"strings"
	"time"
)

type UDPService struct {
	address *net.UDPAddr
	conn    *net.UDPConn
	l       logger.Logger
	stop    chan bool
}

func (s *UDPService) GetConn() *net.UDPConn {
	return s.conn
}

func NewUDPService() *UDPService {
	port, err := strconv.Atoi(conf.C().UDP.Port)
	if err != nil {
		panic(err)
	}
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: port,
	})
	if err != nil {
		panic(err)
	}
	return &UDPService{
		address: &net.UDPAddr{
			IP:   net.IPv4(255, 255, 255, 255),
			Port: port,
		},
		conn: conn,
		l:    zap.L().Named("UDP service"),
		stop: make(chan bool),
	}
}

func (u *UDPService) Start() error {
	myDevices, err := getMyDeviceIP()
	if err != nil {
		return err
	}
	go func() {
		for {
			select {
			case <-u.stop:
				return
			default:
				data := make([]byte, 4096)
				n, remoteAddr, err := u.conn.ReadFromUDP(data)
				if err != nil {
					if strings.Contains(err.Error(), "use of closed network connection") {
						// Ignore the error if the connection is closed
						return
					}
					u.l.Errorf("receive message error: %v", err)
					return
				}
				// 如果接收到的消息是本机发送的消息，则不处理
				// 循环比对本机的所有IP地址，如果接收到的消息是本机的IP地址，则不处理
				// 如果接收到的消息不是本机的IP地址，则处理
				var isMyDevice bool
				for _, address := range myDevices {
					if remoteAddr.String() == address.String() {
						isMyDevice = true
					}
				}

				if isMyDevice {
					fmt.Println("接收到本机消息", remoteAddr, string(data[:n]))
					err = JudgeMessageType(string(data[:n]), remoteAddr)
					if err != nil {
						u.l.Warnf("message error: %s", err)
					}
				} else {
					fmt.Printf("Received from address: %s data: %s\n", remoteAddr, data[:n])
					err = JudgeMessageType(string(data[:n]), remoteAddr)
					if err != nil {
						u.l.Warnf("message error: %s", err)
					}
				}
			}
		}
	}()
	for {
		select {
		case <-u.stop:
			return nil
		default:
			aliveMessage := &broadcast.MessageData{
				Type:    broadcast.AliveType,
				Message: "",
			}
			message, err := json.Marshal(aliveMessage)
			if err != nil {
				u.l.Warnf("marshal error: %s", err)
			}
			_, err = u.conn.WriteToUDP([]byte(message), u.address)
			if err != nil {
				u.l.Warnf("keep alive message send error: %s", err)
				return err
			}
			time.Sleep(10 * time.Second)
		}
	}
}

func (u *UDPService) Stop() {
	u.l.Info("start stop udp conn")
	close(u.stop)
	if err := u.conn.Close(); err != nil {
		u.l.Warnf("stop udp conn error: %v", err)
	}
	u.l.Info("stop udp conn success")
}

// getMyDeviceIP 获取当前设备所有192.168开头的地址
func getMyDeviceIP() ([]broadcast.DeviceInfo, error) {
	// 设备ip数量大致在五个以内
	myDevices := make([]broadcast.DeviceInfo, 5)

	// 获取数据
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, i := range interfaces {
		addrs, err := i.Addrs()
		if err != nil {
			return nil, err
		}
		err = processAddresses(addrs, &myDevices, conf.C().UDP.Port)
		if err != nil {
			return nil, err
		}
	}
	if len(myDevices) == 0 {
		return nil, fmt.Errorf("no IPv4 address starting with 192.168 found")
	}
	return myDevices, nil
}

func processAddresses(addrs []net.Addr, myDevices *[]broadcast.DeviceInfo, port string) error {
	for _, addr := range addrs {
		if isLocalIPv4(addr) {
			*myDevices = append(*myDevices, broadcast.DeviceInfo{
				IP:   addr.(*net.IPNet).IP.String(),
				Port: port,
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

func isValidJSON(s string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}

// JudgeMessageType 根据接收到的消息判断后续操作
func JudgeMessageType(message string, address *net.UDPAddr) error {
	if !isValidJSON(message) {
		return errors.New("invalid JSON format")
	}
	var marshalMessage = broadcast.NewDefaultMessageData()
	err := json.Unmarshal([]byte(message), marshalMessage)
	if err != nil {
		return err
	}
	switch marshalMessage.Type {
	case broadcast.AliveType:
		keepAlive(address)
	case broadcast.ConfirmType:
	case broadcast.AcceptType:
	case broadcast.RefuseType:
	default:
		return errors.New("MessageType not validate")
	}
	return nil
}

func keepAlive(address *net.UDPAddr) {
	// 如果设备地址已存在列表中说明需要新增
	var isExist bool = false
	for _, device := range broadcast.OnlineDevices {
		if device.String() == address.String() {
			isExist = true
			device.ResetTimer()
		}
	}
	if !isExist {
		// 如果状态为false说明设备初次出现，新增一个并且开启定时器
		newDevice := broadcast.DeviceInfo{
			IP:     address.IP.String(),
			Port:   strconv.Itoa(address.Port),
			Tag:    "other",
			Status: true,
		}
		broadcast.OnlineDevices = append(broadcast.OnlineDevices, newDevice)
		newDevice.StartTimer()
	}

}
