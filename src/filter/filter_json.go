package filter

type JSONFilter struct{
    BaseFilter
}

func newJSONFilter() Filter {
    return &JSONFilter{}
}

func (f JSONFilter) Name() string {
    return FilterCORS
}

// set json response
func (f JSONFilter) Post(c Context) (statusCode int, err error) {
    c.GetOriginRequestCtx().Response.Header.Set("Content-Type", "application/json;charset=utf-8")
    return f.BaseFilter.Pre(c)
}
