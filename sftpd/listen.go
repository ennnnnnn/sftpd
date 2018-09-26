package sftpd

import (
	"log"
	"net"

	"github.com/taruti/sftpd"
	"golang.org/x/crypto/ssh"
)

type Sftpd struct {
	host string
	upass *UserPasswd
	config *ssh.ServerConfig
	public ISystemController
}

func NewSftpd(host string) (*Sftpd, error)  {
	tmp := Sftpd{host:host}
	tmp.upass = &UserPasswd{us: map[string]*UserSpace{}, sftpd: &tmp}
	conf, err := NewConfig(tmp.upass)
	if err!= nil {
		return nil, err
	}
	tmp.config = conf
	return &tmp, nil
}

func (this *Sftpd)PublicController(v ISystemController) {
	v.SetName("/公共文件")
	this.public = v
}

func (this *Sftpd)AddUserPasswd(name string, passwd string) *UserSpace {
	us := this.upass.UserSpace(name)
	us.passwd = []byte(passwd)
	return us
}

func (this *Sftpd)Listen() error {
	listener, e := net.Listen("tcp", this.host)
	if e != nil {
		return e
	}
	for {
		conn, e := listener.Accept()
		if e != nil {
			return e
		}
		go this.HandleConnect(conn)
	}
}

func (this *Sftpd)HandleConnect(conn net.Conn) error {
	sc, chans, reqs, e := ssh.NewServerConn(conn, this.config)
	if e != nil {
		return e
	}
	defer sc.Close()

	// The incoming Request channel must be serviced.
	go PrintDiscardRequests(reqs)

	// Service the incoming Channel channel.
	for newChannel := range chans {
		if newChannel.ChannelType() != "session" {
			newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
			continue
		}
		channel, requests, err := newChannel.Accept()
		if err != nil {
			return err
		}

		fs := this.upass.GetUserSpace(sc.Permissions.CriticalOptions["user"]).fs

		go func(in <-chan *ssh.Request) {
			for req := range in {
				ok := false
				switch {
				case sftpd.IsSftpRequest(req):
					ok = true
					go func() {
						e := sftpd.ServeChannel(channel, fs)
						if e != nil {
							log.Println("sftpd servechannel failed:", e)
						}
					}()
				}
				req.Reply(ok, nil)
			}
		}(requests)

	}
	return nil
}

func PrintDiscardRequests(in <-chan *ssh.Request) {
	for req := range in {
		log.Println("Discarding ssh request", req.Type, *req)
		if req.WantReply {
			req.Reply(false, nil)
		}
	}
}
