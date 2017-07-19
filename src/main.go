// +build linux

package main

import (
    //"gateway/src/admin"
    "gateway/src/model"
    "log"
    "gateway/src/proxy"
    "github.com/samuel/go-zookeeper/zk"
    "time"
    "gateway/src/admgateway/handler"
    "strings"
    "gateway/src/thrift"
    "gateway/src/util"
)



func main() {
    // 读取配置文件
    log.SetFlags(log.Llongfile)
    conf := util.GetConfigCenterInstance()

    log.Println("init gateway success")
    // 获取数据库
    store := model.NewMysqlStore(conf.ConfProperties["jdbc"]["db_host"], conf.ConfProperties["jdbc"]["db_username"], conf.ConfProperties["jdbc"]["db_password"], conf.ConfProperties["jdbc"]["db_name"])

    // 转发服务
    h := proxy.NewHttpProxy(store)
    go DataChange(h)
    go thrift.StartThriftServer()
    go handler.Run()
    h.Start()

    // 管理服务
    //s := admin.NewAdminServer(":8080", "luojing", "111111", store)
    //s.Start()

}

func DataChange(h *proxy.HttpProxy)  {
    log.SetFlags(log.Lshortfile)

    conf := util.GetConfigCenterInstance()

    host := strings.Split(conf.ConfProperties["zookeeper"]["zookeeper_server"],",")

    conn, _, err := zk.Connect(host, 10*time.Second)
    if nil != err {
        log.Print("load config error: ", err)
        return
    }

    for {
        _, _, stat, _ := conn.GetW(conf.ConfProperties["zookeeper"]["zookeeper_path"])
        <-stat
        log.Println("data changed")
        h.Init()
    }
}