package main

import (
	"./sftpd"
	"fmt"
)

func main() {

	s, err := sftpd.NewSftpd("127.0.0.1:2323")

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("服务开启")

	tmp := sftpd.NewMapController("/公共文件", `E:`)

	s.PublicController(tmp)

	us := s.AddUserPasswd("test", "1234")
	us.Controller(sftpd.NewMapController("/成员组文件", `E:`))

	mb := sftpd.NewMapController("/b", `E:`)
	mb.SetPermission(0xEF, false)
	us.Controller(mb)
	us.Controller(sftpd.NewMapController("/c", `E:`))

	//us1 := s.AddUserPasswd("abc", "123")
	//us1.Controller(sftpd.NewMapController("/成员组文件", `E:`))
	//us1.Controller(sftpd.NewMapController("/b", `D:`))

	s.Listen()

}