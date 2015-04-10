package uart

import (
	"bytes"
	"fmt"
	"reflect"
	"sync"
)

type (
	TypeId    int64
	TypeTable struct {
		typeTable map[TypeId]*typeItem
		nextId    TypeId
		sync.RWMutex
	}
	typeItem struct {
		t reflect.Type
	}
)

func NewTypeTable() *TypeTable {
	return &TypeTable{typeTable: make(map[TypeId]*typeItem)}
}

func (table *TypeTable) RegisterType(val interface{}) TypeId {
	item := &typeItem{t: reflect.TypeOf(val)}
	table.Lock()
	defer table.Unlock()
	table.nextId++
	table.typeTable[table.nextId] = item
	return table.nextId
}

func (table *TypeTable) RemoveType(val interface{}) {

	table.Lock()
	defer table.Unlock()

	t := reflect.TypeOf(val)
	for id, item := range table.typeTable {
		if isSameType(item.t, t) {
			delete(table.typeTable, id)
			break
		}
	}
}

func (table *TypeTable) NewInstance(id TypeId) interface{} {

	table.RLock()
	defer table.RUnlock()

	if item := table.typeTable[id]; item != nil {
		v := reflect.New(item.t)
		if v.CanInterface() {
			return v.Interface()
		}
	}

	return nil
}

func (table *TypeTable) NewMessage(val interface{}, protocol Protocol) (*Msg, error) {

	table.RLock()
	defer table.RUnlock()

	t := reflect.TypeOf(val)
	for id, item := range table.typeTable {
		sametype := isSameType(item.t, t)
		sametype = sametype || (t.Kind() == reflect.Ptr && isSameType(t.Elem(), item.t))
		if sametype {
			buf := bytes.NewBuffer(nil)
			enc := protocol.NewEncoder(buf)
			if err := enc.Encode(val); err == nil {
				return &Msg{Id: id, Body: buf.Bytes()}, nil
			} else {
				return nil, err
			}
		}
	}

	return nil, fmt.Errorf("Type not registered")
}

func isSameType(a, b reflect.Type) bool {
	return a.PkgPath() == b.PkgPath() && a.Name() == b.Name()
}
