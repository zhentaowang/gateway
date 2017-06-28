package util

import (
	"github.com/samuel/go-zookeeper/zk"
	"time"
	"encoding/binary"
	"log"
	"bytes"
	"io/ioutil"
	"strings"
	"gopkg.in/yaml.v2"
)

type ZookeeperConfig struct {
	ZkServer string `yaml:"zookeeper_server"`
	ZkPath string   `yaml:"zookeeper_path"`
}

func SetData()  {

	configByte, err := ioutil.ReadFile("conf.yml")
	if err != nil {
		log.Fatal(err)
	}

	zkConf := new(ZookeeperConfig)
	err = yaml.Unmarshal(configByte, &zkConf)
	if nil != err {
		log.Panic("load config error: ", err)
		return
	}

	host := strings.Split(zkConf.ZkServer,",")

	conn, _, err := zk.Connect(host, 10*time.Second)
	if nil != err {
		log.Panic("load config error: ", err)
		return
	}

	_, stat, _ := conn.Exists(zkConf.ZkPath)

	cur := time.Now()
	timestamp := cur.UnixNano()

	b_buf := bytes.NewBuffer([]byte{})
	binary.Write(b_buf, binary.BigEndian, timestamp)

	conn.Set(zkConf.ZkPath, b_buf.Bytes(), stat.Version)
}
