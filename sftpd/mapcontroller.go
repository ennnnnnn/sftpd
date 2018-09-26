package sftpd

import (
	"github.com/taruti/sftpd"
	"syscall"
	"io/ioutil"
	"strings"
	"os"
	"fmt"
	"reflect"
)

type MapController struct {
	path string
	EasyController
}
func NewMapController(n, p string) *MapController {
	return CreateMapController(nil, n, p, 0xD0, 0)
}
func CreateMapController(pen ISystemController, n, p string, pa, ba int32) *MapController {
	if p[len(p)-1] == ':' {
		p = p + "/"
	}
	return &MapController{
		path:           p,
		EasyController: CreateEasyController(pen, n, pa, ba),
	}
}
func (this *MapController) OpenDir(name string) (sftpd.Dir, error) {
	fmt.Println("目录OpenDir", name)
	name = strings.Replace(name, `\`, `/`, -1)
	if name != "/" {
		c, p := GetController(this, name)
		if this != c {
			return c.OpenDir(p)
		}
		name = this.path + name
	} else {
		name = this.path
	}
	if this.GetAccess()&PERM_READ_DATA==0 {
		return nil, os.ErrPermission
	}
	dir, err := ioutil.ReadDir(name)
	if err != nil {
		return nil, syscall.ERROR_FILE_NOT_FOUND
	}
	return &Directory{na: FileInfos2NameAttrs(&dir)}, nil
}
func (this *MapController) OpenFile(name string, flags uint32, attr *sftpd.Attr) (sftpd.File, error){
	fmt.Println("目录OpenFile", name)
	name = strings.Replace(name, `\`, `/`, -1)
	if name != "/" {
		c, p := GetController(this, name)
		if this != c {
			return c.OpenFile(p, flags, attr)
		}
		at := this.GetAccess()
		if at&PERM_READ_DATA==0 && at&PERM_WRITE_DATA==0 {
			return nil, os.ErrPermission
		}
		name := this.path + name
		f, err := os.OpenFile(name, os.O_CREATE, 0777)
		if err != nil {
			return nil, os.ErrNotExist
		}
		return &File{name, f}, nil
	}
	return nil, os.ErrNotExist
}
func (this *MapController) Stat(name string, islstat bool) (*sftpd.Attr, error) {
	fmt.Println("目录Stat", name)
	name = strings.Replace(name, `\`, `/`, -1)
	if name != "/" {
		c, p := GetController(this, name)
		if this != c {
			return c.Stat(p, islstat)
		}
		name = this.path + name
	} else {
		name = this.path
	}
	stat, err := os.Stat(name)
	if err != nil {
		return nil, syscall.ERROR_FILE_NOT_FOUND
	}
	return FileInfo2Attr(&stat), nil
}
func (this *MapController) SetStat(name string, attr *sftpd.Attr) error {
	fmt.Println("目录SetStat", name)
	name = strings.Replace(name, `\`, `/`, -1)
	if name != "/" {
		c, p := GetController(this, name)
		if this != c {
			return c.SetStat(p, attr)
		}
		name = this.path + name
	} else {
		name = this.path
	}
	if this.GetAccess()&PERM_WRITE_DATA==0 {
		return os.ErrPermission
	}
	return os.Chtimes(name, attr.ATime, attr.MTime)
}
func (this *MapController) Mkdir(name string, attr *sftpd.Attr) error {
	fmt.Println("目录Mkdir", name)
	name = strings.Replace(name, `\`, `/`, -1)
	if name != "/" {
		c, p := GetController(this, name)
		if this != c {
			return c.Mkdir(p, attr)
		}
		name = this.path + name
	} else {
		name = this.path
	}
	if this.GetAccess()&PERM_CREATE_DIR==0 {
		return os.ErrPermission
	}
	return os.Mkdir(name, 0777)
}
func (this *MapController) Rmdir(name string) error {
	fmt.Println("目录Rmdir", name)
	name = strings.Replace(name, `\`, `/`, -1)
	if name != "/" {
		c, p := GetController(this, name)
		if this != c {
			return c.Rmdir(p)
		}
		name = this.path + name
	} else {
		name = this.path
	}
	if this.GetAccess()&PERM_DELETE_DIR==0 {
		return os.ErrPermission
	}
	return os.Remove(name)
}
func (this *MapController) Remove(name string) error {
	fmt.Println("目录Remove", name)
	name = strings.Replace(name, `\`, `/`, -1)
	if name != "/" {
		c, p := GetController(this, name)
		if this != c {
			return c.Remove(p)
		}
		name = this.path + name
	} else {
		name = this.path
	}
	fmt.Println(this.GetPermission(), this.GetBan(), this.GetAccess(), this.GetAccess()&PERM_DELETE_FILE)
	if this.GetAccess()&PERM_DELETE_FILE==0 {
		return os.ErrPermission
	}
	return os.Remove(name)
}
func (this *MapController) MyRename(oo ISystemController, op, np string, flags uint32) error {
	if reflect.TypeOf(this)!=reflect.TypeOf(oo) {
		return os.ErrPermission
	}
	tmp := oo.(*MapController)
	name1 := this.path+op
	name2 := tmp.path+np
	fmt.Println("目录MyRename", name1, name2)
	if this.GetAccess()&PERM_WRITE_NAME==0 {
		return os.ErrPermission
	}
	return os.Rename(name1, name2)
}