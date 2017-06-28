package handler

import (
	"github.com/valyala/fasthttp"
	"fmt"
	"github.com/buaazp/fasthttprouter"
	"io/ioutil"
	"log"
	"gopkg.in/yaml.v2"
	"path/filepath"
)

type Webconf struct {
	Host string `yaml:"adm_web_host"`
}

func Run() {

	file := filepath.Join( "src","admgateway", "conf.yml")
	configByte, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	wc := new(Webconf)
	err = yaml.Unmarshal(configByte, wc)
	if nil != err {
		log.Panic("load config error: ", err)
		return
	}


	router := fasthttprouter.New()
	InitRoute(router)

	if err := fasthttp.ListenAndServe(wc.Host, router.Handler); err != nil {
		fmt.Println("start fasthttp fail:", err.Error())
	} else {
		fmt.Println("start admweb success")
	}

}
