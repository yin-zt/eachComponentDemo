package main

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"github.com/yin-zt/eachComponentDemo/config"
	"time"
)

var (
	redisLoger log.Logger
)

type RedisPool struct {
	Pool *redis.Pool
}

// NewRedisPooler 用于生成一个redis连接池
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
		Pool: pool,
	}
}

// RedisHLen 用于求出目标用户处于交接状态的资源数量
func (Redis *RedisPool) RedisHLen(username string) (length int, err error) {
	conn := Redis.Pool.Get()
	defer conn.Close()

	val, err := conn.Do("HLEN", fmt.Sprintf("%s.%s", config.RedisHxPrefix, username))
	if err != nil {
		redisLoger.Errorf("RedisHLen连接redis时报错，报错内容为: %v", err)
		return 0, err
	}
	return cast.ToInt(val), err
}

// RedisHKeys 用于查询用户处于交接状态的资源
func (Redis *RedisPool) RedisHKeys(username string) (keys []string, err error) {
	conn := Redis.Pool.Get()
	defer conn.Close()

	val, err := conn.Do("HKEYS", fmt.Sprintf("%s.%s", config.RedisHxPrefix, username))
	if err != nil {
		redisLoger.Errorf("RedisHKeys 连接redis时报错，报错内容为: %v", err)
		return nil, err
	}
	for _, item := range cast.ToSlice(val) {
		keys = append(keys, cast.ToString(item))
	}
	return keys, nil
}

// RedisHExists 用于查询用户指定资源是否存在
func (Redis *RedisPool) RedisHExists(username string, objId string) bool {
	conn := Redis.Pool.Get()
	defer conn.Close()

	val, err := conn.Do("HEXISTS", fmt.Sprintf("%s.%s", config.RedisHxPrefix, username), objId)
	if err != nil {
		redisLoger.Errorf("RedisHExists 连接redis时报错，报错内容为: %v", err)
		return false
	}
	if cast.ToInt(val) == 1 {
		return true
	} else {
		return false
	}
}

// RedisHGet 用于获取用户指定资源id交接信息
func (Redis *RedisPool) RedisHGet(username string, objId string) string {
	conn := Redis.Pool.Get()
	defer conn.Close()

	val, err := conn.Do("HGET", fmt.Sprintf("%s.%s", config.RedisHxPrefix, username), objId)
	if err != nil {
		redisLoger.Errorf("RedisHExists 连接redis时报错，报错内容为: %v", err)
		return ""
	}
	return cast.ToString(val)
}

// RedisHSet 用于记录处于中间状态的交接资源
func (Redis *RedisPool) RedisHSet(oldUser, objId, newUser string) error {
	conn := Redis.Pool.Get()
	defer conn.Close()
	if objId == "" || newUser == "" {
		redisLoger.Errorf("RedisHSet方法中objId和newUser的值不能为空,当前有一个值为空：[%v,%v]", objId, newUser)
		return fmt.Errorf("传参有空值")
	}
	result, err := conn.Do("HSET", fmt.Sprintf("%s.%s", config.RedisHxPrefix, oldUser), objId, newUser)
	if err != nil {
		redisLoger.Errorf("RedisHSet 执行时时报错，报错内容为: %v", err)
		return err
	}
	if cast.ToInt(result) == 0 {
		redisLoger.Infof("RedisHSet 设置用户【%v】交接资源时，此id 【%v】已经处于交接中，即重复设置了,新用户为: 【%v】", oldUser, objId, newUser)
	}
	return err
}

// RedisHGetAll 用于获取用户所有处于交接状态的资源
func (Redis *RedisPool) RedisHGetAll(oldUser string, storeMap map[string]string) error {
	var (
		resp1 any
	)
	defer func() {
		if err := recover(); err != resp1 {
			fmt.Println("捕获到了panic 产生的异常： ", err)
			fmt.Println("捕获到panic的异常了，recover并没有恢复回来")
			redisLoger.Errorf("RedisHGetAll 捕获到panic异常，recover并没有恢复回来了，【err】为：%s", err)
		}
	}()
	conn := Redis.Pool.Get()
	defer conn.Close()

	result, err := conn.Do("HGETALL", fmt.Sprintf("%s.%s", config.RedisHxPrefix, oldUser))
	if err != nil {
		redisLoger.Errorf("RedisHGetAll 执行时时报错，报错内容为: %v", err)
		return err
	}
	resTurn := cast.ToSlice(result)
	for i := 0; i < len(resTurn); i = i + 2 {
		tempKey := cast.ToString(resTurn[i])
		tempVal := cast.ToString(resTurn[i+1])
		storeMap[tempKey] = tempVal
	}
	return err
}

// RedisHDel 用于删除redis中处于中间状态的交接资源
func (Redis *RedisPool) RedisHDel(oldUser, objId string) error {
	conn := Redis.Pool.Get()
	defer conn.Close()

	_, err := conn.Do("HDEL", fmt.Sprintf("%s.%s", config.RedisHxPrefix, oldUser), objId)
	if err != nil {
		redisLoger.Errorf("RedisHDel 执行删除key操作异常, 用户为: %v; 资源ID为: %v", oldUser, objId)
	}
	return err
}

// RedisHDelUserAll 用于删除指定用户在redis上所有处于中间状态的资源
func (Redis *RedisPool) RedisHDelUserAll(oldUser string) error {
	conn := Redis.Pool.Get()
	defer conn.Close()

	hasKeys, err := Redis.RedisHKeys(oldUser)
	if err != nil {
		redisLoger.Errorf("RedisHDelUserAll 返回报错，报错内容为：%v", err)
		return err
	}
	for _, oneKey := range hasKeys {
		er := Redis.RedisHDel(oldUser, oneKey)
		if er != nil {
			redisLoger.Errorf("RedisHDelUserAll 中删除用户指定key是异常，key为: %v, 用户为: %v", oneKey, oldUser)
			continue
		}
	}
	return nil
}
