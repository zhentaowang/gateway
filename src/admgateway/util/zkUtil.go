package util

import (
	"github.com/samuel/go-zookeeper/zk"
	"time"
	"encoding/binary"
	"log"
	"bytes"
	"strings"
	"code.aliyun.com/wyunshare/wyun-zookeeper/go-client/src/conf_center"
)

type ZookeeperConfig struct {
	ZkServer string `yaml:"zookeeper_server"`
	ZkPath string   `yaml:"zookeeper_path"`
}

func SetData()  {

	conf := conf_center.New("gateway")
	conf.Init()

	host := strings.Split(conf.ConfProperties["zookeeper"]["zookeeper_server"],",")

	conn, _, err := zk.Connect(host, 10*time.Second)
	if nil != err {
		log.Panic("load config error: ", err)
		return
	}

	_, stat, _ := conn.Exists(conf.ConfProperties["zookeeper"]["zookeeper_path"])

	cur := time.Now()
	timestamp := cur.UnixNano()

	b_buf := bytes.NewBuffer([]byte{})
	binary.Write(b_buf, binary.BigEndian, timestamp)

	conn.Set(conf.ConfProperties["zookeeper"]["zookeeper_path"], b_buf.Bytes(), stat.Version)
}
