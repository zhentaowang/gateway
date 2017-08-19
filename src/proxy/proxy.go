package proxy

import (
    "github.com/valyala/fasthttp"
    "sync"
    "gateway/src/model"
    "net/http"
    "gateway/src/util"
    "code.aliyun.com/wyunshare/thrift-server/gen-go/server"
    "encoding/json"
    "time"
    "container/list"
    "gateway/src/filter"
    "log"
    "strings"
    "code.aliyun.com/wyunshare/thrift-server/conf"
    "strconv"
)

type HttpProxy struct {
    routeTable     *model.RouteTable
    store          model.Store
    fastHTTPClient *util.FastHTTPClient
    filters        *list.List
}

func NewHttpProxy(store model.Store) *HttpProxy {

    cf := util.GetConfigCenterInstance()

    conf.TConfig = conf.T{}

    conf.TConfig.MaxConnDuration, _= strconv.Atoi(cf.ConfProperties["jdbc"]["max_conn_duration"])
    conf.TConfig.MaxConns, _= strconv.Atoi(cf.ConfProperties["jdbc"]["max_conns"])
    conf.TConfig.MaxIdle, _= strconv.Atoi(cf.ConfProperties["jdbc"]["max_idle"])
    conf.TConfig.MaxIdleConnDuration, _= strconv.Atoi(cf.ConfProperties["jdbc"]["max_idle_conn_duration"])
    conf.TConfig.MaxResponseBodySize, _= strconv.Atoi(cf.ConfProperties["jdbc"]["max_response_body_size"])
    conf.TConfig.ReadTimeout, _ = strconv.Atoi(cf.ConfProperties["jdbc"]["read_timeout"])
    conf.TConfig.ReadTimeout, _ = strconv.Atoi(cf.ConfProperties["jdbc"]["write_timeout"])
    conf.TConfig.ReadBufferSize, _ = strconv.Atoi(cf.ConfProperties["jdbc"]["read_buffer_size"])
    conf.TConfig.WriteBufferSize, _ = strconv.Atoi(cf.ConfProperties["jdbc"]["write_buffer_size"])
    log.Println("jdbc配置成功")

    h := &HttpProxy{
        fastHTTPClient: util.NewFastHTTPClient(&conf.TConfig),
        store: store,
    }

    h.Init()

    return h
}

func (h *HttpProxy) Init() {
    defer util.ErrHandle()
    err := h.InitRouteTable()

    if err != nil {
        log.Panic(err, "init route table error")
    }

    filterNames, _ := h.store.GetFilters(-1)
    h.filters = filter.NewFilters(filterNames)
}

func (h *HttpProxy) InitRouteTable() error {

    h.routeTable = model.NewRouteTable(h.store)
    h.routeTable.Load()

    return nil
}

func (h *HttpProxy) Start() {
    log.Printf("Proxy exit at %s", fasthttp.ListenAndServe(":8888", h.ReverseProxyHandler))
}

func (h *HttpProxy) ReverseProxyHandler(ctx *fasthttp.RequestCtx) {

    defer util.ErrHandle()


    args := ctx.QueryArgs()
    isTest := false

    if strings.Compare(string(args.Peek("type")),"gateway_test")==0 {
        isTest = true
    } else {
        log.Println("网关开始工作，请求的url = " + string(ctx.Request.RequestURI()) + " \n 请求的HEAD=" + ctx.Request.Header.String() + " \n  请求的 body = " + string(ctx.Request.Body()[:]))
    }
    result := h.routeTable.Select(&ctx.Request)

    if nil == result {
        ctx.SetStatusCode(fasthttp.StatusNotFound)
        log.Println("请求的url没有找到,或者请求方法不对，或者url为不可用状态")
        return
    }

    h.doProxy(ctx, nil, result)

    if result.Err != nil {
        if result.API.Mock != nil {
            result.API.RenderMock(ctx)
            if isTest == false {
                    log.Println("网关结束处理  "+string(ctx.Request.RequestURI())+ "，返回的是mock数据， HEAD = " + result.Res.Header.String())
            }
            result.Release()
            return
        }

        ctx.SetStatusCode(result.Code)
        if isTest == false {
            if result.Res!=nil {
                log.Println("网关结束处理  "+string(ctx.Request.RequestURI())+ "，  出错，返回的响应为 HEAD = " + result.Res.Header.String()+",error="+result.Err.Error())
            } else {
                log.Println("网关结束处理  "+string(ctx.Request.RequestURI())+"，  出错，返回的响应为空,error="+result.Err.Error())
            }
        }

        result.Release()
        return
    } else {
        h.writeResult(ctx, result.Res)
        if isTest == false {
            if result.Res!=nil {
                log.Println("网关结束处理  "+string(ctx.Request.RequestURI())+ "，返回的响应为 HEAD = " + result.Res.Header.String())
            } else {
                log.Println("网关结束处理  "+string(ctx.Request.RequestURI())+",返回的响应为空")
            }
        }

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
            if err != nil {
                log.Println("认证中心认证失败  "+err.Error())
            } else {
                log.Println("认证中心认证失败  ")
            }

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
    } else {
        log.Printf("Proxy Filter-Pre<%s> success.", filterName, err)
    }

    // pre filters
    filterName, code, err = filter.DoPreFilters(c, result.API.Filters)
    if nil != err {
        log.Printf("Proxy Filter-Pre<%s> fail.", filterName, err)
        result.Err = err
        result.Code = code
        return
    } else {
        log.Printf("Proxy Filter-Pre<%s> success.", filterName, err)
    }

    service := result.API.Service
    c.SetStartAt(time.Now().UnixNano())
    if strings.ToUpper(string(c.GetOriginRequestCtx().Request.Header.Method())) == "OPTIONS" {
        log.Println("Request Method="+string(c.GetOriginRequestCtx().Request.Header.Method()))
    } else if service.Protocol == "http" {

        outReq.Header.Set("client_id",string(outReq.PostArgs().Peek("client_id")))
        outReq.Header.Set("user_id",string(outReq.PostArgs().Peek("user_id")))

        res, err := h.fastHTTPClient.Do(outReq, service.GetHost())
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
    } else if service.Protocol == "thrift" {
        req := server.NewRequest()

        // 解析serviceName
        req.ServiceName = "businessService"// 默认servicename = businessService
        if serviceProviderName := result.API.ServiceProviderName; len(serviceProviderName) > 0 {
            req.ServiceName = serviceProviderName
        }

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

        // 转化成json
        delete(params, "access_token")
        req.ParamJSON, _ = json.Marshal(params)

        // set operation
        operation := result.API.Name
        if value, ok := params["operation"].(string); ok {
            operation = result.API.GetOperation(value)
            log.Println("Request's Operation="+operation)
        }
        req.Operation = operation

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

        log.Println("网关处理thrift请求，paramJson= "+string(req.ParamJSON)+"  ,operation= "+req.Operation+"  ,ServiceName="+req.ServiceName)
        res, err := rawClient.Send(req)
	if res != nil {
		log.Println("网关结束处理thrift请求，ResponseCode="+strconv.FormatInt(int64(res.ResponeCode),10))
	} else {
		log.Println("网关结束处理thrift请求，返回的响应为空")
	}
        c.SetEndAt(time.Now().UnixNano())

        if err != nil {
            result.Err = err
            log.Println("处理thrift请求失败  "+err.Error())
            return
        }
        result.Res = &fasthttp.Response{}
        if res != nil {
            result.Res.SetStatusCode(int(res.ResponeCode))
            result.Res.SetBody(res.ResponseJSON)
        } else {
            result.Res.SetStatusCode(500)
        }
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
    if res != nil {
        ctx.SetStatusCode(res.StatusCode())
        ctx.Write(res.Body())
    } else {
        log.Println("the reponse is null")
    }
}