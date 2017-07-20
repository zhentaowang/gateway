package utils

import (
	"sync"
	"code.aliyun.com/wyunshare/wyun-zookeeper/go-client/src/conf_center"
	"os"
)

var m conf_center.AppProperties
var once sync.Once

func GetConfigCenterInstance(appName string) conf_center.AppProperties{
	envName := GetEnvName("local_env")
	if(len(envName) > 0){
		appName = appName + "-" + GetEnvName("local_env")
	}
	once.Do(func() {
		m = conf_center.New(appName)
		m.Init()
	})
	return m
}

func GetEnvName(env string) string {
	return os.Getenv(env)
}