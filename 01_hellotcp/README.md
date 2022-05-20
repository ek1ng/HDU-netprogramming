# 基础tcp应用程序

实验要求：（需在实验过程中有明显的体现自己的元素）
要求完成一个基础tcp应用程序，附上wireshark截取的TCP握手包及核心代码，并抓取传输数据，比如“Hello World！”，编程语言及平台不限，最迟提交时间为本学期结束的最后一周。

# 实验环境说明

基于go1.8实现

## 代码实现

```go
//tcp client
func main()  {
    conn ,err := net.Dial("tcp","127.0.0.1:8888")
    if err != nil {
        fmt.Println("connected failed, error message:",err)
        return
    }
    defer conn.Close()
    inputReader := bufio.NewReader(os.Stdout)
    for {
        input, _ := inputReader.ReadString('\n')    //读取用户输入
        inputInfo := strings.Trim(input,"\r\n")
        if strings.ToUpper(inputInfo) == "q"{
            return  //如果输入q就退出
        }
        _,err = conn.Write([]byte(inputInfo))   //发送数据
        if err != nil{
            return 
        }
        buf := [512]byte{}
        n,err := conn.Read(buf[:])
        if err != nil{
            fmt.Println("get information failed, error message",err)
            return 
        }
        fmt.Println(string(buf[:n]))
    }
}

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
```

## 运行结果

![lab1_terminal.png](https://s2.loli.net/2022/05/20/Eyd39MAsoPBp7Lz.png)
![lab1_wireshark1.png](https://s2.loli.net/2022/05/20/xuvqScU5XZiRt38.png)
![lab1_wireshark2.png](https://s2.loli.net/2022/05/20/TlVfX4GSkpZroM9.png)
