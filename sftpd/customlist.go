package sftpd

import (
	"github.com/taruti/sftpd"
	"io"
	"sync"
)

type CustomList struct {
	l []sftpd.NamedAttr
	n int
	mx sync.RWMutex
}

func NewCustomList() *CustomList {
	return &CustomList{
		l: []sftpd.NamedAttr{},
		n: 0,
	}
}

func (this *CustomList) Add(a sftpd.NamedAttr) {
	this.mx.Lock()
	defer this.mx.Unlock()
	this.l = append(this.l, a)
}

func (this *CustomList) Clear() {
	this.mx.Lock()
	defer this.mx.Unlock()
	this.l = []sftpd.NamedAttr{}
	this.n = 0
}

func (this *CustomList) Readdir(count int) ([]sftpd.NamedAttr, error) {
	this.mx.RLock()
	defer this.mx.RUnlock()
	tmp := []sftpd.NamedAttr{}
	l := len(this.l)
	if l == 0 || this.n >= l {
		return nil, io.EOF
	}
	for i := this.n; i < this.n+count; i++ {
		if i >= l {
			this.n = i
			return tmp, nil
		}
		tmp = append(tmp, this.l[i])
	}
	return tmp, nil
}

func (this *CustomList) Close() error {
	this.n = 0
	return nil
}