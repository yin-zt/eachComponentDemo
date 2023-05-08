package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-xorm/xorm"
	"github.com/prometheus/common/log"
	"github.com/samuel/go-zookeeper/zk"
	"github.com/sjqzhang/zksdk"
	"os"
	"strings"
	"time"
)

var (
	cli     = &CliServer{}
	engine  *xorm.Engine
	action  = "init"
	redisDb = "0"
	mysqlDb = "test"
	dbPass  = "Hrzd@2020"
	dbUser  = "root"
	dbPort  = "3306"
	dbHost  = "124.71.146.200"
)

type CliServer struct {
	zkdb *zksdk.ZkSdk
	util *Common
}
type Common struct {
}

func main() {

	zkdb := &zksdk.ZkSdk{}
	zkdb.Init([]string{"124.71.146.200:80"}, time.Second*5)

	cli.zkdb = zkdb

	er := zkdb.Start()
	fmt.Println(er)
	fmt.Println("11111111")

	flag := make(chan bool)

	fmt.Println("pppppppppppp")
	for {
		time.Sleep(3 * time.Second)
		if data, _, er := zkdb.GetMoreWX("/root1", flag, cli); er == nil {
			fmt.Println(string(data))
			fmt.Println(5555555555555)
			if true {
				log.Info(string(data))
			}

			//if c, err := cli.util.ParseZkInfo(string(data), "mysql"); err == nil {
			//	if _enginer, er := cli.util.InitEngine(c); er == nil && _enginer != nil {
			//
			//		engine = _enginer
			//
			//	}
			//}

		} else {
			fmt.Println(er)
			fmt.Println(222222)
			log.Error("Connect to Zookeeper Error")
			if action == "init" {
				os.Exit(1)
			}
		}
	}

}

// Proc GetMoreWX 方法监听的节点修改后，会调用Proc方法，data为节点修改后的值
func (this *CliServer) Proc(data []byte, stat *zk.Stat, delFlag bool) {
	fmt.Println("nnnnnnnnnnnnnnnnnn")
	log.Info("update MySQL Config", string(data))

	if c, err := cli.util.ParseZkInfo(string(data), "mysql"); err == nil {
		fmt.Println(c)
		fmt.Println("cccccccccccc")
		if _enginer, er := cli.util.InitEngine(c); er == nil && _enginer != nil {
			if engine != nil {
				engine.Close()
			}
			engine = _enginer
			log.Info("update MySQL Config", string(data))
		}
	}

}

func (this *Common) ParseZkInfo(conf string, dbtype string) (map[string]string, error) {
	data := make(map[string]string)

	db := ""       // 库名
	host := ""     // 数据库节点
	passowrd := "" // 密码
	port := ""     // 端口
	user := ""     // 用户

	if dbtype == "mysql" {

		data["path"] = ""
		data["cmd"] = ""

		db = ""
		host = "124.71.146.200:2181"
		port = ""
		passowrd = ""
		user = ""
		dbtype = "mysql"

	}
	if dbtype == "redis" {

		data["path"] = ""
		data["cmd"] = ""

		db = "0"
		host = "124.71.146.200:2181"
		port = ""
		passowrd = ""
		user = ""
		dbtype = ""

	}

	dbconf := make(map[string]string)

	dbconf["db"] = db
	dbconf["password"] = passowrd
	dbconf["host"] = host
	dbconf["user"] = user
	dbconf["port"] = port
	dbconf["dbtype"] = dbtype

	var conf2 map[string]interface{}

	if err := json.Unmarshal([]byte(conf), &conf2); err != nil {
		log.Error("ParseDBConfig", err)
		return dbconf, err
	}

	if dbtype == "mysql" {

		if v, ok := conf2[dbHost]; ok {
			host = v.(string)
			hosts := strings.Split(host, ",")
			host = hosts[0]

		}
		if v, ok := conf2[dbPort]; ok {

			port = v.(string)
			ports := strings.Split(port, ",")
			port = ports[0]
		}
		if v, ok := conf2[dbPass]; ok {
			passowrd = v.(string)
		}

		if v, ok := conf2[dbUser]; ok {
			user = v.(string)
		}

		if v, ok := conf2[mysqlDb]; ok {
			db = v.(string)
		}

		if v, ok := conf2["mysql"]; ok {
			dbtype = v.(string)
		}
	}
	if dbtype == "redis" {
		if v, ok := conf2[dbHost]; ok {
			host = v.(string)
			hosts := strings.Split(host, ",")
			host = hosts[0]
		}
		if v, ok := conf2[dbPort]; ok {
			port = v.(string)
			ports := strings.Split(port, ",")
			port = ports[0]
		}
		if v, ok := conf2[dbPass]; ok {
			passowrd = v.(string)
		}

		if v, ok := conf2[dbUser]; ok {
			user = v.(string)
		}

		if v, ok := conf2[redisDb]; ok {
			db = v.(string)
		}

		if v, ok := conf2["redis"]; ok {
			dbtype = v.(string)
		}
	}

	dbconf["db"] = db
	dbconf["password"] = passowrd
	dbconf["host"] = host
	dbconf["user"] = user
	dbconf["port"] = port
	dbconf["dbtype"] = dbtype
	return dbconf, nil
}

func (this *Common) InitEngine(c map[string]string) (*xorm.Engine, error) {

	url := "%s:%s@tcp(%s:%s)/%s?charset=utf8"
	url = fmt.Sprintf(url, c["user"], c["password"], c["host"], c["port"], c["db"])
	dbtype := c["dbtype"]

	if true {
		fmt.Println(url)
		log.Info(url)
	}

	_enginer, er := xorm.NewEngine(dbtype, url)

	if er == nil /*&& this.CheckEnginer(_enginer)*/ {
		_enginer.SetConnMaxLifetime(time.Duration(60) * time.Second)
		_enginer.SetMaxIdleConns(0)
		//		_enginer.ShowSQL(true)
		return _enginer, nil
	} else {
		return nil, er
	}

}
