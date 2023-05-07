package main

import (
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"github.com/yin-zt/eachComponentDemo/config"
	"time"
)

var (
	resp any
)

func main() {
	conn, _, err := zk.Connect([]string{fmt.Sprintf("%s:%s", config.ZkHost, config.ZkPort)}, 5*time.Second)
	if err != nil {
		resp = err
		panic(resp)
	}

	// 1.验证根是否存在
	if hasRoot, _, _ := conn.Exists("/root1"); !hasRoot {
		// 2.新增根
		_, err = conn.Create("/root1", []byte("root_content"), 0, zk.WorldACL(zk.PermAll))
		if err != nil {
			fmt.Println("failed add root, info: ", err.Error())
		}
		fmt.Println("add /root1 success")
	}

	// 3.查询根
	data, stat, err := conn.Get("/root1")
	if err != nil {
		fmt.Println("failed get root, info: ", err.Error())
	}
	fmt.Println("text: ", string(data), stat.Version)

	// 4.修改根
	if _, err = conn.Set("/root1", []byte("update text"), stat.Version); err != nil {
		fmt.Println("failed update root")
	}

	// 5.设置子节点(必须要有根/父节点)
	if _, err = conn.Create("/root1/subnode", []byte("node_text"), 0, zk.WorldACL(zk.PermAll)); err != nil {
		fmt.Println("failed add subnode, info: ", err.Error())
	}
	// 6.获取子节点列表
	childNodes, _, err := conn.Children("/root1")
	if err != nil {
		fmt.Println("failed get node list, info: ", err.Error())
	} else {
		fmt.Println("node list: ", childNodes)
	}

	// 6.删除根(必须先查后删, 删完子才能删父节点)
	_, stat, _ = conn.Get("/root1/subnode")
	if err := conn.Delete("/root1/subnode", stat.Version); err != nil {
		fmt.Println("falied delete node, info: ", stat.Version, err.Error())
	}
	_, stat, _ = conn.Get("/root1")
	if err := conn.Delete("/root1", stat.Version); err != nil {
		fmt.Println("falied delete root, info: ", stat.Version, err.Error())
	}
}

// WatchHostsByPath 监控节点变化
func WatchHostsByPath(path string, conn *zk.Conn) (chan []string, chan error) {
	snapshots := make(chan []string) // 变动后的挂载目标列表
	errors := make(chan error)       // 变动错误信息
	go func() {
		for {
			snapshot, _, events, err := conn.ChildrenW(path)
			if err != nil {
				errors <- err
			}
			snapshots <- snapshot
			// 阻塞直到出现事件消息(可以不用select)
			select {
			case evt := <-events:
				if evt.Err != nil {
					errors <- evt.Err
				}
				fmt.Printf("ChildrenW Event Path:%v, Type:%v\n", evt.Path, evt.Type)
			}
		}
	}()
	return snapshots, errors
}

// WatchDataByPath 监控节点内容变化
func WatchDataByPath(nodePath string, conn *zk.Conn) (chan []byte, chan error) {
	//conn := z.conn
	snapshots := make(chan []byte)
	errors := make(chan error)
	go func() {
		for {
			data, _, events, err := conn.GetW(nodePath)
			if err != nil {
				errors <- err
			}
			snapshots <- data
			select {
			case evt := <-events:
				if evt.Err != nil {
					errors <- evt.Err
					return
				}
				fmt.Printf("GetW Event Path:%v, Type:%v\n", evt.Path, evt.Type)
			}
		}
	}()
	return snapshots, errors
}

// RegistHostOnPath 将主机挂载在路径上(节点上)
func RegistHostOnPath(nodePath, host string, conn *zk.Conn) (err error) {
	// 1. 若路径不存在则新建
	ex, _, err := conn.Exists(nodePath)
	if err != nil {
		return
	}
	if !ex {
		_, err = conn.Create(nodePath, nil, 0, zk.WorldACL(zk.PermAll))
		if err != nil {
			return
		}
	}
	// 2. 将主机进行挂载(路径为永久、主机为临时); 临时会检测主机的存活并清理
	subNodePath := fmt.Sprintf("%s/%s", nodePath, host)
	// 2.1 主机是否已挂载, 没挂则挂
	if ex, _, _ := conn.Exists(subNodePath); !ex {
		_, err = conn.Create(subNodePath, nil, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	}
	return
}

// UpdatePathData 有则更新、无则新建节点并配置内容
func UpdatePathData(nodePath string, config []byte, version int32, conn *zk.Conn) (err error) {
	ex, _, _ := conn.Exists(nodePath)
	if !ex {
		conn.Create(nodePath, config, 0, zk.WorldACL(zk.PermAll))
		return nil
	}
	// 需要版本才能更新
	_, stat, err := conn.Get(nodePath)
	if err != nil {
		return
	}
	_, err = conn.Set(nodePath, config, stat.Version)
	return
}
