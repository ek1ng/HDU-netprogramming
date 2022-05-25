package main

import (
	"bufio"
	"fmt"
	"net"
)

//TCP server

func handleConnection(conn net.Conn) {
	defer conn.Close()
	for {
		reader := bufio.NewReader(conn)
		var buf [128]byte
		n, err := reader.Read(buf[:])
		if err != nil {
			fmt.Println("Failed to connect to client, error message:", err)
		}
		recvStr := string(buf[:n])
		fmt.Println("Receive message:", recvStr)
		conn.Write([]byte(recvStr)) //发送数据
	}
}

func main() {
	// listen on all interfaces
	ln, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		// handle error
		fmt.Println("Listen failed error message:", err)
		return
	}
	for {
		// close the listener when the application closes
		conn, err := ln.Accept()
		if err != nil {
			// handle error
			fmt.Println("Establishing connection failed, error message:", err)
			continue
		}
		go handleConnection(conn)
	}
}
