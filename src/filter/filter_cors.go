package filter

/*
添加Access-Control-Allow-Origin的filter
 */
type CORSFilter struct{
    BaseFilter
}

func newCORSFilter() Filter {
    return &CORSFilter{}
}

func (f CORSFilter) Name() string {
    return FilterCORS
}

// Pre pre filter, before proxy request
func (f CORSFilter) Post(c Context) (statusCode int, err error) {
    c.GetOriginRequestCtx().Response.Header.Add("Access-Control-Allow-Origin", "*")
    c.GetOriginRequestCtx().Response.Header.Add("Access-Control-Allow-Headers", "Content-Type,Access-Token")
    //c.GetOriginRequestCtx().Response.Header.Set("Content-Type", "application/json;charset=utf-8")
    return f.BaseFilter.Pre(c)
}
