package model

import (
    "encoding/base64"
    "github.com/valyala/fasthttp"
    "strings"
    "fmt"
    "regexp"
    "container/list"
    "gateway/src/filter"
    "github.com/labstack/gommon/log"
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
    Mock          *Mock          `json:"mock, omitempty"`
    Desc          string         `json:"desc, omitempty"`
    Pattern       *regexp.Regexp `json:"-"`
    filterNames   string         `json:"filter_names"`
    filters       *list.List     `json:"-"`
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
    return a.Pattern.Match(req.URI().RequestURI())
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

    a.initFilters()

    for _, v := range services {
        if v.ServiceId == a.ServiceId {
            a.Service = v
            return nil
        }
    }
    return ErrServiceNotFound
}

func (a *API) initFilters() {
    a.filters = list.New()
    if a.filterNames == "" {
        return
    }

    for _, filterName := range strings.Split(a.filterNames, ",") {
        f, err := filter.NewFilter(filterName)
        if nil != err {
            log.Panicf("Proxy unknow filter <%+v>", filterName)
        }

        log.Info(f)
        a.filters.PushBack(f)
    }
}

func (a *API) getKey() string {
    key := fmt.Sprintf("%s-%s", a.URI, a.Method)
    return base64.RawURLEncoding.EncodeToString([]byte(key))
}