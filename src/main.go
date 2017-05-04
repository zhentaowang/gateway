package main

import (
    //"gateway/admin"
    "gateway/src/proxy"
    "gateway/src/model"
    "github.com/labstack/gommon/log"
    "github.com/go-yaml/yaml"
    "io/ioutil"
    "gateway/src/config"
)

func main() {
    // 读取配置文件
    configByte, err := ioutil.ReadFile("config.yml")
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
    h.Start()

    // 管理服务
    //s := admin.NewAdminServer(":8080", "luojing", "111111", store)
    //s.Start()

}