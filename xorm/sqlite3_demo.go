package main

import (
	"fmt"
	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/prometheus/common/log"
	"os"
	"time"
)

var (
	dbType = "sqlite3"
	url    = "./test.db"
	action = "init"
)

func main() {
	var (
		err    error
		engine *xorm.Engine
	)

	// 根据cfg.json配置文件中db的类型生成一个相应数据库的驱动; 默认为使用sqlite数据库驱动
	engine, err = xorm.NewEngine(dbType, url)

	if err == nil {
		engine.SetConnMaxLifetime(time.Duration(60) * time.Second)
		engine.SetMaxIdleConns(0)
		fmt.Println(dbType)
	} else {
		fmt.Println(err)
		fmt.Println("Init engine Error")
		log.Error("Init engine Error", err)
		if action == "init" {
			os.Exit(1)
		}
	}

	if err := engine.Sync2(new(TChUser)); err != nil {
		log.Fatalln(err)
	}
}

type TChUser struct {
	Fid         int       `xorm:"not null pk autoincr INT(11)"`
	Femail      string    `xorm:"not null default '' VARCHAR(256)"`
	Fuser       string    `xorm:"not null default '' unique VARCHAR(64)"`
	Fpwd        string    `xorm:"not null default '' VARCHAR(256)"`
	Fip         string    `xorm:"not null default '' VARCHAR(32)"`
	Flogincount int       `xorm:"not null default 0 INT(11)"`
	Ffailcount  int       `xorm:"not null default 0 INT(11)"`
	Flasttime   time.Time `xorm:"updated"`
	Fstatus     int       `xorm:"not null default 0 INT(11)"`
	FmodifyTime time.Time `xorm:"updated index DATETIME"`
	Fversion    int       `xorm:"not null default 0 INT(11)"`
}
