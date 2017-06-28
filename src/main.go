package main

import (
    //"gateway/src/admin"
    "gateway/src/model"
    "gopkg.in/yaml.v2"
    "io/ioutil"
    "gateway/src/config"
    "log"
    "gateway/src/proxy"
    "github.com/samuel/go-zookeeper/zk"
    "time"
    "strings"
    "gateway/src/admgateway/handler"
)



func main() {
    // 读取配置文件
    configByte, err := ioutil.ReadFile("conf.yml")
    if err != nil {
        log.Fatal(err)
    }

    config.TConfig = config.T{}
    err = yaml.Unmarshal(configByte, &config.TConfig)
    if nil != err {
        log.Panic("load config error: ", err)
        return
    }

    // 获取数据库
    store := model.NewMysqlStore(config.TConfig.DBHost, config.TConfig.DBUsername, config.TConfig.DBPassword, config.TConfig.DBName)

    // 转发服务
    h := proxy.NewHttpProxy(store)
    go DataChange(h)
    go handler.Run()
    h.Start()

    // 管理服务
    //s := admin.NewAdminServer(":8080", "luojing", "111111", store)
    //s.Start()

}

func DataChange(h *proxy.HttpProxy)  {

    configByte, err := ioutil.ReadFile("conf.yml")
    if err != nil {
        log.Fatal(err)
    }

    zkConf := new(config.ZookeeperConfig)
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

    for {
        b, _, stat, _ := conn.GetW(zkConf.ZkPath)
        <-stat
        log.Println("data changed")
        println(string(b))
        h.InitRouteTable()
    }
}