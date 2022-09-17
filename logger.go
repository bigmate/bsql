package bsql

import (
	"fmt"
	"io"
)

type Logger interface {
	Errorf(format string, args ...interface{})
}

type logger struct {
	io.Writer
}

func (l logger) Errorf(format string, args ...interface{}) {
	fmt.Fprintf(l, format, args...)
	fmt.Fprintln(l)
}
