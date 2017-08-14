package handler

import (
	"github.com/valyala/fasthttp"
	"fmt"
	"github.com/buaazp/fasthttprouter"
	"gateway/src/util"
)

type Webconf struct {
	Host string `yaml:"adm_web_host"`
}

func Run() {

	conf := util.GetConfigCenterInstance()

	testService := new(Service)
	testService.Name = "0.0.0.0"
	testService.Port = "1323"
	testService.Protocol = "http"
	QueryOneService(testService)
	if testService.ServiceId == 0 {
	    MInsertService(testService)
	}
	QueryOneService(testService)

	testApi := new(Api)
	testApi.Name = "testwebmanager"
	testApi.Uri = "/api.html"
	testApi.Method = "GET"
	testApi.NeedLogin = 0
	testApi.Status = 1
	QueryOneApi(testApi)
	testApi.ServiceId = testService.ServiceId
	if testApi.ApiId == 0 {
		MInsertApi(testApi,0)
	}

	router := fasthttprouter.New()
	InitRoute(router)

	if err := fasthttp.ListenAndServe(conf.ConfProperties["adm_web"]["adm_web_host"], router.Handler); err != nil {
		fmt.Println("start fasthttp fail:", err.Error())
	} else {
		fmt.Println("start admweb success")
	}

}
