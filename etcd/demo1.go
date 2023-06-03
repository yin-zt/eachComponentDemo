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
	_, _ = WriteEtcd(fmt.Sprintf("%s%s", config.EtcdServer, "/keeper/hello1"), "world2", "10", etcdToken)
	//fmt.Println(result)
	//fmt.Println(err)
	//etcdClent := c
	//kapi.Get("")
	GetEtcdNode(fmt.Sprintf("%s%s", config.EtcdServer, "/keeper/hello1"), etcdToken)
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

func GetEtcdNode(url, token string) {
	fmt.Println(url, token)

	// 创建 HTTP 请求对象
	req := httplib.Get(url)

	// 设置用户名和密码
	req.Header("Authorization", token)

	// 发送请求并获取响应
	resp, err := req.Response()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != 200 {
		fmt.Println("Error: Unexpected status code", resp.StatusCode)
		return
	}

	// 解析响应体中的数据
	data := make(map[string]interface{})
	err = req.ToJSON(&data)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("11111111111111")
	fmt.Println(data)

	// 获取节点值
	//value := data["node"].(map[string]interface{})["value"].(string)
	//fmt.Println("Node value:", value)

}
