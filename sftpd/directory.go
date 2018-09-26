package sftpd

import (
	"github.com/taruti/sftpd"
	"io"
)

type Directory struct {
	na *[]sftpd.NamedAttr
	n int
}

func (this *Directory) Readdir(count int) ([]sftpd.NamedAttr, error) {
	//fmt.Println(this, len(*this.na), (*this).n)
	//tmp := []sftpd.NamedAttr{}
	il := len(*this.na)
	if il == 0 || this.n >= il {
		return nil, io.EOF
	}
	i := this.n + count
	if i >= il {
		i = il
	}
	s := this.n
	this.n = i
	return (*this.na)[s:i], nil

	//i := this.n
	//for ; i < this.n+count; i++ {
	//	if i >= l {
	//		(*this).n = i
	//		return tmp, io.EOF
	//	}
	//	tmp = append(tmp, (*this.na)[i])
	//}
	//(*this).n = i
	//return tmp, nil
}

func (this *Directory) Close() error {
	return nil
}