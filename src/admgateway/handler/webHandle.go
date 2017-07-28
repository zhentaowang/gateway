package handler

import (
	"github.com/valyala/fasthttp"
	"path/filepath"
	"html/template"
	"log"
	"bytes"
	"strconv"
	"encoding/json"
	"os"
	"github.com/PuerkitoBio/goquery"
	"gateway/src/util"
)

type WebData struct {
	Title       string
	ApiData     []Api
	ServiceData []Service
	FilterData  []Filter
	Name        string
}


func GetApiFormData(ctx *fasthttp.RequestCtx) (*Api,int) {

	aa := ctx.Request.String()

	println(aa)

	postValues := ctx.PostArgs()
	FormData := new(Api)

	FormData.ApiId, _ = strconv.Atoi(string(postValues.Peek("Api_Api_id")[:]))
	FormData.Desc = string(postValues.Peek("Api_desc")[:])
	FormData.ServiceProviderName = string(postValues.Peek("Api_Service_Provider_name")[:])
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

func GetUploadApiFile(ctx *fasthttp.RequestCtx) []*Api {
	fileData := make([]*Api, 0)
	temp,_ :=ctx.FormFile("fileUpload")
	f ,_:= temp.Open()
	defer f.Close()

	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		log.Println(err.Error())
	}

	//dhead := doc.Find(".module .action")
	doc.Find(".module .page").Each(func(i int, contentSelection *goquery.Selection) {
		contentSelection.Find("h3").Each(func(i int, contentSelection2 *goquery.Selection) {
			if contentSelection2.Text()=="URL" {
				dataT := new(Api)
				dataT.Uri = contentSelection2.Next().Text()
				dataT.Method = contentSelection2.Next().Next().Next().Text()
				dataT.Status = 0
				if len(dataT.Uri)!=0 && (dataT.Method == "POST"||dataT.Method == "GET"||dataT.Method == "*") {
					fileData = append(fileData,dataT)
				}
			}
		})
	})

	return fileData
}

func GetServiceFormData(ctx *fasthttp.RequestCtx)  *Service {
	postValues := ctx.PostArgs()
	FormData := new(Service)

	FormData.ServiceId, _ = strconv.Atoi(string(postValues.Peek("Service_id")[:]))
	FormData.Desc = string(postValues.Peek("Service_desc")[:])
	FormData.Name = string(postValues.Peek("Service_name")[:])
	FormData.Namespace = string(postValues.Peek("Service_namespace")[:])
	FormData.Port = string(postValues.Peek("Service_port")[:])
	FormData.Protocol = string(postValues.Peek("Service_protocol")[:])

	return FormData
}

func GetFilterFormData(ctx *fasthttp.RequestCtx)  *Filter {
	postValues := ctx.PostArgs()
	FormData := new(Filter)

	FormData.FilterId, _ = strconv.Atoi(string(postValues.Peek("Filter_id")[:]))
	FormData.ApiId, _ = strconv.Atoi(string(postValues.Peek("Filter_Api_id")[:]))
	FormData.Name = string(postValues.Peek("Filter_name")[:])
	FormData.Seq, _ = strconv.Atoi(string(postValues.Peek("Filter_seq")[:]))

	return FormData

}

func GetLoginData(ctx *fasthttp.RequestCtx)  *LoginData{
	postValues := ctx.PostArgs()
	FormData := new(LoginData)

	FormData.name = string(postValues.Peek("txtName")[:])
	FormData.password = string(postValues.Peek("txtPwd")[:])

	return FormData
}


func Render(ctx *fasthttp.RequestCtx, url string,data interface{}) {

	defer util.ErrHandle()
	if(url == "") {
		JsonResult(ctx,data)
		return
	}
	pwd, _ := os.Getwd()

	htmlFile := filepath.Join(pwd, "src","admgateway","view", url)

	cssFile := filepath.Join(pwd, "src","admgateway","view", "index.css")

	t, err := template.ParseFiles(htmlFile,cssFile)
	if err != nil {
		log.Fatal(err)
	}

	wr := bytes.NewBufferString("")
	t.ExecuteTemplate(wr,"content",data)
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

func RedirectIndex(ctx *fasthttp.RequestCtx, data interface{})  {

	//url := filepath.Join( "src","admgateway","view", "delete.html")
	//
	//Render(ctx, url, data)

}