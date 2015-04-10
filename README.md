
Go网络编程中经常需要编写从网络连接读取消息和向网络连接写入消息的代码，而且常常需要使用gob、xml、json等编解码器来处理多种不同类型的消息。为了能够并发地读取和写入，通常需要用两个不同的goroutine来处理读取和写入；为了提升性能，服务器端还经常使用另外的goroutine来处理收到的消息，处理完成后再向客户端发送回应。这些代码还要考虑网络连接可能断开的问题。为完成这些功能所需的代码虽然稍微复杂，但是基本模式是不变的。为简化上述处理，笔者编写了这个包。

#使用方法

### 注册类型

```
func initTable() *uart.TypeTable {
	table := uart.NewTypeTable()
	table.RegisterType(mathmsg{})
	table.RegisterType(echomsg(""))
	return table
}
```
注意：不要注册指针类型，虽然本包可以正确处理，但是其行为通常不是你所期望的。

### 创建Uart对象

```
cfg := &uart.UartConfig{
	Table:    table,
	Protocol: uart.JsonProtocol,
	Src:      conn,
	Dst:      conn,
	RecvBuf:  100,
	SendBuf:  100,
}
this := uart.NewUart(cfg)
```
1. table是前面已经建立的类型表
2. uart.JsonProtocol是本包默认提供的三种基本协议之一，此外还有uart.GobProtocol和uart.XmlProtocol可用。
3. conn是已经建立的网络连接，可以是net.Conn、\*net.TCPConn、\*net.UDPConn、、\*websocket.Conn等类型。
4. RecvBuf和SendBuf分别表示收发程道的缓冲区大小，即程道最多可容纳多少条消息。


### 使用Uart对象提供的收发程道来收发消息

```
func doEcho(this *uart.Uart) bool {
	msg := echomsg(randomstr())
	this.SendChnl <- msg
	fmt.Println("Send:", msg)

	recv, ok := <-this.RecvChnl
	if ok {
		if obj, _ := recv.(*echomsg); obj != nil {
			fmt.Println("Recv:", *obj)
		}
	}
	return ok
}
```

* 必须导出结构体中需要交换的字段(字段名称以大写字母开头)，因为编解码器不会处理不导出的字段。
* 可以向发送程道发送已经注册类型的值，或者相应的指针值。
* 从接收程道收到的数据是interface{}类型，可通过类型断言判断具体类型；具体类型会是某种已经注册类型的指针。
* 出现网络错误时Uart对象会停止运行，关闭接收程道，可以通过接收程道已经关闭来断定网络连接已经断开。

#关于包名

本包的作用与硬件中的UART(universal asynchronous receiver-transmitter,通用异步收发器)的作用相似，所以使用了uart作为包名称。

