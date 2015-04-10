package main

import (
	"github.com/yaozijian/uart"
)

type (
	mathmsg struct {
		A        int
		B        int
		Result   int
		Operator operator
	}
	echomsg  string
	operator int
)

const (
	operator_add = iota
	operator_del
	operator_mul
	operator_div
)

func initTable() *uart.TypeTable {
	table := uart.NewTypeTable()
	table.RegisterType(mathmsg{})
	table.RegisterType(echomsg(""))
	return table
}

func (this operator) String() string {
	names := []string{"+", "-", "*", "/"}
	return names[int(this)]
}
