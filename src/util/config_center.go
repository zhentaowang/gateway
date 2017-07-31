package util

import (
	"sync"
	"code.aliyun.com/wyunshare/wyun-zookeeper/go-client/src/conf_center"
	"os"
	"fmt"
)

var m conf_center.AppProperties
var once sync.Once

func GetConfigCenterInstance() conf_center.AppProperties{
	defer func () {
		if err := recover(); err != nil {
			fmt.Println("ERROR!! ,配置中心连接失败 ,",err)
		}
	}()

	once.Do(func() {
		envName := GetEnvName("local_env")
		var appName = "gateway"
		if(len(envName) > 0){
			appName = appName + "-" + envName
		}
		m = conf_center.New(appName)
		m.Init()
	})
	return m
}

func GetEnvName(env string) string {
	return os.Getenv(env)
}
