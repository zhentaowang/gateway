package filter

import "log"

/**
 * 把代理的响应头原封设置到Response Head里面
 */
type ResponseHeaderFilter struct{
    BaseFilter
}

func newResponseHeaderFilter() Filter {
    return &ResponseHeaderFilter{}
}

func (f ResponseHeaderFilter) Name() string {
    return FilterResponseHead
}

// set default response
func (f ResponseHeaderFilter) Post(c Context) (statusCode int, err error) {
    if c.GetProxyResponse()!=nil {
        c.GetProxyResponse().Header.VisitAll(func(key, value []byte) {
            c.GetOriginRequestCtx().Response.Header.Set(string(key),string(value))
        })
    } else {
        log.Println("response is null")
    }
    return f.BaseFilter.Pre(c)
}
