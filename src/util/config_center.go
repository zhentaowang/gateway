package util

import (
	"sync"
	"code.aliyun.com/wyunshare/wyun-zookeeper/go-client/src/conf_center"
	"os"
)

var m conf_center.AppProperties
var once sync.Once

func GetConfigCenterInstancePro(appName string) conf_center.AppProperties{
	envName := GetEnvName("local_env")
	if(len(envName) > 0){
		appName = appName + "-" + envName
	}
	once.Do(func() {
		m = conf_center.New(appName)
		m.Init()
	})
	return m
}

func GetConfigCenterInstance() conf_center.AppProperties{
	return m
}

func GetEnvName(env string) string {
	return os.Getenv(env)
}
