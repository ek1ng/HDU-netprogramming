package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sync"
)

const (
	CONN_PORT = ":8080"
	CONN_TYPE = "tcp"

	MSG_DISCONN = "Connection closed.\n"
)

var waitGroup sync.WaitGroup


func Read(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		str, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf(MSG_DISCONN)
			waitGroup.Done()
			return
		}
		fmt.Print(str)
	}
}


func Write(conn net.Conn) {
	// read from stdin
	reader := bufio.NewReader(os.Stdin)
	// write to connection
	writer := bufio.NewWriter(conn)

	for {
		str, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			waitGroup.Done()
			os.Exit(1)
		}
		_, err = writer.WriteString(str)
		if err != nil {
			fmt.Println(err)
			waitGroup.Done()
			os.Exit(1)
		}
		err = writer.Flush()
		if err != nil {
			fmt.Println(err)
			waitGroup.Done()
			os.Exit(1)
		}
	}
}


func main() {
	// goroutine
	waitGroup.Add(1)

	// connect to the socket
	conn, err := net.Dial(CONN_TYPE, CONN_PORT)
	if err != nil {
		fmt.Println(err)
	}
	
	go Read(conn)
	go Write(conn)

	waitGroup.Wait()
}