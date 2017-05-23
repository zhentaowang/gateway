package filter

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

// Pre pre filter, before proxy request
func (f ResponseHeaderFilter) Post(c Context) (statusCode int, err error) {
    c.GetProxyResponse().Header.VisitAll(func(key, value []byte) {
        c.GetOriginRequestCtx().Response.Header.Set(string(key),string(value))
    })
    return f.BaseFilter.Pre(c)
}
