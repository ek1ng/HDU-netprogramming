package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"os"
)

//TCP client

func main()  {
    // connect to this socket
    conn ,err := net.Dial("tcp","127.0.0.1:8080")
    if err != nil {
        // handle error
        fmt.Println("connected failed, error message:",err)
        return
    }
    // close connection
    defer conn.Close()
    inputReader := bufio.NewReader(os.Stdout)
    // send data to server
    for {
        // Read input from stdin
        input, _ := inputReader.ReadString('\n')
        inputInfo := strings.Trim(input,"\r\n")

        // Write data to connection
        _,err = conn.Write([]byte(inputInfo))
        if err != nil{
            fmt.Println("Write data failed, error message",err)
            return 
        }
        buf := [512]byte{}

        // Read data from connection
        n,err := conn.Read(buf[:])
        if err != nil{
            fmt.Println("Read data failed, error message",err)
            return 
        }

        // check write data with read data
        if (inputInfo == string(buf[:n])){
            fmt.Println("Send message successfully")
        }
    }
}

