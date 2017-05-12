package filter

import (
    "strings"
)

/*
航班更新--系统间调用 需要放行
 */
type UpdateFlightFilter struct{
    BaseFilter
}

func newUpdateFlightFilter() Filter {
    return &UpdateFlightFilter{}
}

func (v UpdateFlightFilter) Name() string {
    return FilterUpdateFlight
}

func (v UpdateFlightFilter) Pre(c Context) (statusCode int, err error) {
    tcUpdate := make([]byte,0,2048)
    tcUpdate = append(c.GetOriginRequestCtx().PostBody()) // 获取到的是 byte[]
    str := "{"
    str += string(tcUpdate[:])
    str = strings.Replace(str,"data=","\"data\":",1)
    str = strings.Replace(str,"&sign=",",\"sign\":\"",1)
    str += "\"}"
    tcUpdate = []byte(str)
    c.GetProxyOuterRequest().SetBody(tcUpdate)
    return
}
