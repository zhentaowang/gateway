package model

import (
    "errors"
    "github.com/valyala/fasthttp"
    "git.apache.org/thrift.git/lib/go/thrift"
    "log"
)

var (
    ErrAPIExists = errors.New("API already exist")
    ErrServiceExists = errors.New("Service already exist")
    ErrAPINotFound = errors.New("API not found")
    ErrServiceNotFound = errors.New("Service not found")
)

type RouteTable struct {
    store Store

    apis map[string]*API
    services map[string]*Service

    transportFactory thrift.TTransportFactory
    protocolFactory *thrift.TBinaryProtocolFactory
}

func NewRouteTable(store Store) *RouteTable {
    rt := &RouteTable{
        apis: make(map[string]*API),
        services: make(map[string]*Service),

        transportFactory: thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory()),
        protocolFactory : thrift.NewTBinaryProtocolFactoryDefault(),

        store: store,
    }

    return rt
}

func (r *RouteTable) Select(req *fasthttp.Request) *RouteResult {
    for _, api := range r.apis {
        if api.matches(req) {
            return &RouteResult{
                API: api,
            }
        }
    }

    return nil
}

func (r *RouteTable) Load() {
    r.loadServices()
    r.loadAPIs()
}

func (r *RouteTable) loadServices() {
    // 先清空 services map
    for k := range r.services {
        delete(r.services, k)
    }

    services, err := r.store.GetServices()
    if nil != err {
        log.Panic(err, "Load services fail.")
    }

    for _, service := range services {
        err := r.AddNewService(service)
        if nil != err {
            log.Panic(err, "Service <%s> add fail", service.getKey())
        }
    }
}

func (r *RouteTable) loadAPIs() {
    // 先清空 apis map
    for k := range r.apis {
        delete(r.apis, k)
    }

    apis, err := r.store.GetAPIs()
    if nil != err {
        log.Panic(err, "Load apis fail.")
        return
    }

    for _, api := range apis {
        err := r.AddNewAPI(api)
        if nil != err {
            log.Panic(err, "API <%s> add fail", api.getKey())
        }
    }
}

// AddNewService add a new service
func (r *RouteTable) AddNewService(service *Service) error {
    key := service.getKey()
    _, ok := r.services[key]

    if ok {
        return ErrServiceExists
    }

    err := service.init(r)
    if nil != err {
        log.Panic(err, "Service init error")
    }

    r.services[key] = service

    log.Printf("Service <%s-%s> added", service.Namespace, service.Name)
    return nil
}

// AddNewAPI add a new API
func (r *RouteTable) AddNewAPI(api *API) error {
    apiKey := api.getKey()
    _, ok := r.apis[apiKey]

    if ok {
        return ErrAPIExists
    }

    err := api.init(r.services)
    if nil != err {
        log.Panic(err, "API init error")
    }

    r.apis[apiKey] = api

    log.Printf("API <%s-%s> added", api.Method, api.URI)

    return nil
}

type RouteResult struct {
    API *API
    Err error
    Code int
    Res *fasthttp.Response
}

func (r *RouteResult) Release() {
    if nil != r.Res {
        fasthttp.ReleaseResponse(r.Res)
    }
}