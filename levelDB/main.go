package main

import (
	"encoding/json"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"log"
)

var (
	CONST_LEVELDB_FILE_NAME = "./test.db"
	opts                    = &opt.Options{
		CompactionTableSize: 1024 * 1024 * 20,
		WriteBuffer:         1024 * 1024 * 20,
	}
	//UserRelationField = make(map[string]int)
)

type LevelDb struct {
	DB *leveldb.DB
}

func main() {
	ldb := NewLevelDb()
	ldb.SetKeyInLevelDB("hello", "world")
	result, err := ldb.GetKeyFromDB("hello")
	fmt.Println(result, err)
	ldb.DelKeyFromDB("hello")
	result1, _ := ldb.GetKeyFromDB("hello")
	fmt.Println(result1)
}

func NewLevelDb() *LevelDb {
	db, err := leveldb.OpenFile(CONST_LEVELDB_FILE_NAME, opts)
	if err != nil {
		log.Println("NewLevelDb error, with: %v", err)
	}
	return &LevelDb{db}
}

func (this LevelDb) GetKeyFromDB(key string) (string, error) {
	if data, err := this.DB.Get([]byte(key), nil); err != nil {
		return "", err
	} else {
		return string(data), nil
	}
}

// DelKeyFromDB 从LevelDB中清理指定key
func (this LevelDb) DelKeyFromDB(key string) (string, error) {
	if err := this.DB.Delete([]byte(key), nil); err != nil {
		return "", err
	} else {
		return "", nil
	}
}

func (this LevelDb) SetKeyInLevelDB(key string, value interface{}) error {
	switch trueVal := value.(type) {
	case string:
		err := this.DB.Put([]byte(key), []byte(trueVal), nil)
		return err
	case []string, map[string]interface{}, []map[string]string:
		valStr, err := json.Marshal(trueVal)
		if err != nil {
			log.Fatalln("[]string格式数据序列化失败, 值为: %v", trueVal)
			return err
		}
		err = this.DB.Put([]byte(key), valStr, nil)
		return err
	default:
		log.Fatalln("Unknow the type of value, value is: %v", trueVal)
		return nil
	}
}
