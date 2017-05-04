package proxy

import (
    "github.com/labstack/gommon/log"
    "github.com/valyala/fasthttp"
    "sync"
    "gateway/model"
    "net/http"
    "gateway/util"
    "gateway/thrift/gen-go/server"
    "encoding/json"
    "gateway/config"
    "time"
)

type HttpProxy struct {
    routeTable     *model.RouteTable
    store          model.Store
    fastHTTPClient *util.FastHTTPClient
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
}

func (h *HttpProxy) initRouteTable() error {

    h.routeTable = model.NewRouteTable(h.store)
    h.routeTable.Load()

    return nil
}

func (h *HttpProxy) Start() {
    log.Error(fasthttp.ListenAndServe(":8081", h.ReverseProxyHandler), "Proxy exit at %s", )
}

func (h *HttpProxy) ReverseProxyHandler(ctx *fasthttp.RequestCtx) {
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

    // pre filters
    filterName, code, err := result.API.DoPreFilters(c)
    if nil != err {
        log.Warnf("Proxy Filter-Pre<%s> fail.", filterName, err)
        result.Err = err
        result.Code = code
        return
    }

    service := result.API.Service
    c.SetStartAt(time.Now().UnixNano())
    if service.Protocol == "http" {
        res, err := h.fastHTTPClient.Do(outReq, service.GetHost())
        log.Info(outReq)
        c.SetEndAt(time.Now().UnixNano())
        result.Res = res

        if err != nil || res.StatusCode() >= fasthttp.StatusInternalServerError {
            resCode := http.StatusServiceUnavailable

            if nil != err {
                log.Infof("Proxy Fail <%s>", service.GetHost(), err)
            } else {
                resCode = res.StatusCode()
                log.Infof("Proxy Fail <%s>, Code <%d>", service.GetHost(), res.StatusCode(), err)
            }

            result.Err = err
            result.Code = resCode
            return
        }
        log.Infof("Backend server[%s] responsed, code <%d>, body<%s>", service.GetHost(), res.StatusCode(), res.Body())
    } else if service.Protocol == "thrift" {
        req := server.NewRequest()
        req.ServiceName = service.GetHost()

        // 解析参数，转化成json格式
        params := make(map[string]interface{})
        var f = func(k []byte, v []byte) {
            params[string(k)] = string(v)
        }

        if nil != outReq.Body() {
            if err = json.Unmarshal(outReq.Body(), &params); nil != err {
                log.Info("body json parse error")
            }
        }
        outReq.PostArgs().VisitAll(f)
        outReq.URI().QueryArgs().VisitAll(f)
        // 如果header中有User-Id, 需要加入参数中
        if userId := outReq.Header.Peek("User-Id"); userId != nil {
            params["user_id"] = string(userId)
        }

        req.ParamJSON, _ = json.Marshal(params)

        pooledClient, err := service.Pool.Get()
        if err != nil {
            result.Err = err
            log.Error("Thrift pool get client error", err)
            return
        }
        defer service.Pool.Put(pooledClient, false)

        rawClient, ok := pooledClient.(*server.MyServiceClient)
        if !ok {
            log.Error("convert to raw client failed")
            return
        }
        res, err := rawClient.Send(req)
        c.SetEndAt(time.Now().UnixNano())

        if err != nil {
            result.Err = err
            log.Error(err)
            return
        }
        result.Res = &fasthttp.Response{}
        result.Res.SetStatusCode(int(res.ResponeCode))
        result.Res.SetBody(res.ResponseJSON)
        log.Info(res)
    } else {
        return
    }

    // post filters
    filterName, code, err = result.API.DoPostFilters(c)
    if nil != err {
        log.Infof("Proxy Filter-Post<%s> fail: %s ", filterName, err.Error())

        result.Err = err
        result.Code = code
        return
    }
}

func (h *HttpProxy) writeResult(ctx *fasthttp.RequestCtx, res *fasthttp.Response) {
    ctx.SetStatusCode(res.StatusCode())
    ctx.Write(res.Body())
}
