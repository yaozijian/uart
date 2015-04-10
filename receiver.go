package uart

import (
	"io"
	"sync"
	"sync/atomic"
)

type (
	Receiver struct {
		cfg      ReceiverConfig
		waitstop sync.WaitGroup
		stop     chan int
		stoping  int32
		stoponce sync.Once
	}

	ReceiverConfig struct {
		Table    *TypeTable
		Protocol Protocol
		Policy   ReceiverPolicy
		//---
		Src io.Reader
		Dst chan<- interface{}
		//---
		OnError func(interface{}, error) bool
		Ctx     interface{}
	}

	ReceiverPolicy int
)

const (
	Receive_sync = iota
	Receive_discard
)

func NewReceiver(cfg *ReceiverConfig) (this *Receiver) {

	if cfg == nil || cfg.Table == nil || cfg.Protocol == nil {
		return
	} else if cfg.Src == nil || cfg.Dst == nil {
		return
	}

	this = &Receiver{cfg: *cfg, stop: make(chan int)}

	this.waitstop.Add(1)
	go this.loop()

	return
}

func (this *Receiver) Stop(wait bool) {
	this.stoponce.Do(func() { close(this.stop); atomic.StoreInt32(&this.stoping, 1) })
	if wait {
		this.waitstop.Wait()
	}
}

func (this *Receiver) loop() {

	defer this.waitstop.Done()
	defer close(this.cfg.Dst)

	cfg := &this.cfg

	// 如果源支持关闭,则等到了停止请求时就关闭源，
	// 这样阻塞的dec.Decode()调用就会返回
	if closer, _ := cfg.Src.(io.Closer); closer != nil {
		go func() {
			<-this.stop
			closer.Close()
		}()
	}

	var err error
	var obj interface{}

	dec := cfg.Protocol.NewDecoder(cfg.Src)

	for err == nil && atomic.LoadInt32(&this.stoping) == 0 {
		msg := &Msg{}
		if err = dec.Decode(msg); err == nil {
			if obj, err = msg.Decode(cfg.Table, cfg.Protocol); err == nil {
				this.onmessage(obj)
			}
		}
		if err != nil && cfg.OnError != nil {
			if cfg.OnError(cfg.Ctx, err) {
				err = nil //错误回调函数表示忽略错误,继续运行
			}
		}
	}
}

func (this *Receiver) onmessage(val interface{}) {
	switch this.cfg.Policy {
	case Receive_sync:
		select {
		case this.cfg.Dst <- val:
		case <-this.stop:
		}
	case Receive_discard:
		select {
		case this.cfg.Dst <- val:
		default:
		}
	}
}
