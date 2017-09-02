package filter

import (
    "net/http"

    "github.com/valyala/fasthttp"
    "container/list"
    "log"
)


// Context filter context
type Context interface {
    SetStartAt(startAt int64)
    SetEndAt(endAt int64)
    GetStartAt() int64
    GetEndAt() int64

    GetProxyOuterRequest() *fasthttp.Request
    GetProxyResponse() *fasthttp.Response

    GetOriginRequestCtx() *fasthttp.RequestCtx
}

// Filter filter interface
type Filter interface {
    Name() string

    Pre(c Context) (statusCode int, err error)
    Post(c Context) (statusCode int, err error)
    PostErr(c Context)
}

// BaseFilter base filter support default implemention
type BaseFilter struct{}

// Pre execute before proxy
func (f BaseFilter) Pre(c Context) (statusCode int, err error) {
    return http.StatusOK, nil
}

// Post execute after proxy
func (f BaseFilter) Post(c Context) (statusCode int, err error) {
    return http.StatusOK, nil
}

// PostErr execute proxy has errors
func (f BaseFilter) PostErr(c Context) {

}

func DoPreFilters(c Context, filters *list.List) (filterName string, statusCode int, err error) {
    for item := filters.Front(); item != nil; item = item.Next() {
        f, _ := item.Value.(Filter)
        filterName = f.Name()

        statusCode, err = f.Pre(c)
        if nil != err {
            log.Printf("Proxy Filter-Pre<%s> fail.", filterName, err)
            return filterName, statusCode, err
        }
        log.Printf("Proxy Filter-Pre<%s> success.", filterName)
    }

    return "", http.StatusOK, nil
}

func DoPostFilters(c Context, filters *list.List) (filterName string, statusCode int, err error) {
    for item := filters.Back(); item != nil; item = item.Prev() {
        f, _ := item.Value.(Filter)

        statusCode, err = f.Post(c)
        if nil != err {
            log.Printf("Proxy Filter-Post<%s> fail: %s ", filterName, err.Error())
            return filterName, statusCode, err
        }

        log.Printf("Proxy Filter-Post<%s> success ", filterName)
    }

    return "", http.StatusOK, nil
}

func DoPostErrFilters(c Context, filters *list.List) {
    for item := filters.Back(); item != nil; item = item.Prev() {
        f, _ := item.Value.(Filter)

        f.PostErr(c)
    }
}
