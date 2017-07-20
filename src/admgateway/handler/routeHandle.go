package handler

import (
	"github.com/valyala/fasthttp"
	"gateway/src/util"
	"log"
	"strings"
)



func indexHandler(ctx *fasthttp.RequestCtx) {

	url := "api.html"
	var LoginInfo LoginData = LoginData{}

	data := new(WebData)

	GetCookie(ctx,&LoginInfo)

	ToOauth(&LoginInfo,data)

	Render(ctx, url, data)
}

func ToFilter(ctx *fasthttp.RequestCtx)  {

	url := "filter.html"
	var LoginInfo LoginData = LoginData{}

	data := new(WebData)

	GetCookie(ctx,&LoginInfo)

	ToOauth(&LoginInfo,data)

	ctx.Response.Header.Set("Location", "/")

	Render(ctx, url, data)
}

func ToService(ctx *fasthttp.RequestCtx)  {

	url := "service.html"
	var LoginInfo LoginData = LoginData{}

	data := new(WebData)

	GetCookie(ctx,&LoginInfo)

	ToOauth(&LoginInfo,data)

	ctx.Response.Header.Set("Location", "/")

	Render(ctx, url,data)

}

func Login(ctx *fasthttp.RequestCtx)  {
	log.Printf("logining...")

	indexHandler(ctx)

	cookie := fasthttp.AcquireCookie()
	LoginInfo := GetLoginData(ctx)

	cookie.SetKey("login")
	log.Printf("name=%s,password=%s",LoginInfo.name,LoginInfo.password)
	cookie.SetValue(LoginInfo.name+"&"+LoginInfo.password)
	ctx.Response.Header.SetCookie(cookie)
}

func ToOauth(LoginInfo *LoginData,data *WebData)  {
	allowed, _ := util.GetPermission(LoginInfo.name,LoginInfo.password,"gateway_list")

	if allowed==true {
		data.Title = "Gateway Manager"
		data.ApiData = MQueryApi(new(Api))
		data.ServiceData = MQueryService(new(Service))
		data.FilterData = MQueryFilter(new(Filter))
		data.Name = LoginInfo.name
	} else {
		data.Name = "登陆"
	}
}

func GetCookie(ctx *fasthttp.RequestCtx,LoginInfo *LoginData)  {
	clientCookie := ctx.Request.Header.Cookie("login")
	//log.Printf(string(clientCookie))
	if clientCookie!=nil&&len(clientCookie)!=0 {
		info := strings.Split(string(clientCookie),"&")
		LoginInfo.name = info[0]
		LoginInfo.password = info[1]
	}
}


func AddApi(ctx *fasthttp.RequestCtx) {

	form_data, filter_seq := GetApiFormData(ctx)

	MInsertApi(form_data, filter_seq)

	indexHandler(ctx)
}


func AddService(ctx *fasthttp.RequestCtx)  {

	form_data := GetServiceFormData(ctx)

	MInsertService(form_data)

	indexHandler(ctx)

}


func AddFilter(ctx *fasthttp.RequestCtx)  {

	form_data := GetFilterFormData(ctx)

	MInsertFilter(form_data)

	indexHandler(ctx)

}


func QueryApi(ctx *fasthttp.RequestCtx) {

	form_data , _:= GetApiFormData(ctx)

	result := MQueryApi(form_data)

	Render(ctx, "", result)

}


func QueryService(ctx *fasthttp.RequestCtx)  {

	form_data := GetServiceFormData(ctx)

	MQueryService(form_data)

}


func QueryFilter(ctx *fasthttp.RequestCtx)  {

	form_data := GetFilterFormData(ctx)

	MQueryFilter(form_data)

}


func ModifyApi(ctx *fasthttp.RequestCtx) {

	form_data, _ := GetApiFormData(ctx)

	MModifyApi(form_data)
}


func ModifyService(ctx *fasthttp.RequestCtx)  {

	form_data := GetServiceFormData(ctx)

	MModifyService(form_data)

}


func ModifyFilter(ctx *fasthttp.RequestCtx)  {

	form_data := GetFilterFormData(ctx)

	MModifyFilter(form_data)

}


func DeleteApi(ctx *fasthttp.RequestCtx)  {

	form_data , _:= GetApiFormData(ctx)
	MDeleteApi(form_data)

	indexHandler(ctx)

}


func DeleteService(ctx *fasthttp.RequestCtx)  {

	form_data := GetServiceFormData(ctx)

	MDeleteService(form_data)

	indexHandler(ctx)

}


func DeleteFilter(ctx *fasthttp.RequestCtx)  {

	form_data := GetFilterFormData(ctx)

	MDeleteFilter(form_data)

	indexHandler(ctx)

}