package sftpd

import (
	"github.com/taruti/sshutil"
	"golang.org/x/crypto/ssh"
	"github.com/pkg/errors"
	"crypto/subtle"
	"fmt"
)

type UserSpace struct {
	name string
	passwd []byte
	fs ISystemController
}

func (this *UserSpace)SetFileSystem(fs ISystemController) {
	this.fs = fs
}

func (this *UserSpace)Controller(v ISystemController) {
	this.fs.Controller(v)
}

type UserPasswd struct {
	sftpd *Sftpd
	us map[string]*UserSpace
}

func (this *UserPasswd)GetUserSpace(name string) *UserSpace {
	if _, ok := this.us[name]; ok {
		return this.us[name]
	}
	return nil
}

func (this *UserPasswd)UserSpace(name string) *UserSpace {
	if _, ok := this.us[name]; !ok {
		sc := NewSystemController()
		sc.Controller(this.sftpd.public)
		fmt.Println(this.sftpd.public)
		tmp := UserSpace{
			name:name,
			passwd: []byte(""),
			fs: sc,
		}
		this.us[name] = &tmp
	}
	return this.us[name]
}

func (this *UserPasswd)RemUserSpace(name string) {
	delete(this.us, name)
}

func (this *UserPasswd)passwordCheck(conn ssh.ConnMetadata, passwd []byte) (*ssh.Permissions, error) {
	us := this.GetUserSpace(conn.User())
	if us == nil || subtle.ConstantTimeCompare(us.passwd, passwd) != 1 {
		return nil, errors.New("用户密码不存在")
	}
	return &ssh.Permissions{
		CriticalOptions: map[string]string{
			"user": us.name,
		},
	}, nil
}

func NewConfig(upass *UserPasswd) (*ssh.ServerConfig, error) {
	config := &ssh.ServerConfig{
		PasswordCallback: upass.passwordCheck,
	}
	hkey, e := sshutil.KeyLoader{Flags: sshutil.Create | sshutil.RSA2048}.Load()
	if e != nil {
		return nil, e
	}
	config.AddHostKey(hkey)
	return config, nil
}

