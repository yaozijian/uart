package main

import (
	"fmt"
	"net"
	"os"

	"github.com/yaozijian/uart"
)

var (
	table *uart.TypeTable
)

func main() {
	port := ":4444"
	if len(os.Args) > 1 {
		port = ":" + os.Args[1]
	}

	table = initTable()

	l, _ := net.Listen("tcp4", port)
	for {
		c, _ := l.Accept()
		go proc(c)
	}
}

func proc(conn net.Conn) {

	addr := fmt.Sprintf("%s -> %s", conn.RemoteAddr(), conn.LocalAddr())
	fmt.Println("Connected:", addr)

	defer conn.Close()
	defer fmt.Println("Disconnected:", addr)

	cfg := &uart.UartConfig{
		Table:    table,
		Protocol: uart.JsonProtocol,
		Src:      conn,
		Dst:      conn,
		RecvBuf:  100,
		SendBuf:  100,
	}
	this := uart.NewUart(cfg)

	for {
		msg, ok := <-this.RecvChnl
		if ok {
			onmsg(this, msg)
		} else {
			break
		}
	}
}

func onmsg(this *uart.Uart, msg interface{}) {
	switch val := msg.(type) {
	case *echomsg:
		this.SendChnl <- val
	case *mathmsg:
		switch val.Operator {
		case operator_add:
			val.Result = val.A + val.B
		case operator_del:
			val.Result = val.A - val.B
		case operator_mul:
			val.Result = val.A * val.B
		case operator_div:
			val.Result = val.A / val.B
		}
		this.SendChnl <- val
	}
}
