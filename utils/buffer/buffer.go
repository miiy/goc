package hhosai

import (
	"errors"
	"fmt"
	"io"
	"sync"
)

type Buffer struct {
	sync.Mutex
	buf []interface{}
	off int
}

func (b *Buffer) Write(i interface{}) error {
	fmt.Println("Buffer write()")
	if i == nil {
		return errors.New("cannot write nil value")
	}

	b.Lock()
	defer b.Unlock()
	b.buf = append(b.buf, i)
	return nil
}

func (b *Buffer) Read() (interface{}, error) {
	b.Lock()
	defer b.Unlock()
	if len(b.buf) == 0 || b.off >= len(b.buf) {
		return nil, io.EOF
	}

	data := b.buf[b.off]
	b.off++
	return data, nil
}

func (b *Buffer) Reset() error {
	b.Lock()
	defer b.Unlock()
	b.buf = nil
	b.off = 0
	return nil
}
