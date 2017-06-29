package handler

import (
	"github.com/valyala/fasthttp"
	"fmt"
	"github.com/buaazp/fasthttprouter"
	"code.aliyun.com/wyunshare/wyun-zookeeper/go-client/src/conf_center"
)

type Webconf struct {
	Host string `yaml:"adm_web_host"`
}

func Run() {

	conf := conf_center.New("gateway")
	conf.Init()

	router := fasthttprouter.New()
	InitRoute(router)

	if err := fasthttp.ListenAndServe(conf.ConfProperties["adm_web"]["adm_web_host"], router.Handler); err != nil {
		fmt.Println("start fasthttp fail:", err.Error())
	} else {
		fmt.Println("start admweb success")
	}

}
