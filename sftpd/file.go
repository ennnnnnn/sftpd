package sftpd

import (
	"github.com/taruti/sftpd"
	"os"
)

type File struct {
	name string
	f *os.File
}
func (this *File)FStat() (*sftpd.Attr, error){
	fileInfo, err := os.Stat(this.name)
	if err != nil {
		return nil, err
	}
	return FileInfo2Attr(&fileInfo), nil
}
func (this *File)FSetStat(*sftpd.Attr) error{
	return nil
}
func (this *File)Close() error{
	return this.f.Close()
}
func (this *File)ReadAt(p []byte, off int64) (n int, err error) {
	return this.f.ReadAt(p, off)
}
func (this *File)WriteAt(p []byte, off int64) (n int, err error){
	return this.f.WriteAt(p, off)
}