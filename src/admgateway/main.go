package main

import (
	"fmt"

	"github.com/valyala/fasthttp"
	"github.com/buaazp/fasthttprouter"

	_ "github.com/go-sql-driver/mysql"
	"admgateway/handdler"
)






func main() {

	router := fasthttprouter.New()
	handdler.InitRoute(router)

	if err := fasthttp.ListenAndServe("0.0.0.0:1323", router.Handler); err != nil {
		fmt.Println("start fasthttp fail:", err.Error())
	}

}
