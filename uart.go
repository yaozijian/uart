package uart

import (
	"io"
	"sync"
)

type (
	Uart struct {
		cfg      UartConfig
		waitstop sync.WaitGroup
		stop     chan int
		stoponce sync.Once
		//--
		SendChnl    chan<- interface{}
		RecvChnl    <-chan interface{}
		receiver    *Receiver
		transmitter *Transmitter
	}

	UartConfig struct {
		Table    *TypeTable
		Protocol Protocol
		Policy   ReceiverPolicy
		//---
		Src     io.Reader
		Dst     io.Writer
		SendBuf int
		RecvBuf int
	}
)

func NewUart(cfg *UartConfig) (this *Uart) {

	if cfg == nil || cfg.Table == nil || cfg.Protocol == nil {
		return
	} else if cfg.Src == nil || cfg.Dst == nil || cfg.SendBuf < 0 || cfg.RecvBuf < 0 {
		return
	}

	this = &Uart{cfg: *cfg, stop: make(chan int)}

	sendChnl := make(chan interface{}, cfg.SendBuf)
	recvChnl := make(chan interface{}, cfg.RecvBuf)

	recvCfg := &ReceiverConfig{
		Table:    cfg.Table,
		Protocol: cfg.Protocol,
		Policy:   cfg.Policy,
		Src:      cfg.Src,
		Dst:      recvChnl,
	}
	receiver := NewReceiver(recvCfg)

	sendCfg := &TransmitterConfig{
		Table:    cfg.Table,
		Protocol: cfg.Protocol,
		Src:      sendChnl,
		Dst:      cfg.Dst,
	}
	sender := NewTransmitter(sendCfg)

	this.receiver = receiver
	this.transmitter = sender
	this.SendChnl = sendChnl
	this.RecvChnl = recvChnl

	return
}

func (this *Uart) Stop(wait bool) {
	this.stoponce.Do(func() {
		this.receiver.Stop(false)
		this.transmitter.Stop(false)
	})
	if wait {
		waitstop := make(chan int, 2)
		go func() {
			this.receiver.Stop(true)
			waitstop <- 1
		}()
		go func() {
			this.transmitter.Stop(true)
			waitstop <- 1
		}()
		<-waitstop
		<-waitstop
		close(waitstop)
	}
}
