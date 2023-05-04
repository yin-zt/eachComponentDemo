package main

import (
	"encoding/base64"
	"fmt"
	"github.com/astaxie/beego/httplib"
	"github.com/yin-zt/eachComponentDemo/config"
	"log"
	"time"
)

var (
	str = config.EtcdUser + ":" + config.EtcdPass // 加密调用
)

func main() {
	etcdToken := "Basic " + Base64Encode(str)
	result, err := WriteEtcd(fmt.Sprintf("%s%s", config.EtcdServer, "/keeper/hello"), "world", "10", etcdToken)
	fmt.Println(result)
	fmt.Println(err)
	//etcdClent := c
	//kapi.Get("")
}

func Base64Encode(str string) string {

	return base64.StdEncoding.EncodeToString([]byte(str))
}

func WriteEtcd(url string, value string, ttl string, token string) (string, error) {

	req := httplib.Post(url)

	req.Header("Authorization", token)
	req.Param("value", value)
	req.Param("ttl", ttl)
	req.SetTimeout(time.Second*10, time.Second*5)
	str, err := req.String()
	//	fmt.Println(str)
	if err != nil {
		print(err)
		log.Fatalln(err)
	}
	return str, err
}
