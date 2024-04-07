package transfer

import (
	"MyTransfer/conf"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

var (
	AppName = "transfer"
)

func InitSendFromTcp(path string) {
	conn, err := net.Dial("tcp", conf.C().TCP.TcpAddr())
	if err != nil {
		fmt.Println("Failed to connect:", err)
	}
	defer conn.Close()
	sendFile(conn, path)
}

func sendFile(conn net.Conn, path string) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Failed to open file:", err)
		return
	}
	defer file.Close()
	time.Sleep(1 * time.Second)
	buffer := make([]byte, 1024)
	for {
		n, err := file.Read(buffer)
		if err != nil {
			fmt.Println("Failed to read file:", err)
			break
		}
		conn.Write(buffer[:n])
	}

}

func InitReciveFromTcp(path string) {
	listen, err := net.Listen("tcp", conf.C().TCP.TcpAddr())
	if err != nil {
		fmt.Println("Failed to listen:", err)
	}
	defer listen.Close()
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("Failed to accept:", err)
			continue
		}

		go handleRequest(conn, path)
	}
}

func handleRequest(conn net.Conn, path string) {
	defer conn.Close()

	//reader := bufio.NewReader(conn)
	//msg, _ := reader.ReadString('\n')
	//fmt.Println("Received request:", msg)

	receiveFile(conn, path)
}
func receiveFile(conn net.Conn, path string) {
	file, err := os.Create(path)
	if err != nil {
		fmt.Println("Failed to create file:", err)
		return
	}
	defer file.Close()

	_, err = io.Copy(file, conn)
	if err != nil {
		fmt.Println("Failed to copy file:", err)
		return
	}

	fmt.Println("File received successfully")
}
