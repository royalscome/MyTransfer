package main

import (
	"MyTransfer/apps"
	_ "MyTransfer/apps/all"
	"MyTransfer/apps/websocket/impl"
	"MyTransfer/conf"
	"MyTransfer/protocol"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var (
	config   *conf.Config
	filePath string = "etc/config.toml"
)

// 用于管理所有需要启动的服务
type manage struct {
	http *protocol.HttpService
	udp  *protocol.UDPService
	l    logger.Logger
}

func newManage() *manage {
	return &manage{
		http: protocol.NewHttpService(),
		udp:  protocol.NewUDPService(),
		l:    zap.L().Named("MAIN"),
	}
}

func (m *manage) Start() error {
	go func() {
		m.udp.Start()
	}()
	return m.http.Start(m.udp.GetConn())
}

// Stop 处理来自外部的终端信号
func (m *manage) waitStop(ch <-chan os.Signal) {
	for v := range ch {
		switch v {
		default:
			m.l.Infof("received signal %s", v)
			m.udp.Stop()
			m.http.Stop()
		}
	}
}

func loadGlobalLogger() error {
	var (
		logInitMsg string
		level      zap.Level
	)
	lc := conf.C().Log
	lv, err := zap.NewLevel(lc.Level)
	if err != nil {
		logInitMsg = fmt.Sprintf("%s, use default level INFO", err)
		level = zap.InfoLevel
	} else {
		level = lv
		logInitMsg = fmt.Sprintf("log level: %s", lv)
	}
	zapConfig := zap.DefaultConfig()
	zapConfig.Level = level
	// 配置日志轮转
	zapConfig.Files.RotateOnStartup = false
	switch lc.To {
	case conf.ToStdout:
		zapConfig.ToStderr = true
		zapConfig.ToFiles = false
	case conf.ToFile:
		zapConfig.Files.Name = "api.log"
		zapConfig.Files.Path = lc.PathDir
	}
	switch lc.Format {
	case conf.JSONFormat:
		zapConfig.JSON = true
	}
	if err := zap.Configure(zapConfig); err != nil {
		return err
	}
	zap.L().Named("INIT").Info(logInitMsg)
	return nil
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Make sure we close the connection when the function returns
	defer ws.Close()
	// Initialize a new Connection
	conn, err := impl.InitConnection(ws)
	if err != nil {
		log.Fatal("init connection:", err)
	}

	// Loop indefinitely
	for {
		// Read message from browser
		msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}

		// Print the message to the console
		log.Printf("recv: %s", msg)

		// TODO: Process the message
		// For now, we'll just echo the same message back
		processedMsg := msg

		// Write message back to browser
		err = conn.WriteMessage(processedMsg)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func main() {
	// 加载配置文件
	err := conf.LoadConfigFromToml(filePath)
	config = conf.C()
	http.HandleFunc("/ws", handleConnections)
	log.Fatal(http.ListenAndServe(":2000", nil))
	return

	// 开启http服务
	// 加载日志
	if err = loadGlobalLogger(); err != nil {
		panic(err)
	}
	// 注册HostService的实例
	apps.InitImpl()
	svc := newManage()
	ch := make(chan os.Signal, 1)
	defer close(ch)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP, syscall.SIGINT)
	go svc.waitStop(ch)
	if err = svc.Start(); err != nil {
		panic(err)
	}
}

// GOOS=windows GOARCH=amd64 go build -o myTransfer.exe main.go
