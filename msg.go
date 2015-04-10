package uart

import (
	"bytes"
	"fmt"
)

type (
	Msg struct {
		Id   TypeId
		Body []byte
	}
)

func (msg *Msg) Decode(table *TypeTable, protocol Protocol) (obj interface{}, err error) {

	obj = table.NewInstance(msg.Id)
	if obj == nil {
		err = fmt.Errorf("Type not registered")
		return
	}

	buf := bytes.NewBuffer(msg.Body)
	dec := protocol.NewDecoder(buf)
	if err = dec.Decode(obj); err != nil {
		obj = nil
	}

	return
}
