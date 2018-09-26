package sftpd

import (
	"sort"
	"sync"
	"github.com/taruti/sftpd"
	"path/filepath"
	"strings"
	"fmt"
	"os"
)

const (
	PERM_CREATE_DIR = 1
	PERM_CREATE_FILE = 2
	PERM_DELETE_DIR = 4
	PERM_DELETE_FILE = 8
	PERM_READ_NAME = 16
	PERM_WRITE_NAME = 32
	PERM_READ_DATA = 64
	PERM_WRITE_DATA = 128
)

type ISystemController interface {
	sftpd.FileSystem
	SetName(string)
	GetName() string
	Controller(ISystemController)
	Match(string) (ISystemController, string)
	MyRename(ISystemController, string, string, uint32) error
	GetPermission() int32
	GetBan() int32
	SetPermission(int32, bool)
	SetBan(int32)
	setParent(ISystemController)
	GetParent() ISystemController
	GetAccess() int32
}
type ISystemControllers []ISystemController
func (this ISystemControllers) Len() int {
	return len(this)
}
func (this ISystemControllers) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}
func (this ISystemControllers) Less(i, j int) bool {
	return len(this[i].GetName()) > len(this[j].GetName())
}

type EasyController struct {
	parent ISystemController
	name string
	p int32
	b int32
	l ISystemControllers
	sftpd.EmptyFS
	mx sync.RWMutex
}
func NewEasyController(name string) EasyController {
	return EasyController{
		name: name,
	}
}
func CreateEasyController(pen ISystemController, n string, p, b int32) EasyController {
	return EasyController{
		parent: pen,
		name:   n,
		p:      p,
		b:      b,
	}
}
func (this *EasyController) setParent(v ISystemController)  {
	this.parent = v
}
func (this *EasyController) GetParent() ISystemController {
	return this.parent
}
func (this *EasyController) GetName() string  {
	return this.name
}
func (this *EasyController) SetName(name string)  {
	this.name = name
}
func (this *EasyController) Match(p string) (ISystemController, string) {
	this.mx.RLock()
	defer this.mx.RUnlock()
	p = p + "/"
	l := len(p)
	for _, v := range this.l {
		nt := v.GetName()+"/"
		tl := len(nt)
		if l < tl {
			continue
		}
		if p == nt {
			return v, "/"
		}
		if nt == p[:tl] {
			return v, "/" + p[tl:len(p)-1]
		}
	}
	return nil, p
}
func (this *EasyController) Controller(v ISystemController) {
	this.mx.Lock()
	defer this.mx.Unlock()
	v.setParent(this)
	this.l = append(this.l, v)
	sort.Sort(this.l)
}
func (this *EasyController) RealPath(path string) (string, error) {
	path = filepath.Clean(path)
	path = strings.Replace(path, "\\", "/", -1)
	fmt.Println("文件RealPath" , path)
	return path, nil
}
func (this *EasyController) Rename(old string, new string, flags uint32) error {
	on := strings.Replace(old, `\`, `/`, -1)
	c1, p1 := GetController(this, on)
	nn := strings.Replace(new, `\`, `/`, -1)
	c2, p2 := GetController(this, nn)
	return c2.MyRename(c1, p1, p2, flags)
}
func (this *EasyController) MyRename(oo ISystemController, op, np string, flags uint32) error {
	return os.ErrPermission
}
func (this *EasyController) GetPermission() int32 {
	return this.p
}
func (this *EasyController) SetPermission(v int32, e bool) {
	this.p = v
	if !e {
		this.b = ^v
	}
}
func (this *EasyController) GetBan() int32 {
	return this.b
}
func (this *EasyController) SetBan(v int32) {
	this.b = v
}
func (this *EasyController) GetAccess() int32 {
	if this.parent == nil {
		return this.p &^ this.b
	} else {
		return (this.parent.GetAccess() | this.p) &^ this.b
	}
}