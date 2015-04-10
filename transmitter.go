package uart

import (
	"io"
	"sync"
)

type (
	Transmitter struct {
		cfg      TransmitterConfig
		waitstop sync.WaitGroup
		stop     chan int
		stoponce sync.Once
	}

	TransmitterConfig struct {
		Table    *TypeTable
		Protocol Protocol
		//---
		Src <-chan interface{}
		Dst io.Writer
		//---
		OnError func(interface{}, error) bool
		Ctx     interface{}
	}
)

func NewTransmitter(cfg *TransmitterConfig) (this *Transmitter) {

	if cfg == nil || cfg.Table == nil || cfg.Protocol == nil {
		return
	} else if cfg.Src == nil || cfg.Dst == nil {
		return
	}

	this = &Transmitter{cfg: *cfg, stop: make(chan int)}

	this.waitstop.Add(1)
	go this.loop()

	return
}

func (this *Transmitter) Stop(wait bool) {
	this.stoponce.Do(func() { close(this.stop) })
	if wait {
		this.waitstop.Wait()
	}
}

func (this *Transmitter) loop() {

	defer func() {
		if closer, _ := this.cfg.Dst.(io.Closer); closer != nil {
			closer.Close()
		}
		this.waitstop.Done()
	}()

	cfg := &this.cfg
	enc := cfg.Protocol.NewEncoder(cfg.Dst)

	var msg *Msg
	var err error

	for err == nil {
		select {
		case <-this.stop:
			return
		case input, ok := <-cfg.Src:
			if ok {
				msg, err = cfg.Table.NewMessage(input, cfg.Protocol)
				if err == nil {
					err = enc.Encode(msg)
				}
				if err != nil && cfg.OnError != nil {
					if cfg.OnError(cfg.Ctx, err) {
						err = nil //错误回调函数表示忽略错误,继续运行
					}
				}
			} else {
				return
			}
		}
	}
}
