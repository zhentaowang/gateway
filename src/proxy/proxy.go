package proxy

import (
    "github.com/valyala/fasthttp"
    "sync"
    "gateway/src/model"
    "net/http"
    "gateway/src/util"
    "gateway/src/thrift/gen-go/server"
    "encoding/json"
    "gateway/src/config"
    "time"
    "container/list"
    "gateway/src/filter"
    "log"
    "strings"
)

type HttpProxy struct {
    routeTable     *model.RouteTable
    store          model.Store
    fastHTTPClient *util.FastHTTPClient
    filters        *list.List
}

func NewHttpProxy(store model.Store) *HttpProxy {
    h := &HttpProxy{
        fastHTTPClient: util.NewFastHTTPClient(&config.TConfig),
        store: store,
    }

    h.init()

    return h
}

func (h *HttpProxy) init() {
    err := h.initRouteTable()

    if err != nil {
        log.Panic(err, "init route table error")
    }

    filterNames, _ := h.store.GetFilters(-1)
    h.filters = filter.NewFilters(filterNames)
}

func (h *HttpProxy) initRouteTable() error {

    h.routeTable = model.NewRouteTable(h.store)
    h.routeTable.Load()

    return nil
}

func (h *HttpProxy) Start() {
    log.Printf("Proxy exit at %s", fasthttp.ListenAndServe(":8888", h.ReverseProxyHandler))
}

func (h *HttpProxy) ReverseProxyHandler(ctx *fasthttp.RequestCtx) {
    log.Println(string(ctx.Request.RequestURI()))
    log.Println(string(ctx.Request.Body()[:]))
    result := h.routeTable.Select(&ctx.Request)

    if nil == result {
        ctx.SetStatusCode(fasthttp.StatusNotFound)
        return
    }

    h.doProxy(ctx, nil, result)

    if result.Err != nil {
        if result.API.Mock != nil {
            result.API.RenderMock(ctx)
            result.Release()
            return
        }

        ctx.SetStatusCode(result.Code)
        result.Release()
        return
    } else {
        h.writeResult(ctx, result.Res)
        result.Release()
        return
    }
}

func (h *HttpProxy) doProxy(ctx *fasthttp.RequestCtx, wg *sync.WaitGroup, result *model.RouteResult) {
    if nil != wg {
        defer wg.Done()
    }

    outReq := copyRequest(&ctx.Request)

    c := model.NewContext(h.routeTable, ctx, outReq, result)

    // 验证用户权限，同时获取用户id
    if result.API.NeedLogin {
        ok, err := h.CheckToken(outReq, result)
        if err != nil || !ok {
            result.Err = err
            result.Code = http.StatusForbidden
            return
        }
    }

    // 系统统一的filters
    filterName, code, err := filter.DoPreFilters(c, h.filters)
    if nil != err {
        log.Printf("Proxy Filter-Pre<%s> fail.", filterName, err)
        result.Err = err
        result.Code = code
        return
    }

    // pre filters
    filterName, code, err = filter.DoPreFilters(c, result.API.Filters)
    if nil != err {
        log.Printf("Proxy Filter-Pre<%s> fail.", filterName, err)
        result.Err = err
        result.Code = code
        return
    }

    service := result.API.Service
    c.SetStartAt(time.Now().UnixNano())
    if strings.ToUpper(string(c.GetOriginRequestCtx().Request.Header.Method())) == "OPTIONS" {

    } else if service.Protocol == "http" {

        outReq.Header.Set("client_id",string(outReq.PostArgs().Peek("client_id")))
        outReq.Header.Set("user_id",string(outReq.PostArgs().Peek("user_id")))

        res, err := h.fastHTTPClient.Do(outReq, service.GetHost())
        log.Println(outReq)
        c.SetEndAt(time.Now().UnixNano())
        result.Res = res

        if err != nil || res.StatusCode() >= fasthttp.StatusInternalServerError {
            resCode := http.StatusServiceUnavailable

            if nil != err {
                log.Printf("Proxy Fail <%s>", service.GetHost(), err)
            } else {
                resCode = res.StatusCode()
                log.Printf("Proxy Fail <%s>, Code <%d>", service.GetHost(), res.StatusCode(), err)
            }

            result.Err = err
            result.Code = resCode
            return
        }
        log.Printf("Backend server[%s] responsed, code <%d>, body<%s>", service.GetHost(), res.StatusCode(), res.Body())
    } else if service.Protocol == "thrift" {
        req := server.NewRequest()
        req.ServiceName = "businessService"

        // 解析参数，转化成json格式
        params := make(map[string]interface{})
        var f = func(k []byte, v []byte) {
            params[string(k)] = string(v)
        }

        if nil != outReq.Body() {
            if err = json.Unmarshal(outReq.Body(), &params); nil != err {
                log.Println("body json parse error")
            }
        }

        outReq.URI().QueryArgs().VisitAll(f)
        outReq.PostArgs().VisitAll(f)

        req.ParamJSON, _ = json.Marshal(params)
        req.Operation = result.API.Name

        pooledClient, err := service.Pool.Get()
        if err != nil {
            result.Err = err
            log.Println("Thrift pool get client error", err)
            return
        }
        defer service.Pool.Put(pooledClient, false)

        rawClient, ok := pooledClient.(*server.MyServiceClient)
        if !ok {
            log.Println("convert to raw client failed")
            return
        }
        res, err := rawClient.Send(req)
        c.SetEndAt(time.Now().UnixNano())

        if err != nil {
            result.Err = err
            log.Println(err)
            return
        }
        result.Res = &fasthttp.Response{}
        result.Res.SetStatusCode(int(res.ResponeCode))
        result.Res.SetBody(res.ResponseJSON)
    } else {
        return
    }

    // 系统统一的 post filters
    filterName, code, err = filter.DoPostFilters(c, h.filters)
    if nil != err {
        log.Printf("Proxy Filter-Post<%s> fail: %s ", filterName, err.Error())

        result.Err = err
        result.Code = code
        return
    }

    // api 自己的 post filters
    filterName, code, err = filter.DoPostFilters(c, result.API.Filters)
    if nil != err {
        log.Printf("Proxy Filter-Post<%s> fail: %s ", filterName, err.Error())

        result.Err = err
        result.Code = code
        return
    }

}

func (h *HttpProxy) writeResult(ctx *fasthttp.RequestCtx, res *fasthttp.Response) {
    ctx.SetStatusCode(res.StatusCode())
    ctx.Write(res.Body())
}
