package model

import (
    "encoding/base64"
    "github.com/valyala/fasthttp"
    "strings"
    "fmt"
    "regexp"
    "container/list"
    "gateway/src/filter"
)

const (
    // APIStatusDown down status
    APIStatusDown = iota // 0
    //APIStatusUp up status
    APIStatusUp // 1
)

// Mock mock
type Mock struct {
    Value         string             `json:"value"`
    ContentType   string             `json:"contentType, omitempty"`
    Headers       []*MockHeader      `json:"headers, omitempty"`
    Cookies       []string           `json:"cookies, omitempty"`
    ParsedCookies []*fasthttp.Cookie `json:"-"`
}

// MockHeader header
type MockHeader struct {
    Name  string `json:"name"`
    Value string `json:"value"`
}

// API a api define
type API struct {
    APIId         int            `json:"api_id"`
    Name          string         `json:"name, omitempty"`
    URI           string         `json:"uri"`
    Method        string         `json:"method"`
    NeedLogin     bool           `json:"need_login"`
    Status        int            `json:"status, omitempty"`
    ServiceId     int            `json:"service_id"`
    Service       *Service       `json:"service"`
    ServiceProviderName string   `json:"service_provider_name"`
    Mock          *Mock          `json:"mock, omitempty"`
    Desc          string         `json:"desc, omitempty"`
    filterNames   []string       `json:"-"`
    Pattern       *regexp.Regexp `json:"-"`
    Filters       *list.List     `json:"-"`
}

func (a *API) isUp() bool {
    return a.Status == APIStatusUp
}

func (a *API) matches(req *fasthttp.Request) bool {
    return a.isUp() && a.isMethodMatches(req) && a.isURIMatches(req)
}

func (a *API) isMethodMatches(req *fasthttp.Request) bool {
    return a.Method == "*" || strings.ToUpper(string(req.Header.Method())) == a.Method
}

func (a *API) isURIMatches(req *fasthttp.Request) bool {
        uri := strings.Split(string(req.URI().RequestURI()), "?")

        return a.Pattern.Match([]byte(uri[0]))
}

// RenderMock dender mock response
func (a *API) RenderMock(ctx *fasthttp.RequestCtx) {
    if a.Mock == nil {
        return
    }

    ctx.Response.Header.SetContentType(a.Mock.ContentType)

    if a.Mock.Headers != nil && len(a.Mock.Headers) > 0 {
        for _, header := range a.Mock.Headers {
            ctx.Response.Header.Add(header.Name, header.Value)
        }
    }

    if a.Mock.ParsedCookies != nil && len(a.Mock.ParsedCookies) > 0 {
        for _, ck := range a.Mock.ParsedCookies {
            ctx.Response.Header.SetCookie(ck)
        }
    }

    ctx.WriteString(a.Mock.Value)
}

func (a *API) init(services map[string]*Service) error {
    a.Pattern = regexp.MustCompile(a.URI)

    a.Filters = filter.NewFilters(a.filterNames)

    for _, v := range services {
        if v.ServiceId == a.ServiceId {
            a.Service = v
            return nil
        }
    }
    return ErrServiceNotFound
}

func (a *API) getKey() string {
    key := fmt.Sprintf("%s-%s", a.URI, a.Method)
    return base64.RawURLEncoding.EncodeToString([]byte(key))
}

// thrift 操作 设置operate,
// 设置规则：
//  1、当前api对象name属性为"dynamic-operate"
//  2、请求url或请求体必须要有key为operation的参数
func (a *API) GetOperation(operation string) string {
    if "dynamic-operate" == a.Name && len(operation) > 0{
        return operation
    }

    return a.Name
}