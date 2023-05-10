package main

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/yin-zt/eachComponentDemo/config"
	"time"
)

var (
	resp any
)

func initRedis() (dial redis.Conn, err error) {
	dial, err = redis.Dial("tcp", config.RedisHost)
	if err != nil {
		resp = err
		panic(resp)
	}
	return dial, nil
}

func testSetGet(dial redis.Conn) {
	// 设置一个Key和value
	dial.Do("set", "abc", "this is a test")

	// 设置过期时间
	dial.Do("expire", "abc", 2)

	// 读取指定key的value
	reply, err := dial.Do("get", "abc")
	str, err := redis.String(reply, err)
	if err != nil {
		resp = err
		panic(resp)
	}
	fmt.Println(str)
}

func testHSetGet(dial redis.Conn) {
	key := "abc"
	value := "this is a book test"

	// 设置一个hash表来存储
	reply, err := dial.Do("hset", "books", key, value)
	if err != nil {
		resp = err
		panic(resp)
	}
	fmt.Println(reply)
	reply2, err := dial.Do("hget", "books", key)
	str, err := redis.String(reply2, err)
	if err != nil {
		resp = err
		panic(resp)
	}
	fmt.Println(str)
}

func testMSetGet(dial redis.Conn) {
	key := "abc"
	value := "this is a test book"
	key2 := "golang"
	value2 := "this is a golang book"

	// 一次设置多个key
	dial.Do("mset", "books", key, value, key2, value2)
	str, err := redis.Strings(dial.Do("mget", "books", key, key2))
	if err != nil {
		resp = err
		panic(resp)
	}
	fmt.Println(str)
}

func main() {
	dial, _ := initRedis()
	testSetGet(dial)
	testHSetGet(dial)
	testMSetGet(dial)

	time.Sleep(4 * time.Second)
	reply, err := dial.Do("get", "abc")
	str, err := redis.String(reply, err)
	if err != nil {
		resp = err
		panic(resp)
	}
	fmt.Println("get abc", str)

	defer dial.Close()
}
