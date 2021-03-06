package handler

import (
	"github.com/valyala/fasthttp"
	"gateway/src/util"
	"log"
	"strings"
	"strconv"
)

var fsHandle  = fasthttp.FSHandler("src/admgateway/view",1)

func indexHandler(ctx *fasthttp.RequestCtx) {

	url := "api.html"

	data := new(WebData)

    cookiesData := new(LoginData)
    GetCookie(ctx,cookiesData)

    data.Title = "Gateway Manager"

    if len(cookiesData.name)!=0 {
        data.ApiData = MQueryApi(new(Api))
        data.ServiceData = MQueryService(new(Service))
        data.FilterData = MQueryFilter(new(Filter))
        data.Name = cookiesData.name
    } else {
        data.Name = "登陆"
    }

	Render(ctx, url, data)
}

func ToFilter(ctx *fasthttp.RequestCtx)  {

    url := "filter.html"

    data := new(WebData)

    cookiesData := new(LoginData)
    GetCookie(ctx,cookiesData)

    data.Title = "Gateway Manager"

    if len(cookiesData.name)!=0 {
        data.ApiData = MQueryApi(new(Api))
        data.ServiceData = MQueryService(new(Service))
        data.FilterData = MQueryFilter(new(Filter))
        data.Name = cookiesData.name
    } else {
        data.Name = "登陆"
    }

    ctx.Response.Header.Set("Location", "/")

    Render(ctx, url, data)
}

func ToService(ctx *fasthttp.RequestCtx)  {

	url := "service.html"

	data := new(WebData)

        cookiesData := new(LoginData)
        GetCookie(ctx,cookiesData)

        data.Title = "Gateway Manager"

        if len(cookiesData.name)!=0 {
            data.ApiData = MQueryApi(new(Api))
            data.ServiceData = MQueryService(new(Service))
            data.FilterData = MQueryFilter(new(Filter))
            data.Name = cookiesData.name
        } else {
            data.Name = "登陆"
        }

	ctx.Response.Header.Set("Location", "/")

	Render(ctx, url,data)

}

func Login(ctx *fasthttp.RequestCtx)  {
	log.Printf("logining...")

        cookiesData := new(LoginData)
        data := new(WebData)

	cookie := fasthttp.AcquireCookie()
	LoginInfo := GetLoginData(ctx)

        GetCookie(ctx,cookiesData)

            allowed, _ := util.GetPermission(LoginInfo.name, LoginInfo.password, "gateway_list")
            log.Print("allowed " + LoginInfo.name+"   "+LoginInfo.password)
            log.Print(allowed)

            if allowed==true {
                data.Title = "Gateway Manager"
                data.ApiData = MQueryApi(new(Api))
                data.ServiceData = MQueryService(new(Service))
                data.FilterData = MQueryFilter(new(Filter))
                data.Name = LoginInfo.name

                cookie.SetKey("login")
                cookie.SetValue(LoginInfo.name+"&"+LoginInfo.password)
                ctx.Response.Header.SetCookie(cookie)

                ctx.Response.Header.Set("Location","/api.html")
                ctx.Response.Header.Set("status","true")
                return
            } else {
                data.Name = "登陆"

                ctx.Response.Header.Set("status","false")
                ctx.Response.Header.DelClientCookie("login")
                ctx.Response.Header.Set("Location","/api.html")
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
	lg := util.GetCommonLog()
	var LoginInfo LoginData = LoginData{}
	GetCookie(ctx,&LoginInfo)

	MInsertApi(form_data, filter_seq)

	lg.Println(LoginInfo.name+"  insert api " +form_data.Uri+ " success ")

	indexHandler(ctx)
}

func InsertUploadFile(ctx *fasthttp.RequestCtx)  {

	data := GetUploadApiFile(ctx)
	lg := util.GetCommonLog()
	var LoginInfo LoginData = LoginData{}
	GetCookie(ctx,&LoginInfo)

	MutiInsertApi(data)

	lg.Println(LoginInfo.name+"  insert api from file ")

}


func AddService(ctx *fasthttp.RequestCtx)  {

	form_data := GetServiceFormData(ctx)
	lg := util.GetCommonLog()
	var LoginInfo LoginData = LoginData{}
	GetCookie(ctx,&LoginInfo)

	MInsertService(form_data)
	lg.Println(LoginInfo.name+"  insert service " +form_data.Namespace+"."+ form_data.Name +":"+form_data.Port+ " success ")

	indexHandler(ctx)

}


func AddFilter(ctx *fasthttp.RequestCtx)  {

	form_data := GetFilterFormData(ctx)
	lg := util.GetCommonLog()
	var LoginInfo LoginData = LoginData{}
	GetCookie(ctx,&LoginInfo)

	MInsertFilter(form_data)
	lg.Println(LoginInfo.name + "  insert filter " + form_data.Name + " success ")

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
	lg := util.GetCommonLog()
	var LoginInfo LoginData = LoginData{}
	GetCookie(ctx,&LoginInfo)

	MDeleteApi(form_data)
	lg.Println(LoginInfo.name+"  have deleted api " + form_data.Uri +"  id="+strconv.Itoa(form_data.ApiId)+ "  success ")

	indexHandler(ctx)

}


func DeleteService(ctx *fasthttp.RequestCtx)  {

	form_data := GetServiceFormData(ctx)
	lg := util.GetCommonLog()
	var LoginInfo LoginData = LoginData{}
	GetCookie(ctx,&LoginInfo)

	MDeleteService(form_data)
	lg.Println(LoginInfo.name+"  have deleted service " +form_data.Namespace+"."+ form_data.Name +":"+form_data.Port+"  id="+strconv.Itoa(form_data.ServiceId)+ " success ")

	indexHandler(ctx)

}


func DeleteFilter(ctx *fasthttp.RequestCtx)  {

	form_data := GetFilterFormData(ctx)
	lg := util.GetCommonLog()
	var LoginInfo LoginData = LoginData{}
	GetCookie(ctx,&LoginInfo)

	MDeleteFilter(form_data)
	lg.Println(LoginInfo.name+"  haved deleted filter " + form_data.Name+"  id="+strconv.Itoa(form_data.FilterId) + "  success ")

	indexHandler(ctx)

}