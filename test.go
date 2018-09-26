package main

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"time"
	"fmt"
	"io/ioutil"
)

func MD5(text string) string{
	ctx := md5.New()
	ctx.Write([]byte(text))
	return hex.EncodeToString(ctx.Sum(nil))
}

// return len=8  salt
func GetRandomSalt() string {
	return GetRandomString(8)
}

//生成随机字符串
func GetRandomString(ll int) string{
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < ll; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func main()  {

	fi,_:= ioutil.ReadDir(`E:\`)
	for _,v:=range fi{
		fmt.Println(v.Name())
	}
	fmt.Println(48&^16)


	//for n:=0;n<1500;n++{
	//	f, err := os.OpenFile(`E:\a\`+GetRandomSalt()+".txt", os.O_CREATE|os.O_SYNC, 0)
	//	if err!=nil {
	//		continue
	//	}
	//	f.Write([]byte("hello world!"))
	//	f.Close()
	//}

	//a := "123456789"
	//
	//fmt.Println(len(a), a[0:3], a[3:9])
	//fmt.Println(strings.Replace(p, `\`, `/`, -1))
}
