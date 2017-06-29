package main

import (
    //"gateway/src/admin"
    "gateway/src/model"
    "log"
    "gateway/src/proxy"
    "github.com/samuel/go-zookeeper/zk"
    "time"
    "code.aliyun.com/wyunshare/wyun-zookeeper/go-client/src/conf_center"
    "gateway/src/admgateway/handler"
    "strings"
)



func main() {
    // 读取配置文件
    conf := conf_center.New("gateway")
    conf.Init()

    // 获取数据库
    store := model.NewMysqlStore(conf.ConfProperties["jdbc"]["db_host"], conf.ConfProperties["jdbc"]["db_username"], conf.ConfProperties["jdbc"]["db_password"], conf.ConfProperties["jdbc"]["db_name"])

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

    conf := conf_center.New("gateway")
    conf.Init()

    host := strings.Split(conf.ConfProperties["zookeeper"]["zookeeper_server"],",")

    conn, _, err := zk.Connect(host, 10*time.Second)
    if nil != err {
        log.Panic("load config error: ", err)
        return
    }

    for {
        b, _, stat, _ := conn.GetW(conf.ConfProperties["zookeeper"]["zookeeper_path"])
        <-stat
        log.Println("data changed")
        println(string(b))
        h.InitRouteTable()
    }
}