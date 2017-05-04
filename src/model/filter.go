package model

import (
    "github.com/valyala/fasthttp"
    "net/http"
    "gateway/src/filter"
)

func (a *API) DoPreFilters(c filter.Context) (filterName string, statusCode int, err error) {
    for item := a.filters.Front(); item != nil; item = item.Next() {
        f, _ := item.Value.(filter.Filter)
        filterName = f.Name()

        statusCode, err = f.Pre(c)
        if nil != err {
            return filterName, statusCode, err
        }
    }

    return "", http.StatusOK, nil
}

func (a *API) DoPostFilters(c filter.Context) (filterName string, statusCode int, err error) {
    for item := a.filters.Back(); item != nil; item = item.Prev() {
        f, _ := item.Value.(filter.Filter)

        statusCode, err = f.Post(c)
        if nil != err {
            return filterName, statusCode, err
        }
    }

    return "", http.StatusOK, nil
}

func (a *API) DoPostErrFilters(c filter.Context) {
    for item := a.filters.Back(); item != nil; item = item.Prev() {
        f, _ := item.Value.(filter.Filter)

        f.PostErr(c)
    }
}

type proxyContext struct {
    startAt   int64
    endAt     int64
    result    *RouteResult
    outerReq  *fasthttp.Request
    originCtx *fasthttp.RequestCtx
    rt        *RouteTable
}

func NewContext(rt *RouteTable, originCtx *fasthttp.RequestCtx, outerReq *fasthttp.Request, result *RouteResult) filter.Context {
    return &proxyContext{
        result:    result,
        originCtx: originCtx,
        outerReq:  outerReq,
        rt:        rt,
    }
}

func (c *proxyContext) GetStartAt() int64 {
    return c.startAt
}

func (c *proxyContext) SetStartAt(startAt int64) {
    c.startAt = startAt
}

func (c *proxyContext) GetEndAt() int64 {
    return c.endAt
}

func (c *proxyContext) SetEndAt(endAt int64) {
    c.endAt = endAt
}

func (c *proxyContext) GetProxyOuterRequest() *fasthttp.Request {
    return c.outerReq
}

func (c *proxyContext) GetProxyResponse() *fasthttp.Response {
    return c.result.Res
}

func (c *proxyContext) GetOriginRequestCtx() *fasthttp.RequestCtx {
    return c.originCtx
}
