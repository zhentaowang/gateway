package handler

import (
	"github.com/valyala/fasthttp"
)



func indexHandler(ctx *fasthttp.RequestCtx) {

	url := "api.html"

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

	ctx.Response.Header.Set("Location", "/")

	Render(ctx, url, data)
}

func ToService(ctx *fasthttp.RequestCtx)  {

	url := "filter.html"

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

func ToFilter(ctx *fasthttp.RequestCtx)  {

	url := "service.html"

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