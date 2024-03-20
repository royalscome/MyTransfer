package main

import (
	"MyTransfer/apps/broadcast"
)

func main() {
	err := broadcast.StartBroadcast()
	if err != nil {
		panic(err)
		return
	}
}

// GOOS=windows GOARCH=amd64 go build -o myTransfer.exe main.go
