package handler

import (
	"github.com/valyala/fasthttp"
	"os"
	"path/filepath"
	"html/template"
	"log"
	"bytes"
	"strconv"
	"encoding/json"
)




func GetApiFormData(ctx *fasthttp.RequestCtx) (*Api,int) {

	postValues := ctx.PostArgs()
	FormData := new(Api)

	FormData.ApiId, _ = strconv.Atoi(string(postValues.Peek("Api_Api_id")[:]))
	FormData.Desc = string(postValues.Peek("Api_desc")[:])
	FormData.DisplayName = string(postValues.Peek("Api_display_name")[:])
	FormData.Filters = string(postValues.Peek("Api_filters")[:])
	FormData.Method = string(postValues.Peek("Api_method")[:])
	FormData.Mock = string(postValues.Peek("Api_mock")[:])
	FormData.Name = string(postValues.Peek("Api_name")[:])
	FormData.NeedLogin, _ = strconv.Atoi(string(postValues.Peek("Api_need_login")[:]))
	FormData.ServiceId, _ = strconv.Atoi(string(postValues.Peek("Api_Service_id")[:]))
	FormData.Status, _ = strconv.Atoi(string(postValues.Peek("Api_status")[:]))
	FormData.Uri = string(postValues.Peek("Api_uri")[:])
	FilterSeq, _ := strconv.Atoi(string(postValues.Peek("Filter_seq")[:]))



	return FormData, FilterSeq

}

func GetServiceFormData(ctx *fasthttp.RequestCtx)  *Service {
	postValues := ctx.PostArgs()
	FormData := new(Service)

	FormData.Desc = string(postValues.Peek("Service_desc")[:])
	FormData.Name = string(postValues.Peek("Service_name")[:])
	FormData.Namespace = string(postValues.Peek("Service_namespace")[:])
	FormData.Port = string(postValues.Peek("Service_port")[:])
	FormData.Protocol = string(postValues.Peek("Service_protocol")[:])
	FormData.ServiceId, _ = strconv.Atoi(string(postValues.Peek("Service_Service_id")[:]))

	return FormData
}

func GetFilterFormData(ctx *fasthttp.RequestCtx)  *Filter {
	postValues := ctx.PostArgs()
	FormData := new(Filter)

	FormData.FilterId, _ = strconv.Atoi(string(postValues.Peek("Filter_Filter_id")[:]))
	FormData.ApiId, _ = strconv.Atoi(string(postValues.Peek("Filter_Api_id")[:]))
	FormData.Name = string(postValues.Peek("Filter_name")[:])
	FormData.Seq, _ = strconv.Atoi(string(postValues.Peek("Filter_seq")[:]))

	return FormData

}


func Render(ctx *fasthttp.RequestCtx, url string, data interface{}) {

	if(url == "") {
		JsonResult(ctx,data)
		return
	}
	pwd, _ := os.Getwd()
	htmlFile := filepath.Join(pwd, url)

	t, err := template.ParseFiles(htmlFile)
	if err != nil {
		log.Fatal(err)
	}

	wr := bytes.NewBufferString("")
	err = t.Execute(wr, data)
	if err != nil {
		log.Fatal(err)
	}

	ctx.SetContentType("text/html")
	ctx.WriteString(wr.String())
}


func JsonResult(ctx *fasthttp.RequestCtx, data interface{}) {

	b, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	ctx.WriteString(string(b))
}
