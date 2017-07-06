package filter

/*
添加Access-Control-Allow-Origin的filter
 */
type TextFilter struct{
    BaseFilter
}

func newTextFilter() Filter {
    return &TextFilter{}
}

func (f TextFilter) Name() string {
    return FilterCORS
}

// Pre pre filter, before proxy request
func (f TextFilter) Post(c Context) (statusCode int, err error) {
    c.GetOriginRequestCtx().Response.Header.Set("Content-Type", "text/html;charset=utf-8")
    return f.BaseFilter.Pre(c)
}
