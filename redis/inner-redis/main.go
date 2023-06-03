package main

import (
	"fmt"
	"github.com/alicebob/miniredis"
	"github.com/yin-zt/eachComponentDemo/config"
	"log"
	"strings"
	"time"
)

var (
	action = "init"
	resp   any
)

// main 起了一个内建的redis
func main() {

	// 在给定时间内持续清理redis中TTL过期的key
	CleanExpireKeys := func(redis *miniredis.Miniredis) {
		for {
			t := time.Tick(time.Second * 60)
			<-t
			redis.FastForward(time.Second * 60)
		}
	}

	RunRedisServer := func(port string) {

		redis_server := miniredis.NewMiniRedis()

		go CleanExpireKeys(redis_server)

		var er error

		// 使用redis 0 号数据库
		redis_server.DB(0)

		// 如果配置文件中redis密码不为空，则配置redis密码访问
		if config.RedisPassword != "" {
			redis_server.RequireAuth(config.RedisPassword)

			er = redis_server.StartAddr(":" + port)
		} else {
			er = redis_server.StartAddr(":" + port)
		}
		if er != nil {
			fmt.Println(er)
		}

	}

	infos := strings.Split(config.RedisHost, ":")
	if len(infos) <= 1 {
		msg := "Redis address must be contain port"
		if action == "init" {
			resp = msg
			panic(resp)
		} else {
			log.Fatalln(msg)
		}
	}
	fmt.Println("127.0.0.1:" + infos[1])
	go RunRedisServer(infos[1])

	time.Sleep(100 * time.Second)
}
