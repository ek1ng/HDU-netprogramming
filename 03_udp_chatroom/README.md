# UDP_Chatroom

## 实验环境

go 1.8

## 代码实现

主要参考github仓库 https://github.com/pawalt/gopunch

### 客户端代码

```go
package main

import (
 "bufio"
 "flag"
 "fmt"
 "net"
 "os"
 "strings"
)

var token = flag.String("token", "test", "token to use for matching up pairs")
var serverAddr = flag.String("serverAddr", "127.0.0.1:1338", "stun server to connect to")

func main() {
 flag.Parse()
 p := make([]byte, 2048)
 addr := net.UDPAddr{
  IP: net.ParseIP("0.0.0.0"),
 }
 ser, err := net.ListenUDP("udp", &addr)
 if err != nil {
  panic(err)
 }

 serverTypedAddr, err := net.ResolveUDPAddr("udp", *serverAddr)
 if err != nil {
  panic(err)
 }
 fmt.Printf("Sending STUN request to %v\n", serverTypedAddr)
 _, err = ser.WriteToUDP([]byte(*token), serverTypedAddr)
 if err != nil {
  panic(err)
 }

 _, _, err = ser.ReadFromUDP(p)
 if err != nil {
  panic(err)
 }

 resp := string(p)
 resp = strings.ReplaceAll(resp, "local:", "")
 resp = strings.ReplaceAll(resp, "remote:", "")
 addresses := strings.Split(resp, "\n")

 localStrAddr := addresses[0]
 remoteStrAddr := addresses[1]

 localAddr, err := net.ResolveUDPAddr("udp", localStrAddr)
 if err != nil {
  panic(err)
 }
 remoteAddr, err := net.ResolveUDPAddr("udp", remoteStrAddr)
 if err != nil {
  panic(err)
 }

 _, err = ser.WriteToUDP([]byte(fmt.Sprintf("Connected to host at %v", localAddr)), remoteAddr)
 if err != nil {
  panic(err)
 }

 go func() {
  for {
   n, _, err := ser.ReadFromUDP(p)
   if err != nil {
    panic(err)
   }
   fmt.Printf("%s\n", p[:n])
  }
 }()

 go func() {
  reader := bufio.NewReader(os.Stdin)
  for {
   text, err := reader.ReadString('\n')
   if err != nil {
    panic(err)
   }
   text = strings.TrimSpace(text)
   _, err = ser.WriteToUDP([]byte(text), remoteAddr)
   if err != nil {
    panic(err)
   }
  }
 }()

 for {
 }
}
```

### 服务端代码

```go
package main

import (
 "fmt"
 "net"
)

func main() {
 keys := make(map[string]*net.UDPAddr)

 p := make([]byte, 2048)
 addr := net.UDPAddr{
  Port: 1338,
  IP:   net.ParseIP("0.0.0.0"),
 }
 ser, err := net.ListenUDP("udp", &addr)
 if err != nil {
  panic(err)
 }
 for {
  _, remoteaddr, err := ser.ReadFromUDP(p)
  if err != nil {
   panic(err)
  }
  fmt.Printf("Read message from %v %s \n", remoteaddr, p)
  msg := string(p)
  if addr, found := keys[msg]; found {
   err = sendInfo(ser, remoteaddr, addr)
   if err != nil {
    panic(err)
   }
   delete(keys, msg)
  } else {
   keys[msg] = remoteaddr
  }
 }
}

func sendInfo(conn *net.UDPConn, first *net.UDPAddr, sec *net.UDPAddr) error {
 _, err := conn.WriteToUDP([]byte(fmt.Sprintf(`local:%v
remote:%v`, first, sec)), first)
 if err != nil {
  return err
 }

 _, err = conn.WriteToUDP([]byte(fmt.Sprintf(`local:%v
remote:%v`, sec, first)), sec)
 if err != nil {
  return err
 }

 return nil
}
```

## 运行结果
可以实现多对客户端之间的1对1UDP连接聊天
![图 2](https://s2.loli.net/2022/05/24/VoxXasTCJQlPK9z.png)  
