package sftpd

import (
	"github.com/taruti/sftpd"
	"strings"
	"syscall"
	"fmt"
	"reflect"
	"os"
)

type SystemController struct {
	cl *CustomList
	EasyController
}
func NewSystemController() *SystemController {
	return &SystemController{
		cl:             NewCustomList(),
		EasyController: CreateEasyController(nil, "/", 0x50, ^0x50),
	}
}
func (this *SystemController) Controller(v ISystemController) {
	fmt.Println("传入", v.GetName(), v.GetAccess())
	if v == nil {
		fmt.Println("成员为空")
		return
	}
	this.EasyController.Controller(v)
	this.cl.Clear()
	for _, v := range this.EasyController.l {
		fmt.Println("迭代", v.GetName(), v.GetAccess())
		if v.GetAccess()&PERM_READ_NAME==0 {
			continue
		}
		name := v.GetName()
		name = name[1:]
		if strings.IndexByte(name, '/') >= 0 {
			continue
		}
		stat, err := v.Stat("/", true)
		if err != nil {
			continue
		}
		this.cl.Add(sftpd.NamedAttr{
			Name: name,
			Attr: *stat,
		})
	}
}
func (this *SystemController) Stat(name string, islstat bool) (*sftpd.Attr, error) {
	fmt.Println("文件系统Stat", name)
	name = strings.Replace(name, `\`, `/`, -1)
	if name!="/" {
		c, p := GetController(this, name)
		fmt.Println("文件系统Stat子", this != c)
		if this != c {
			fmt.Println("文件系统Stat子", this.GetName(), c.GetName(), reflect.TypeOf(c))
			return c.Stat(p, islstat)
		}
		return nil, syscall.ERROR_FILE_NOT_FOUND
	}
	return &sftpd.Attr{}, nil
}
func (this *SystemController) OpenDir(name string) (sftpd.Dir, error) {
	fmt.Println("文件系统OpenDir", name)
	name = strings.Replace(name, `\`, `/`, -1)
	if name != "/" {
		c, p := GetController(this, name)
		if this != c {
			return c.OpenDir(p)
		}
		return nil, syscall.ERROR_FILE_NOT_FOUND
	}
	if this.GetAccess()&PERM_READ_DATA==0 {
		return nil, os.ErrPermission
	}
	return this.cl, nil
}
func (this *SystemController) OpenFile(name string, flags uint32, attr *sftpd.Attr) (sftpd.File, error) {
	fmt.Println("文件系统Stat", name)
	name = strings.Replace(name, `\`, `/`, -1)
	if name != "/" {
		c, p := GetController(this, name)
		if this != c {
			return c.OpenFile(p, flags, attr)
		}
	}
	return nil, syscall.ERROR_FILE_NOT_FOUND
}
func (this *SystemController) Remove(name string) error {
	name = strings.Replace(name, `\`, `/`, -1)
	if name != "/" {
		c, p := GetController(this, name)
		if this != c {
			return c.Remove(p)
		}
	}
	return os.ErrPermission
}
func (this *SystemController) Mkdir(name string, attr *sftpd.Attr) error {
	name = strings.Replace(name, `\`, `/`, -1)
	if name != "/" {
		c, p := GetController(this, name)
		if this != c {
			return c.Mkdir(p, attr)
		}
	}
	return os.ErrPermission
}
func (this *SystemController) Rmdir(name string) error {
	name = strings.Replace(name, `\`, `/`, -1)
	if name != "/" {
		c, p := GetController(this, name)
		if this != c {
			return c.Rmdir(p)
		}
	}
	return os.ErrPermission
}
func (this *SystemController) SetStat(name string, attr *sftpd.Attr) error {
	name = strings.Replace(name, `\`, `/`, -1)
	if name != "/" {
		c, p := GetController(this, name)
		if this != c {
			return c.SetStat(p, attr)
		}
	}
	return os.ErrPermission
}
func (this *SystemController) MyRename(oo ISystemController, op, np string, flags uint32) error {
	return os.ErrPermission
}
//func (this FileSystem) ReadLink(path string) (string, error) {
//	return this.r.ReadLink(path)
//}
//func (this FileSystem) CreateLink(path string, target string, flags uint32) error {
//	return this.r.CreateLink(path, target, flags)
//}
