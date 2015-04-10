package main

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"

	"github.com/yaozijian/uart"
)

func main() {

	addr := "127.0.0.1:4444"
	if len(os.Args) > 1 {
		addr = os.Args[1]
	}

	conn, err := net.Dial("tcp4", addr)
	if err != nil {
		fmt.Println("连接失败:", err)
		return
	}
	addr = fmt.Sprintf("%s -> %s", conn.LocalAddr(), conn.RemoteAddr())
	fmt.Println("Connected:", addr)

	rand.Seed(time.Now().Unix())

	table := initTable()

	cfg := &uart.UartConfig{
		Table:    table,
		Protocol: uart.JsonProtocol,
		Src:      conn,
		Dst:      conn,
		RecvBuf:  100,
		SendBuf:  100,
	}
	this := uart.NewUart(cfg)

	ok := true

	for ok {
		if rand.Int()%2 == 0 {
			ok = doEcho(this)
		} else {
			ok = doMath(this)
		}
		fmt.Println("")
		time.Sleep(time.Second * 2)
	}

	fmt.Println("Disconnected:", addr)
}

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

func doMath(this *uart.Uart) bool {
	msg := &mathmsg{
		A:        int(rand.Uint32()%100 + 1),
		B:        int(rand.Uint32()%100 + 1),
		Operator: operator(int(rand.Uint32() % 4)),
	}
	this.SendChnl <- msg

	recv, ok := <-this.RecvChnl
	if ok {
		if obj, _ := recv.(*mathmsg); obj != nil {
			fmt.Printf("%2d %s %2d = %d\n", obj.A, obj.Operator, obj.B, obj.Result)
		}
	}
	return ok
}

func randomstr() string {
	l := rand.Intn(10) + 1
	a := make([]rune, l)
	for i := 0; i < l; i++ {
		a[i] = rune(rand.Intn(26) + 'a')
	}
	return string(a)
}
