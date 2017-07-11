package util

import (
	"sync"
	"code.aliyun.com/wyunshare/wyun-zookeeper/go-client/src/conf_center"
)


var m conf_center.AppProperties
var once sync.Once

func GetConfigCenterInstance() conf_center.AppProperties{
	once.Do(func() {
		m = conf_center.New("gateway")
		m.Init()
	})
	return m
}
