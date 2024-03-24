package main

import (
	"MyTransfer/apps"
	_ "MyTransfer/apps/all"
	"MyTransfer/conf"
	"MyTransfer/protocol"
	"fmt"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
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
	return m.http.Start()
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

func main() {
	// 加载配置文件
	err := conf.LoadConfigFromToml(filePath)
	config = conf.C()
	//return

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
