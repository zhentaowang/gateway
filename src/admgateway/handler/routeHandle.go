package handler

import (
	"github.com/valyala/fasthttp"
	"path/filepath"
)



func indexHandler(ctx *fasthttp.RequestCtx) {

	url := filepath.Join( "src","admgateway","view", "index.html")

	data := struct {
		Title string
		ApiData []Api
		ServiceData []Service
		FilterData []Filter
	}{
		Title: "Gateway Manager",
		ApiData: MQueryApi(new(Api)),
		ServiceData: MQueryService(new(Service)),
		FilterData: MQueryFilter(new(Filter)),
	}


	Render(ctx, url, data)
}

func deleteHandler(ctx *fasthttp.RequestCtx)  {

	url := filepath.Join( "src","admgateway","view", "delete.html")

	data := struct {
		Title string
		ApiData []Api
		ServiceData []Service
		FilterData []Filter
	}{
		Title: "Gateway Manager",
		ApiData: MQueryApi(new(Api)),
		ServiceData: MQueryService(new(Service)),
		FilterData: MQueryFilter(new(Filter)),
	}

	Render(ctx, url, data)

}


func AddApi(ctx *fasthttp.RequestCtx) {

	form_data, filter_seq := GetApiFormData(ctx)

	MInsertApi(form_data, filter_seq)

}


func AddService(ctx *fasthttp.RequestCtx)  {

	form_data := GetServiceFormData(ctx)

	MInsertService(form_data)

}


func AddFilter(ctx *fasthttp.RequestCtx)  {

	form_data := GetFilterFormData(ctx)

	MInsertFilter(form_data)

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

	var filter_data = Filter{0, form_data.ApiId, "", 0}

	MDeleteFilter(&filter_data)
}


func DeleteService(ctx *fasthttp.RequestCtx)  {

	form_data := GetServiceFormData(ctx)


	MDeleteService(form_data)

}


func DeleteFilter(ctx *fasthttp.RequestCtx)  {

	form_data := GetFilterFormData(ctx)

	MDeleteFilter(form_data)

}