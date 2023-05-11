package main

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"github.com/yin-zt/eachComponentDemo/config"
	"time"
)

type RedisPool struct {
	pool *redis.Pool
}

func NewRedisPooler() *RedisPool {
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
	return &RedisPool{
		pool: pool,
	}
}

func main() {

	myPool := NewRedisPooler()
	conn := myPool.pool.Get()
	defer conn.Close()
	//var cursor uint64

	conn.Do("hset", "books1", "hello", "world")
	conn.Do("hset", "books2", "hello", "world")
	conn.Do("hset", "books3", "hello", "world")
	conn.Do("hset", "books4", "hello", "world")
	//val, err := redis.String(conn.Do("hget", "books", "hello"))
	//val, err := redis.String(conn.Do("scan", cursor, "books"))
	//fmt.Println(val)
	res, err := myPool.RedisKeys("book*")
	fmt.Println(res)
	fmt.Println(err)
}

func (Redis *RedisPool) RedisKeys(key string) (keys []interface{}, err error) {
	cursor := "0"

	conn := Redis.pool.Get()
	defer conn.Close()
	for {
		res, err := conn.Do("SCAN", cursor, "match", key, "count", 100)
		if err != nil {
			break
		}
		fmt.Println(res)
		dataList := res.([]interface{})[1].([]interface{})
		if len(dataList) > 0 {
			// 获取相关的key
			keysTmp := make([]interface{}, len(dataList))
			i := 0
			for _, item := range dataList {
				keysTmp[i] = cast.ToString(item)
				i++
			}
			keys = append(keys, keysTmp...)
		}
		// 获取下个游标
		cursor = cast.ToString(res.([]interface{})[0])
		if cursor == "0" {
			break
		}
	}
	return
}
