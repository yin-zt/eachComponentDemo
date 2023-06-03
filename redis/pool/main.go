package main

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	log "github.com/sirupsen/logrus"
	"github.com/yin-zt/eachComponentDemo/config"
	"time"
)

var (
	resp any
)

func main() {
	pool := &redis.Pool{
		MaxIdle:     16,
		MaxActive:   1024,
		IdleTimeout: 240 * time.Second,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", config.RedisHost,
				redis.DialConnectTimeout(time.Duration(5)*time.Second),
				redis.DialPassword(config.RedisPassword),
				redis.DialDatabase(0),
			)
			if err != nil {
				fmt.Println(err)
				log.Error(err)
			}
			return conn, err

		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("ping")
			if err != nil {
				log.Error(err)
				return err
			}
			return err
		},
	}

	conn := pool.Get()
	defer conn.Close()
	//var cursor uint64

	//conn.Do("hset", "books", "hello", "world")
	//val, err := redis.String(conn.Do("hget", "books", "hello"))
	//val, err := redis.String(conn.Do("scan", cursor, "books"))

	var keys []string
	cursor := 0
	values, err := redis.Values(conn.Do("SCAN", cursor, "MATCH", "book*"))
	if err != nil {
		resp = err
		panic(resp)
	}
	fmt.Println(cursor)
	cursor, _ = redis.Int(values[0], nil)
	scanKeys, _ := redis.Strings(values[1], nil)
	keys = append(keys, scanKeys...)

	if err != nil {
		log.Error(err)
	}
	conn.Close()
	//fmt.Println(val)
}
