package util

import (
	"github.com/samuel/go-zookeeper/zk"
	"time"
	"log"
	"strings"
	"strconv"
)

type ZookeeperConfig struct {
	ZkServer string `yaml:"zookeeper_server"`
	ZkPath string   `yaml:"zookeeper_path"`
}

func SetData()  {

	defer ErrHandle()
	conf := GetConfigCenterInstance()

	host := strings.Split(conf.ConfProperties["zookeeper"]["zookeeper_server"],",")

	conn, _, err := zk.Connect(host, 10*time.Second)
	if nil != err {
		log.Panic("load config error: ", err)
		return
	}

	_, stat, _ := conn.Exists(conf.ConfProperties["zookeeper"]["zookeeper_path"])

	cur := time.Now()
	timestamp := cur.UnixNano()
	timeStr:=strconv.FormatInt(timestamp,10)

	//b_buf := bytes.NewBuffer([]byte{})
	//binary.Write(b_buf, binary.BigEndian, timestamp)

	conn.Set(conf.ConfProperties["zookeeper"]["zookeeper_path"], []byte(timeStr), stat.Version)

	log.Println("设置zookeepers时间戳 "+timeStr)
}
