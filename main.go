package main

import (
	"MyTransfer/apps/broadcast"
	"MyTransfer/conf"
)

var (
	config   *conf.Config
	filePath string = "etc/config.toml"
)

func main() {
	// 加载配置文件
	err := conf.LoadConfigFromToml(filePath)
	config = conf.C()

	// 开始广播
	err = broadcast.StartBroadcast(config.UDP)
	if err != nil {
		panic(err)
		return
	}
}

// GOOS=windows GOARCH=amd64 go build -o myTransfer.exe main.go
