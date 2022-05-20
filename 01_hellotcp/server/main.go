package main

import (
	"bufio"
	"fmt"
	"net"
)

//TCP server端

func process(conn net.Conn) {
	defer conn.Close() //关闭连接
	for {
		reader := bufio.NewReader(conn)
		var buf [128]byte
		n, err := reader.Read(buf[:]) //读取数据
		if err != nil {
			fmt.Println("Failed to connect to client, error message:", err)
		}
		recvStr := string(buf[:n])
		fmt.Println("Receive client information:", recvStr)
		conn.Write([]byte(recvStr)) //发送数据
	}
}

func main() {
	listen, err := net.Listen("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("Listen failed error message:", err)
		return
	}
	for {
		conn, err := listen.Accept() //建立连接
		if err != nil {
			fmt.Println("Establishing connection failed, error message:", err)
			continue
		}
		go process(conn) //启动一个goroutine处理连接
	}
}
