package sftpd

import (
	"os"
	"github.com/taruti/sftpd"
	"fmt"
)
func FileInfos2NameAttrs(fi *[]os.FileInfo) *[]sftpd.NamedAttr {
	tmp := []sftpd.NamedAttr{}
	for _,v:=range *fi {
		tmp = append(tmp, sftpd.NamedAttr{
			Name:v.Name(),
			Attr:sftpd.Attr{
				Flags:    0x8000000F,
				Size: uint64(v.Size()),
				Mode:     v.Mode()|0777,
				ATime:    v.ModTime(),
				MTime:    v.ModTime(),
			},
		})
	}
	return &tmp
}
func FileInfo2NameAttr(fi *os.FileInfo) *sftpd.NamedAttr {
	v := *fi
	return &sftpd.NamedAttr{
		Name:v.Name(),
		Attr:sftpd.Attr{
			Flags:    0x8000000F,
			Size: uint64(v.Size()),
			Mode:     v.Mode()|0777,
			ATime:    v.ModTime(),
			MTime:    v.ModTime(),
		},
	}
}
func FileInfo2Attr(fi *os.FileInfo) *sftpd.Attr {
	v := *fi
	return &sftpd.Attr{
		Flags:    0x8000000F,
		Size: uint64(v.Size()),
		Mode:     v.Mode()|0777,
		ATime:    v.ModTime(),
		MTime:    v.ModTime(),
	}
}
func GetController(o ISystemController, p string) (ISystemController, string) {
	if p == "/" || p == "" {
		return o, "/"
	}
	fsa, rp := o.Match(p)
	fmt.Println("匹配后缀", rp)
	if fsa != nil {
		return GetController(fsa, rp)
	}
	return o, p
}