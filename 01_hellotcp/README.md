# 基础tcp应用程序

实验要求：（需在实验过程中有明显的体现自己的元素）
要求完成一个基础tcp应用程序，附上wireshark截取的TCP握手包及核心代码，并抓取传输数据，比如“Hello World！”，编程语言及平台不限，最迟提交时间为本学期结束的最后一周。

# 实验环境说明

基于go1.8实现

## 代码实现

读写确认 + go routine多端连接通信

### 客户端代码

```go
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
```

### 服务端代码

```go
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

```

## 运行结果

![lab1_terminal.png](https://s2.loli.net/2022/05/20/Eyd39MAsoPBp7Lz.png)
![lab1_wireshark1.png](https://s2.loli.net/2022/05/20/xuvqScU5XZiRt38.png)
![lab1_wireshark2.png](https://s2.loli.net/2022/05/20/TlVfX4GSkpZroM9.png)
