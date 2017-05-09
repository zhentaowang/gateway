package model

import (
    "github.com/valyala/fasthttp"
    "gateway/src/filter"
)



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
